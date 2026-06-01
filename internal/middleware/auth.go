package middleware

import (
	"errors"
	"net/http"
	"strings"

	"eao/internal/config"

	jwtpkg "eao/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "userID"
)

var (
	ErrorUserNotLogin = errors.New("当前用户未登录")
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	// 从请求中读取 Authorization 或 token，校验并解析 JWT。
	// 解析成功后，将用户ID写入上下文：key 为 ContextUserIDKey。
	// 若未携带 token 或 token 无效，则返回统一错误响应并中止后续处理。
	return func(c *gin.Context) {
		// 0) 跳过 CORS 预检与白名单请求
		if c.Request.Method == "OPTIONS" || isRequestWhitelisted(c) {
			c.Next()
			return
		}
		// 1) 优先从 Authorization 头读取，支持 "Bearer <token>" 格式
		authHeader := c.GetHeader("Authorization")
		var token string
		if authHeader != "" {
			if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				token = strings.TrimSpace(authHeader[7:])
			} else {
				// 非 Bearer 直接作为 token 尝试解析
				token = strings.TrimSpace(authHeader)
			}
		}

		// 2) 其次尝试从查询参数读取
		if token == "" {
			token = strings.TrimSpace(c.Query("token"))
		}
		// 3) 再次尝试从 Cookie 读取
		if token == "" {
			if ck, err := c.Cookie("token"); err == nil {
				token = strings.TrimSpace(ck)
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1006,
				"message": "未登录或登录过期",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 解析并验证 token
		claims, err := jwtpkg.ParseToken(token)
		if err != nil || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1006,
				"message": "登录状态失效",
				"data":    nil,
			})
			c.Abort()
			return
		}
		if !claims.IsAccessToken() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1006,
				"message": "登录状态失效",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将用户ID写入上下文，供后续处理使用
		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}

// isRequestWhitelisted 判断当前请求是否在访问白名单内。
// 白名单元素支持两种格式：
// 1) "METHOD:/path"（严格匹配方法与路径）
// 2) "/path"（仅匹配路径，方法不限）
// 另外支持方法为 "*" 的写法："*:/path" 表示任意方法匹配该路径。
func isRequestWhitelisted(c *gin.Context) bool {
	cfg := config.GetConfig()
	if cfg == nil || len(cfg.Auth.Whitelist) == 0 {
		return false
	}
	reqMethod := strings.ToUpper(c.Request.Method)
	reqPath := c.FullPath()
	if reqPath == "" {
		// 当未使用命名路由时，FullPath 可能为空，退回到实际请求路径
		reqPath = c.Request.URL.Path
	}
	for _, item := range cfg.Auth.Whitelist {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		// 格式：METHOD:/path
		if strings.Contains(item, ":") {
			parts := strings.SplitN(item, ":", 2)
			method := strings.ToUpper(strings.TrimSpace(parts[0]))
			path := strings.TrimSpace(parts[1])
			if (method == reqMethod || method == "*") && path == reqPath {
				return true
			}
			continue
		}
		// 仅路径匹配
		if item == reqPath {
			return true
		}
	}
	return false
}
