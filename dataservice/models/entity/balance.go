package entity

type Balance struct {
	Current   float32 `json:"current" db:"order_number"`
	Withdrawn float32 `json:"withdrawn" db:"order_number"`
}
