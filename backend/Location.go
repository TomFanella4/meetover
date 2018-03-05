package main

// Address is a our location metric
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Area  string `json:"area,omitempty"`
}

// Geolocation stores the user's location in lat, long and time from last update
type Geolocation struct {
	ID        string `json:"uid"`
	Coord     Coord  `json:"coord,omitempty"`
	TimeStamp int64  `json:"timestamp,omitempty"`
}

// Coord - latitide and longitude
type Coord struct {
	Lat  string `json:"lat,omitempty"`
	Long string `json:"long,omitempty"`
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
	testCoords1 := Coord{
		Lat:  "12101",
		Long: "20302",
	}
	testCoords2 := Coord{
		Lat:  "10101",
		Long: "20302",
	}
	testCoords3 := Coord{
		Lat:  "10101",
		Long: "2023302",
	}
	testCoords4 := Coord{
		Lat:  "10101001",
		Long: "20302",
	}
	testCoords5 := Coord{
		Lat:  "1",
		Long: "20302",
	}
	addGeolocation("loc-test1", testCoords1, makeTimestamp())
	addGeolocation("loc-test2", testCoords2, makeTimestamp())
	addGeolocation("loc-test3", testCoords3, makeTimestamp())
	addGeolocation("loc-test4", testCoords4, makeTimestamp())
	addGeolocation("loc-test5", testCoords5, makeTimestamp())
}
