package tests

import (
	"encoding/json"
	"github.com/xdire/weathersum/handlers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// mockNWSServer creates a test server that mimics the National Weather Service API
type mockNWSServer struct {
	server *httptest.Server
}

// newMockNWSServer sets up a mock NWS server with predefined responses
func newMockNWSServer() *mockNWSServer {
	mock := &mockNWSServer{}

	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/points/40.7128,-74.006":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"properties": map[string]interface{}{
					"gridId": "NYK",
					"gridX":  33,
					"gridY":  37,
				},
			})
		case "/gridpoints/NYK/33,37/forecast":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"properties": map[string]interface{}{
					"periods": []map[string]interface{}{
						{
							"name":          "Today",
							"temperature":   65,
							"shortForecast": "Partly Cloudy",
						},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))

	parsedURL, err := url.Parse(mock.server.URL)
	if err != nil {
		panic("failed to parse mock server url")
	}
	http.DefaultClient = &http.Client{
		Transport: &mockRoundTripper{
			originalTransport: http.DefaultTransport,
			mockBaseURL:       parsedURL.Host,
		},
	}

	return mock
}

// mockRoundTripper intercepts and redirects NWS API calls to our mock server
type mockRoundTripper struct {
	originalTransport http.RoundTripper
	mockBaseURL       string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Replace the original NWS URL with our mock server URL
	req.URL.Scheme = "http"
	req.URL.Host = m.mockBaseURL

	return m.originalTransport.RoundTrip(req)
}

// Close shuts down the mock server
func (m *mockNWSServer) Close() {
	m.server.Close()
}

// TestWeatherHandler tests the weather endpoint
func TestWeatherHandler(t *testing.T) {
	// Setup mock NWS server
	mockNWS := newMockNWSServer()
	defer mockNWS.Close()

	// Test cases
	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Coordinates",
			url:            "/weather?lat=40.7128&lon=-74.0060",
			expectedStatus: http.StatusOK,
			expectedBody:   "For Today expecting moderate temperature and Partly Cloudy",
		},
		{
			name:           "Missing Latitude",
			url:            "/weather?lon=-74.0060",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing Longitude",
			url:            "/weather?lat=40.7128",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Latitude",
			url:            "/weather?lat=invalid&lon=-74.0060",
			expectedStatus: http.StatusBadRequest,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the handler directly
			handlers.SimplifiedWeather(rr, req)

			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// If we expect a successful response, validate the body
			if tc.expectedStatus == http.StatusOK {
				var forecast handlers.ForecastResult
				err := json.Unmarshal(rr.Body.Bytes(), &forecast)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Compare key fields
				if forecast.Forecast != tc.expectedBody {
					t.Errorf("Incorrect short forecast: got %v want %v",
						forecast.Forecast, tc.expectedBody)
				}
			}
		})
	}
}
