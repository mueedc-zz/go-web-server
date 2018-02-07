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

type weatherUnderground struct {
	apiKey string
}

func (w weatherUnderground) temperature (city string) (float64, error) {
	resp, err := http.Get("http://api.weatherunderground.com/api/" + w.apiKey + "/conditions/q/" + city + ".json")
	if err != nil {
		return 0, nil
	}

	defer resp.Body.Close()

	var d struct {
		Observation struct {
			Celsius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	kelvin := d.Observation.Celsius + 273.15
	log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
	return kelvin, nil
}

type multiWeatherProvider []weatherProvider

func (w  multiWeatherProvider) temperature (city string) {
	// individual channels for temperatures and errors. each will push a value into only one
	temps := make(chan float64, len(w))
	errs := make(chan error, len(w))

	// for each provider start go routine that invokes temperature method and forwards the response
	for _, provider := range w {
		go func(p weatherProvider) {
			k, err := p.temperature(city)
			if err != nil {
				errs <- err
				return
			}
			temps <- k
		}(provider)
	}

	sum := 0.0

	// collect temperature from each provider
	for i := 0; i < len(w); i++ {
		select {
		case temp := <- temps:
			sum += temp
		case err := <- errs:
			return 0, err
		}
	}

	// return the average temperature
	return sum / float64(len(providers)), nil
}
