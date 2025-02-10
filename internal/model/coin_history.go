package model

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	User   string `json:"fromUser"`
	Amount int    `json:"amount"`
}

type SentTransaction struct {
	User   string `json:"toUser"`
	Amount int    `json:"amount"`
}
