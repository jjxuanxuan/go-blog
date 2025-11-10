// Package util/hash 提供密码哈希与校验封装（bcrypt）。
package util

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 生成密码哈希，默认成本（Cost）。
func HashPassword(pw string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

// CheckPassword 比对明文密码与哈希是否匹配。
func CheckPassword(hashed, plain string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
