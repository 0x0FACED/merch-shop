package model

type Item struct {
	Type     string `json:"type" db:"type"`
	Quantity uint   `json:"quantity" db:"quantity"`
}
