package modelrequests

type CreateRequest struct {
	Title    string   `json:"title" validate:"required"`
	Content  string   `json:"content" validate:"required"`
	Category string   `json:"category" validate:"required"`
	Tags     []string `json:"tags" validate:"required"`
}

type UpdateRequest struct {
	Title    string   `json:"title" validate:"required"`
	Content  string   `json:"content" validate:"required"`
	Category string   `json:"category" validate:"required"`
	Tags     []string `json:"tags" validate:"required"`
}
