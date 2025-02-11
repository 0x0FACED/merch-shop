package model

type Inventory struct {
	Items []Item `json:"inventory"`
}

type Item struct {
	Type     string `json:"type" db:"type"`
	Quantity uint   `json:"quantity" db:"quantity"`
}
