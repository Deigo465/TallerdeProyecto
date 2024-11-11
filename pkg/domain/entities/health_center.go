package entities

type HealthCenter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func NewHealthCenter(id int, name, district, address string) HealthCenter {
	return HealthCenter{id, name, district, address}
}

// helper method for testing purposes
func NewFakeHealthCenter() HealthCenter {
	return HealthCenter{1, "name", "district", "address"}
}
