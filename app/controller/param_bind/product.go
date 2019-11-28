package param_bind

type ProductAdd struct {
	Name string `json:"name" form:"name" binding:"required"`
}
