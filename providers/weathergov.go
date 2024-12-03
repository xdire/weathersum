package providers

import (
	"encoding/json"
	"fmt"
	"github.com/xdire/weathersum/forecast"
	"net/http"
	"strconv"
	"strings"
)

const (
	wgovPeriodToday         string = "today"
	wgovPeriodThisAfternoon string = "this afternoon"
	wgovPeriodTonight       string = "night"
)

type WeatherGov struct {
	identity string
	client   *http.Client
}

func NewWeatherGov(identity string) *WeatherGov {
	return &WeatherGov{
		identity: identity,
		// TODO Provide a customization to build this provider with the custom client if needed
		client: http.DefaultClient,
	}
}

// GetGridpoint fetches the gridpoint data and provides gridId and X,Y coordinates in the output
func (pwg *WeatherGov) GetGridpoint(lat, lon float64) (string, string, error) {
	// Format float should trim trailing zeros as Sprintf printing out additional precision in the mantissa
	// @see https://stackoverflow.com/questions/31289409/format-a-float-to-n-decimal-places-without-trailing-zeros
	url := fmt.Sprintf("https://api.weather.gov/points/%s,%s",
		strconv.FormatFloat(lat, 'f', -1, 64),
		strconv.FormatFloat(lon, 'f', -1, 64))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("User-Agent", pwg.identity)
	resp, err := pwg.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var pointData struct {
		Properties struct {
			GridId string `json:"gridId"`
			GridX  int    `json:"gridX"`
			GridY  int    `json:"gridY"`
		} `json:"properties"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&pointData); err != nil {
		return "", "", err
	}

	return pointData.Properties.GridId,
		fmt.Sprintf("%d,%d", pointData.Properties.GridX, pointData.Properties.GridY),
		nil
}

// GetForecast retrieves the forecast for a specific grid point and coordinates (@see GetGridpoint)
func (pwg *WeatherGov) GetForecast(gridId, gridPoint string) (forecast.Forecast, error) {
	url := fmt.Sprintf("https://api.weather.gov/gridpoints/%s/%s/forecast", gridId, gridPoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", pwg.identity)
	resp, err := pwg.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var forecastData struct {
		Properties struct {
			Periods []struct {
				Name          string `json:"name"`
				Temperature   int    `json:"temperature"`
				ShortForecast string `json:"shortForecast"`
			} `json:"periods"`
		} `json:"properties"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
		return nil, err
	}

	// Get today's forecast (first period)
	if len(forecastData.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast data available")
	}

	out := &forecast.Simplified{}

	// Do transform of the provider data to weathersum API type of the data
	for _, period := range forecastData.Properties.Periods {
		normalizedPeriod := strings.ToLower(period.Name)
		representation := forecast.SimplifiedPeriod{}
		if strings.Contains(normalizedPeriod, wgovPeriodThisAfternoon) {
			representation.SetKind(forecast.PeriodAfternoon)
		} else if strings.Contains(normalizedPeriod, wgovPeriodToday) {
			representation.SetKind(forecast.PeriodToday)
		} else if strings.Contains(normalizedPeriod, wgovPeriodTonight) {
			representation.SetKind(forecast.PeriodTonight)
		} else {
			// We exhausted the available variances for the next forecast
			break
		}
		representation.SetTemp(period.Temperature)
		representation.SetShortDesc(period.ShortForecast)
		out.AddPeriod(representation)
	}

	return out, nil
}
