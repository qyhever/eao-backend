package model

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminRefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type AdminLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
