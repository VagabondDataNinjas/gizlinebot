package domain

type Location struct {
	Id        uint64  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationPrice struct {
	Location
	Price float64 `json:"price"`
}

type LocationThai struct {
	Location
	NameThai string `json:"name_thai"`
}
