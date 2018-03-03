package main

// Address is a our location metric
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Area  string `json:"area,omitempty"`
}

// QueryLocation will return the location for the given coordinates
func QueryLocation(coords string) (Address, error) {
	// long := strings.Split(coords, ",")[0]
	// lat := strings.Split(coords, ",")[1]

	// Temprary place holder
	location := Address{City: "Chicago", State: "IL", Area: "ORD"}
	return location, nil
}
