package model

import "time"

type SalesRecord struct {
	UserName         string    `json:"user_name"`
	Age              int       `json:"age"`
	Height           float32   `json:"height"`
	Gender           string    `json:"gender"`
	SaleAmount       float32   `json:"sale_amount"`
	LastPurchaseDate time.Time `json:"last_purchase_date"`
}
