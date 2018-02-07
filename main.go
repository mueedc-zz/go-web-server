package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"flag"
	"log"
	"time"
)

func main() {
	wundergroundAPIKey := flag.String("wunderground.api.key", "0123456789abcdef", "wunderground.com API key")
	flag.Parse()

	mw := multiWeatherProvider{
		openWeatherMap{},
		weatherUnderground{apiKey: *wundergroundAPIKey},
	}

	http.HandleFunc("/weather/", func (w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		temp, err := mw.temperature(city) {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"city": city,
			"temp": temp,
			"took": time.Since(begin).String()
		})
	})

	http.ListenAndServe(":8080", nil)
}

type weatherProvider interface {
	temperature(city string) (float64, error)
}

type openWeatherMap struct{}

func (w openWeatherMap) temperature (city string) (float64, error) {
	resp, err := htpp.Get("http://api.openweathermap.org/data/2.5/weather?APIID=YOUR_API_KEY&q" + city)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}


func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?ADDID=YOUR_API_KEY&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()
	
	var d weatherData

	if err := json.NewDecoder((resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	})

	return d, nil
}

