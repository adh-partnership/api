package dto

type StaffingRequest struct {
	Date             string `json:"date" binding:"required"`
	DepartureAirport string `json:"departureAirport" binding:"required"`
	ArrivalAirport   string `json:"arrivalAirport" binding:"required"`
	Pilots           int    `json:"pilots" binding:"required"`
	Comments         string `json:"comments"`
}
