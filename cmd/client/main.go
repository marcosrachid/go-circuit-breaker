package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var settings gobreaker.Settings
	settings.Name = "HTTP GET"
	settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 10 && failureRatio >= 0.6
	}
	settings.Timeout = 2 * time.Millisecond
	settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		if to == gobreaker.StateOpen {
			log.Error().Msg("State Open!")
		}
		if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
			log.Info().Msg("Going from Open to Half Open")
		}
		if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
			log.Info().Msg("Goind from Half Open to Closed!")
		}
	}
	cb = gobreaker.NewCircuitBreaker(settings)
}

func main() {
	urlIncorrect := "http://localhost:8081"
	urlCorrect := "http://localhost:8080"
	var body []byte
	var err error
	for i := 0; i < 20; i++ {
		body, err = Get(urlIncorrect)
		if err != nil {
			log.Error().Err(err).Msg("Error")
		}
		fmt.Println(string(body))
		if i > 15 {
			urlIncorrect = urlCorrect
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("http Get request gave error")
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
	if err != nil {
		return nil, err
	}
	return body.([]byte), nil
}
