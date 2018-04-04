package main

import (
	"fmt"
	"math"
	"time"
)

// Address is a our location metric
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Area  string `json:"area,omitempty"`
}

// Geolocation - latitide and longitude and last time of update
type Geolocation struct {
	Lat       float64 `json:"lat,omitempty"`
	Long      float64 `json:"long,omitempty"`
	TimeStamp float64 `json:"timestamp,omitempty"`
}

// QueryLocation will return the location for the given coordinates
func QueryLocation(coords string) (Address, error) {
	// long := strings.Split(coords, ",")[0]
	// lat := strings.Split(coords, ",")[1]

	// Temprary place holder
	location := Address{City: "Chicago", State: "IL", Area: "ORD"}
	return location, nil
}

// InRadius - checks if distance between two points is within radius
func InRadius(center Geolocation, point Geolocation, radius float64) bool {
	radius = radius * 1000 // convert to meters
	dist := Distance(center.Lat, center.Long, point.Lat, point.Long)
	if dist < radius {
		return true
	}
	fmt.Printf("dist: %f, radius: %f\n", dist, radius)
	return false
}

// hsin - haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance returns the meters between coord1 and coord2
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180
	r = 6378100 // Earth radius in METERS
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
