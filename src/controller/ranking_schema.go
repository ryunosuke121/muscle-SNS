package controller

type GetUserTotalWeightInMonthRequestSchema struct {
	Year  uint `json:"year" validate:"required"`
	Month uint `json:"month" validate:"required"`
}

type GetMonthRankingInGroupRequestSchema struct {
	Year  uint `json:"year" validate:"required"`
	Month uint `json:"month" validate:"required"`
}

type GetMonthRankingInGroupByMenuResponseSchema struct {
	MenuID uint `json:"menu_id" validate:"required"`
	Year   uint `json:"year" validate:"required"`
	Month  uint `json:"month" validate:"required"`
}
