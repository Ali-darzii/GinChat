package serializer

import "mime/multipart"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	PhoneNo string `binding:"required,phone_validator,max=11,min=11" json:"phone_no"`
}

type LoginRequest struct {
	PhoneNo string `binding:"required,phone_validator" json:"phone_no"`
	Token   int    `binding:"required" json:"token"`
	Name    string `binding:"name_validator" json:"name"`
}

type ProfileUpdateRequest struct {
	ID       uint                  `json:"id" form:"id"`
	Avatar   *multipart.FileHeader `binding:"image_validator" json:"avatar" form:"avatar"`
	Name     string                `binding:"required,name_validator" json:"name" form:"name"`
	Username string                `json:"username" form:"username"`
}

type UpdatedProfile struct {
	ID       uint   `json:"id" form:"id"`
	Avatar   string `json:"avatar" form:"avatar"`
	Name     string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
}
