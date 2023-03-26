package dto

type AccountStoreRequest struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type AccountUpdateRequest struct {
	Password string `json:"password" validate:"required,min=8,max=32"`
}
