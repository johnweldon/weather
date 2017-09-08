package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	baseURL           = "https://api.openweathermap.org/data/2.5/"
	appID             = ""
	latitude  float64 = -20.272967
	longitude float64 = 30.934364
)

func init() {
	flag.StringVar(&baseURL, "baseurl", baseURL, "Base URL for weather service")
	flag.StringVar(&appID, "appid", appID, "APPID token")
	flag.Float64Var(&latitude, "latitude", latitude, "latitude")
	flag.Float64Var(&longitude, "longitude", longitude, "longitude")
}

func main() {
	flag.Parse()
	p := weatherParams{
		BaseURL: baseURL,
		APPID:   appID,
		Coord: Coordinates{
			Lat: latitude,
			Lon: longitude,
		},
	}

	cur, err := currentWeather(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cur)

	fore, err := forecast(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", fore)
}

func currentWeather(p weatherParams) (*CurrentWeather, error) {
	raw, err := apiCall(p, "weather")
	if err != nil {
		return nil, err
	}

	cur := &CurrentWeather{}
	err = json.Unmarshal(raw, cur)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func forecast(p weatherParams) (*Forecast, error) {
	raw, err := apiCall(p, "forecast")
	if err != nil {
		return nil, err
	}

	cur := &Forecast{}
	err = json.Unmarshal(raw, cur)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

type weatherParams struct {
	BaseURL string
	APPID   string
	Coord   Coordinates
}

func apiCall(p weatherParams, path string) ([]byte, error) {
	u, err := url.Parse(p.BaseURL + path)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("APPID", p.APPID)
	q.Set("lat", fmt.Sprintf("%f", p.Coord.Lat))
	q.Set("lon", fmt.Sprintf("%f", p.Coord.Lon))

	u.RawQuery = q.Encode()

	client := http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("Got %q", resp.Status)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

type CurrentWeather struct {
	Coord      Coordinates `json:"coord"`
	Weather    []Info      `json:"weather"`
	Base       string      `json:"base"`
	Main       Reading     `json:"main"`
	Visibility int         `json:"visibility"`
	Wind       Wind        `json:"wind"`
	Clouds     Clouds      `json:"clouds"`
	Dt         int         `json:"dt"`
	Sys        Sys2        `json:"sys"`
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Cod        int         `json:"cod"`
}

type Info struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Reading struct {
	Temp     float64 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

type Sys2 struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

type Forecast struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []Item  `json:"list"`
	City    City    `json:"city"`
}

type Item struct {
	Dt      int       `json:"dt"`
	Main    Main      `json:"main"`
	Weather []Weather `json:"weather"`
	Clouds  Clouds    `json:"clouds"`
	Wind    Wind      `json:"wind"`
	Sys     Sys       `json:"sys"`
	DtTxt   string    `json:"dt_txt"`
	Rain    Rain      `json:"rain,omitempty"`
}
type Main struct {
	Temp      float64 `json:"temp"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Clouds struct {
	All int `json:"all"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type Sys struct {
	Pod string `json:"pod"`
}

type Rain struct {
	ThreeH float64 `json:"3h"`
}

type City struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Coord      Coordinates `json:"coord"`
	Country    string      `json:"country"`
	Population int         `json:"population"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
