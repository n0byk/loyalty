package entity

import "time"

type Order struct {
	OrderNumber string    `json:"number" db:"order_number"`
	OrderState  string    `json:"status" db:"order_state"`
	Accrual     float32   `json:"accrual"`
	UpdateTime  time.Time `json:"uploaded_at" db:"update_time"`
}

type OrderIDNumber struct {
	OrderID     string `json:"order_id" db:"order_id"`
	OrderNumber string `json:"order_number" db:"order_number"`
}

type OrderWithdrawals struct {
	OrderID     string    `json:"order" db:"order_number"`
	OrderSum    float32   `json:"sum" db:"order_number"`
	ProcessedAt time.Time `json:"uploaded_at" db:"update_time"`
}
