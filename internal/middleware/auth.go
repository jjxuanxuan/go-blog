package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

const jwtIssuer = "go-blog"

func jwtSecret() []byte {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "change_me_dev_secret"
	}
	return []byte(sec)
}

func accessTTL() time.Duration {
	minStr := os.Getenv("ACCESS_TOKEN_TTL")
	if mins, err := strconv.Atoi(minStr); err == nil && mins > 0 {
		return time.Duration(mins) * time.Minute
	}
	return 120 * time.Minute
}

// GenerateToken 用于登录成功后签发 access_token
func GenerateToken(userId uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "sub": userId,
            "exp": time.Now().Add(accessTTL()).Unix(),
            "iat": time.Now().Unix(),
            "iss": jwtIssuer,
        })
    return token.SignedString(jwtSecret())
}

// AuthMiddleware 校验 Authorization: Bearer <token>
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "缺少或非法Token"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return jwtSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "无效Token"})
			c.Abort()
			return
		}

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "解析Token失败"})
            c.Abort()
            return
        }

        // 校验签发方（Issuer）
        if iss, ok := claims["iss"].(string); !ok || iss != jwtIssuer {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "无效Token"})
            c.Abort()
            return
        }

		sub, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token中缺少sub"})
			c.Abort()
			return
		}
		c.Set("user_id", uint(sub))
		c.Next()
	}
}
