package jwt

import (
	"eao/internal/config"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID    uint64 `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// getSecret 返回用于签名/校验 JWT 的密钥。
// 延迟从配置中取值，避免在包初始化阶段访问未加载的配置导致空指针。
func getSecret() []byte {
	cfg := config.GetConfig()
	if cfg == nil || cfg.JWT.Secret == "" {
		return []byte("")
	}
	return []byte(cfg.JWT.Secret)
}

// keyFunc 提供给 jwt.Parse/ParseWithClaims 使用的密钥回调。
func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return getSecret(), nil
}

func parseDurationOrFallback(raw string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return duration
}

// GenToken 生成access token 和 refresh token
func GenToken(userID uint64) (aToken, rToken string, err error) {
	// 获取配置
	cfg := config.GetConfig()

	// 解析访问令牌过期时间
	accessExpireDuration := parseDurationOrFallback(cfg.JWT.AccessExpiresIn, time.Hour*8)

	// 解析刷新令牌过期时间
	refreshExpireDuration := parseDurationOrFallback(cfg.JWT.RefreshExpiresIn, time.Hour*24*3)

	// 创建一个我们自己的声明
	c := MyClaims{
		UserID:    userID,
		TokenType: TokenTypeAccess,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessExpireDuration).Unix(), // 过期时间
			Issuer:    "jeve",                                      // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(getSecret())

	rc := MyClaims{
		UserID:    userID,
		TokenType: TokenTypeRefresh,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshExpireDuration).Unix(), // 过期时间
			Issuer:    "jeve",                                       // 签发人
		},
	}
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rc).SignedString(getSecret())
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

// IsAccessToken reports whether the token is an access token.
func (c *MyClaims) IsAccessToken() bool {
	return c != nil && c.TokenType == TokenTypeAccess
}

// IsRefreshToken reports whether the token is a refresh token.
func (c *MyClaims) IsRefreshToken() bool {
	return c != nil && c.TokenType == TokenTypeRefresh
}
