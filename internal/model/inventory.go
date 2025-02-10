package model

type Inventory struct {
	Type     string `json:"type" db:"type"`
	Quantity uint   `json:"quantity" db:"quantity"`
}
