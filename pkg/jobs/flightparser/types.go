package flightparser

import "github.com/adh-partnership/api/pkg/geo"

type Facility struct {
	ID       string      `json:"id"`
	Boundary [][]float64 `json:"coords"`
	Polygon  geo.Polygon
}
