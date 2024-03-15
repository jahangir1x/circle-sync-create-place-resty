package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
	"os"
	"strconv"
	"strings"
)

func login() (string, error) {
	type LoginReq struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		LongLived bool   `json:"long_lived"`
	}

	type LoginResp struct {
		UserId       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	payload := LoginReq{
		Email:     "jahangir64r@gmail.com",
		Password:  "Passw0rd",
		LongLived: true,
	}

	var loginResp LoginResp
	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(&payload).
		SetResult(&loginResp).
		Post("http://16.171.41.66:6977/api/log-in")
	if err != nil {
		return "", err
	}

	fmt.Println("Status:", response.Status())
	return loginResp.AccessToken, nil
}

type CsvParsedData struct {
	LocationName string  `csv:"Location name"`
	Coordinates  string  `csv:"Coordinates"`
	Radius       float64 `csv:"Radius (Meters)"`
}

func parseCsvToStruct(csvPath string) ([]CsvParsedData, error) {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}

	var csvParsedData []CsvParsedData
	if err := gocsv.UnmarshalFile(csvFile, &csvParsedData); err != nil {
		return nil, err
	}
	return csvParsedData, nil
}

func createPlace(accessToken string, placeName string, latitude float64, longitude float64, radius float64) error {
	type CreatePlaceReq struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Radius    float64 `json:"radius"`
	}

	payload := CreatePlaceReq{
		Name:      placeName,
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radius,
	}

	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(&payload).
		Post("http://16.171.41.66:6977/api/place")
	fmt.Println("Status:", response.Status())
	if err != nil {
		return err
	}
	fmt.Println(response.String())
	return nil
}

func main() {
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) != 1 {
		fmt.Println("Usage: go run main.go <csv_file_path>")
		return
	}

	accessToken, err := login()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Logged in successfully")

	csvFilePath := arguments[0]
	fmt.Println("Processing file: ", csvFilePath)
	csvParsed, err := parseCsvToStruct(csvFilePath)
	if err != nil {
		fmt.Println(err)
		return

	}
	for _, place := range csvParsed {
		latitude, err := strconv.ParseFloat(strings.Split(place.Coordinates, ", ")[0], 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		longitude, err := strconv.ParseFloat(strings.Split(place.Coordinates, ", ")[1], 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = createPlace(accessToken, place.LocationName, latitude, longitude, place.Radius)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
