package dto

type StaffingRequest struct {
	Date             string `json:"date" binding:"required"`
	Start            string `json:"start" binding:"required"`
	End              string `json:"end" binding:"end"`
	DepartureAirport string `json:"departure" binding:"required"`
	ArrivalAirport   string `json:"arrival" binding:"required"`
	Pilots           string `json:"pilots" binding:"required"`
	Comments         string `json:"comments"`
}
