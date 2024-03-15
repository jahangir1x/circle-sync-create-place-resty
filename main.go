package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
	"os"
	"resty_experiments/serializer"
	"strconv"
	"strings"
)

func login() (string, error) {
	payload := serializer.LoginReq{
		Email:     "jahangir64r@gmail.com",
		Password:  "Passw0rd",
		LongLived: true,
	}

	var loginResp serializer.LoginResp
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

func parseCsvToStruct(csvPath string) ([]serializer.CsvData, error) {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}

	var csvParsedData []serializer.CsvData
	if err := gocsv.UnmarshalFile(csvFile, &csvParsedData); err != nil {
		return nil, err
	}
	return csvParsedData, nil
}

func createPlace(accessToken string, placeName string, latitude float64, longitude float64, radius float64) error {
	payload := serializer.CreatePlaceReq{
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

func createMultiplePlaces(accessToken string, places []serializer.CsvData) error {
	for _, place := range places {
		latitude, err := strconv.ParseFloat(strings.Split(place.Coordinates, ", ")[0], 64)
		if err != nil {
			return err
		}
		longitude, err := strconv.ParseFloat(strings.Split(place.Coordinates, ", ")[1], 64)
		if err != nil {
			return err
		}

		err = createPlace(accessToken, place.LocationName, latitude, longitude, place.Radius)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseCsvPathFromCmd() (string, error) {
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) != 1 {
		return "", fmt.Errorf("Usage: go run main.go <csv_file_path>")
	}
	return arguments[0], nil
}

func main() {
	csvFilePath, err := parseCsvPathFromCmd()
	if err != nil {
		fmt.Println(err)
		return
	}

	accessToken, err := login()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Logged in successfully")

	fmt.Println("Processing file: ", csvFilePath)
	csvData, err := parseCsvToStruct(csvFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = createMultiplePlaces(accessToken, csvData)
	if err != nil {
		fmt.Println(err)
		return
	}
}
