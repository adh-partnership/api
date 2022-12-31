package dto

type AirportWeatherDTO struct {
	ID    string `json:"id"`
	METAR string `json:"metar"`
	TAF   string `json:"taf"`
}
