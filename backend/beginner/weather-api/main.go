package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	glide "github.com/valkey-io/valkey-glide/go/v2"
	"github.com/valkey-io/valkey-glide/go/v2/config"
	"github.com/valkey-io/valkey-glide/go/v2/servermodules/glidejson"
	"golang.org/x/time/rate"
)

var ApiKey = os.Getenv("WEATHER_VISUALCROSSING")

func setWeather(client *glide.Client, ctx context.Context, key string, bw []byte) error {
	ok, err := glidejson.JsonSet(client, ctx, key, "$", string(bw))
	if err != nil || ok != "OK" {
		return err
	}
	return nil
}

// key city:scc (state or country code)
func getWeather(client *glide.Client, ctx context.Context, key string) ([]byte, error) {
	resp, err := glidejson.JsonGet(client, ctx, key)
	if err != nil {
		fmt.Println("here, not in the http handler: ", err)
		return nil, err
	}

	var body []byte
	if resp != "" {
		return []byte(resp), nil

	} else {
		temp := strings.Split(key, ":")
		city := temp[0]
		scc := temp[1]

		res, err := http.Get(fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s,%s?unitGroup=us&key=%s&contentType=json", city, scc, ApiKey))
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			break
		case http.StatusTooManyRequests:
			fmt.Println("Rate limit exceeded")
			return nil, err
		default:
			panic(fmt.Sprintf("unexpected code: %d", res.StatusCode))
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("reading resp body into byte string: ", err)
			return nil, err
		}

		err = setWeather(client, ctx, key, body)
		if err != nil {
			fmt.Println("set inside get weather: ", err)
		}

	}
	return body, nil
}

func handleWeather(client *glide.Client) http.HandlerFunc {
	// 3 per second, bursting 1
	var limiter = rate.NewLimiter(rate.Limit(3), 1)
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded"))
			return
		}
		ctx := context.Background()
		city := r.URL.Query().Get("city")
		scc := r.URL.Query().Get("scc")

		temp := fmt.Sprintf("%s:%s", city, scc)
		weather, err := getWeather(client, ctx, temp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error retrieving weather"))
			fmt.Println("weather cache/retrieval in http handler: ", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(weather)
	}
}

func main() {

	vhost := "localhost"
	vport := 6379
	config := config.NewClientConfiguration().WithAddress(&config.NodeAddress{Host: vhost, Port: vport})

	client, err := glide.NewClient(config)
	if err != nil {
		fmt.Println("valkey new client: ", err)
		return
	}
	defer client.Close()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	//localhost:8080/weather?city=Newyork&scc=US
	http.HandleFunc("/weather", handleWeather(client))

	server.ListenAndServe()

}
