package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/http"
)

type weatherData struct {
	weather *jester.Data
	time    time.Time
}

var weatherCache = map[string]weatherData{}

// GetWeather fetches the weather data for a given location.
func GetWeather(location string) (*jester.Data, error) {
	if w, ok := weatherCache[location]; ok && time.Since(w.time) < time.Minute*5 {
		return w.weather, nil
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?appid=533d00f29739edba54c5fb2b1a8262b3&q=%s",
		url.QueryEscape(location),
	)

	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		return nil, err
	}

	data := resp.Json()

	if data.Get("cod").MustInt() != 200 {
		return nil, fmt.Errorf("invalid location")
	}

	weatherCache[location] = weatherData{
		weather: data,
		time:    time.Now(),
	}

	return data, nil
}
