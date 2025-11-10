// Package middleware/auth 负责解析与校验 JWT，将用户ID与角色注入上下文。
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

// GenerateToken 用于登录成功后签发 access_token（旧版简单实现，仍兼容）。
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

// AuthMiddleware 校验 Authorization: Bearer <token> 并注入上下文字段：
// user_id：uint 类型的用户ID（支持 sub 为数字或字符串）
// role：字符串角色（可选）
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

        if subNum, ok := claims["sub"].(float64); ok {
            c.Set("user_id", uint(subNum))
        } else if subStr, ok := claims["sub"].(string); ok {
            var n uint64
            for i := 0; i < len(subStr); i++ {
                ch := subStr[i]
                if ch < '0' || ch > '9' { n = 0; break }
                n = n*10 + uint64(ch-'0')
            }
            if n == 0 {
                c.JSON(http.StatusUnauthorized, gin.H{"message": "Token中缺少sub"})
                c.Abort()
                return
            }
            c.Set("user_id", uint(n))
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Token中缺少sub"})
            c.Abort()
            return
        }
        if role, rok := claims["role"].(string); rok {
            c.Set("role", role)
        }
        c.Next()
    }
}
