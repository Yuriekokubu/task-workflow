package model

import "github.com/Yuriekokubu/workflow/internal/constant"

type RequestItem struct {
	Title    string `json:"title"`
	Amount   int    `json:"amount"`
	Quantity uint   `json:"quantity"`
	OwnerID  uint   `json:"owner_id"`
}

type RequestFindItem struct {
	Statuses constant.ItemStatus `form:"status"`
}

type RequestUpdateItem struct {
	Status constant.ItemStatus
}

type RequestLogin struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}
