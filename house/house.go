package house

//House database struct
type house struct {
	HouseID *int `json:"house_id,omitempty"`
	Name *string `json:"house_name,omitempty"`
	Founder *string `json:"founder,omitempty"`
	Animal *string `json:"animal,omitempty"`
}