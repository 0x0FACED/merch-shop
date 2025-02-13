package model

type User struct {
	ID       uint   `db:"id"`
	Username string `db:"username"`
	Password string `db:"password_hash"`
}

type UserInfo struct {
	Coins       uint `db:"balance"`
	Inventory   []Item
	CoinHistory CoinHistory
}
