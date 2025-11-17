package dto

type CategoryResp struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentId *uint  `json:"parent_id,omitempty"`
}

type CreateCategoryReq struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Slug     string `json:"slug" binding:"required,min=1,max=100"`
	ParentId *uint  `json:"parent_id"`
}
