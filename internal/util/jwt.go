// Package util/jwt 封装 JWT 的生成与解析，包含 Access/Refresh TTL 与角色信息。
package util

import (
    "github.com/golang-jwt/jwt/v5"
    "os"
    "strconv"
    "time"
)

const jwtIssuer = "go-blog"

// Claims 自定义声明，包含角色与标准注册字段。
type Claims struct {
    Role string `json:"role"`
    jwt.RegisteredClaims
}

func jwtSecret() []byte {
    sec := os.Getenv("JWT_SECRET")
    if sec == "" {
        sec = "change_me_dev_secret"
    }
    return []byte(sec)
}

func ttlFromEnv(key string, defMins int) time.Duration {
    if v := os.Getenv(key); v != "" {
        if m, err := strconv.Atoi(v); err == nil && m > 0 {
            return time.Duration(m) * time.Minute
        }
    }
    return time.Duration(defMins) * time.Minute
}

// AccessTTL 访问令牌有效期（分钟），默认 120 分钟，可通过 ACCESS_TOKEN_TTL 配置。
func AccessTTL() time.Duration  { return ttlFromEnv("ACCESS_TOKEN_TTL", 120) }
// RefreshTTL 刷新令牌有效期（分钟），默认 7 天，可通过 REFRESH_TOKEN_TTL 配置。
func RefreshTTL() time.Duration { return ttlFromEnv("REFRESH_TOKEN_TTL", 7*24*60) }

// GenerateAccessToken 生成短期访问令牌，包含用户ID与角色。
func GenerateAccessToken(userID uint, role string) (string, error) {
    now := time.Now()
    claims := &Claims{
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    jwtIssuer,
            Subject:   strconv.FormatUint(uint64(userID), 10),
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(AccessTTL())),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret())
}

// GenerateRefreshToken 生成长期刷新令牌，包含唯一 jti，适合黑名单/旋转策略。
func GenerateRefreshToken(userID uint, role string) (string, error) {
    now := time.Now()
    claims := &Claims{
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    jwtIssuer,
            Subject:   strconv.FormatUint(uint64(userID), 10),
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTTL())),
            ID:        strconv.FormatInt(now.UnixNano(), 10),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret())
}

// ParseToken 解析并校验 token，返回自定义 Claims。
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        return jwtSecret(), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrTokenInvalidClaims
}
