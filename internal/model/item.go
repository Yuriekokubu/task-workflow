package model

import (
	"time"

	"github.com/Yuriekokubu/workflow/internal/constant"
)

type Item struct {
	ID        uint                `json:"id" gorm:"primaryKey"`
	Title     string              `json:"title"`
	Amount    int                 `json:"amount"`
	Quantity  uint                `json:"quantity"`
	Status    constant.ItemStatus `json:"status"`
	OwnerID   uint                `json:"owner_id"`
	UpdatedAt time.Time           `json:"updated_at"`
}
