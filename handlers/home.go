package handlers

import "net/http"

func APIHome(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`
{
	"paths": {
		"weather_v1": "/v1/weather"
	}
}
	`))
	if err != nil {
		return
	}
}
