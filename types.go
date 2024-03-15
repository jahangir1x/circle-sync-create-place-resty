package main

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

type CsvData struct {
	LocationName string  `csv:"Location name"`
	Coordinates  string  `csv:"Coordinates"`
	Radius       float64 `csv:"Radius (Meters)"`
}

type CreatePlaceReq struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}
