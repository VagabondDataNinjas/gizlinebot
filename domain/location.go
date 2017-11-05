package domain

type Location struct {
	Id        uint64
	Name      string
	Latitude  float64
	Longitude float64
}

type LocationPrice struct {
	Location
	Price float64
}
