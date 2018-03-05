package main

import "time"

// Address is a our location metric
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Area  string `json:"area,omitempty"`
}

// Geolocation - latitide and longitude and last time of update
type Geolocation struct {
	ID        string `json:"uid,omitempty"` // TODO: we probably don't need this feild once the GeoLocation is in the Person Struct
	Lat       string `json:"lat,omitempty"`
	Long      string `json:"long,omitempty"`
	TimeStamp int64  `json:"timestamp,omitempty"`
}

// QueryLocation will return the location for the given coordinates
func QueryLocation(coords string) (Address, error) {
	// long := strings.Split(coords, ",")[0]
	// lat := strings.Split(coords, ",")[1]

	// Temprary place holder
	location := Address{City: "Chicago", State: "IL", Area: "ORD"}
	return location, nil
}

// AddTestCoordsToDB used in the test/ endpoint to see db stores the coordinates
func AddTestCoordsToDB() {
	testCoords1 := Geolocation{
		ID:        "test1",
		Lat:       "12101",
		Long:      "20302",
		TimeStamp: makeTimestamp(),
	}
	testCoords2 := Geolocation{
		ID:        "test2",
		Lat:       "10101",
		Long:      "20302",
		TimeStamp: makeTimestamp(),
	}
	testCoords3 := Geolocation{
		ID:        "test3",
		Lat:       "10101",
		Long:      "2023302",
		TimeStamp: makeTimestamp(),
	}
	testCoords4 := Geolocation{
		ID:        "test4",
		Lat:       "10101001",
		Long:      "20302",
		TimeStamp: makeTimestamp(),
	}
	testCoords5 := Geolocation{
		ID:        "test5",
		Lat:       "1",
		Long:      "20302",
		TimeStamp: makeTimestamp(),
	}
	addGeolocation(testCoords1)
	addGeolocation(testCoords2)
	addGeolocation(testCoords3)
	addGeolocation(testCoords4)
	addGeolocation(testCoords5)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
