package serializer

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
