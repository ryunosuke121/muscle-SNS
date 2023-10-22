package controller

type CreatePostRequestSchema struct {
	Comment  string                `json:"comment"`
	Training *CreateTrainingSchema `json:"training"`
}

type CreateTrainingSchema struct {
	MenuID uint `json:"menu_id"`
	Times  uint `json:"times"`
	Weight uint `json:"weight"`
	Sets   uint `json:"sets"`
}
