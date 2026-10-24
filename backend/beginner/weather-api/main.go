package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	glide "github.com/valkey-io/valkey-glide/go/v2"
	"github.com/valkey-io/valkey-glide/go/v2/config"
)

type Days struct {
	DateTime    string  `json:"datetime"`
	TempMax     float64 `json:"tempmax"`
	TempMin     float64 `json:"tempmin"`
	Temp        float64 `json:"temp"`
	FeelsLike   float64 `json:"feelslike"`
	Humidity    float64 `json:"humidity"`
	WindSpeed   float64 `json:"windspeed"`
	Conditions  string  `json:"conditions"`
	Description string  `json:"description"`
}
type Weather struct {
	Address  string `json:"address"`
	Timezone string `json:"timezone"`
	Days     []Days `json:"days"`
}

func main() {

	host := "localhost"
	port := 6379

	config := config.NewClientConfiguration().WithAddress(&config.NodeAddress{Host: host, Port: port})

	client, err := glide.NewClient(config)
	if err != nil {
		fmt.Println("valkey new client: ", err)
		return
	}
	defer client.Close()
	ctx := context.Background()
	//TODO: working on valkey(redis)
	client.Set(ctx, "", "")

	apiKey := os.Getenv("WEATHER_VISUALCROSSING")

	resp, err := http.Get("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Tucson,AZ?unitGroup=us&key=" + apiKey + "&contentType=json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		fmt.Println("Rate limit exceeded")
		return
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("reading resp body into byte string: ", err)
		return
	}

	wRes := Weather{}
	if err := json.Unmarshal(body, &wRes); err != nil {
		fmt.Println("unmarshaling json: ", err)
		return
	}

	fmt.Println(wRes)

}
