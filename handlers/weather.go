package handlers

import (
	"encoding/json"
	"github.com/xdire/weathersum/providers"
	"net/http"
	"strconv"
)

type ForecastResult struct {
	Forecast string `json:"forecast"`
}

func SimplifiedWeather(w http.ResponseWriter, r *http.Request) {
	// Ensure only GET requests are accepted
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	// Validate coordinates
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	// TODO Add provider selector for the cases when customer wants different forecast service
	// TODO Provider selector also can select different provider if current provider cannot provide forecast
	wp := providers.NewWeatherGov("weathersum")

	// Get gridpoint
	gridId, gridPoint, err := wp.GetGridpoint(lat, lon)
	if err != nil {
		http.Error(w, "Error getting gridpoint: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get forecast
	forecast, err := wp.GetForecast(gridId, gridPoint)
	if err != nil {
		http.Error(w, "Error getting forecast: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ForecastResult{Forecast: forecast.AsString()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
