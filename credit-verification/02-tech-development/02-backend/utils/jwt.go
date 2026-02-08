package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义JWT声明（包含用户核心信息）
type Claims struct {
	UserID               uint64 `json:"user_id"` // 用户ID
	Address              string `json:"address"` // 以太坊地址
	Role                 string `json:"role"`    // 用户角色
	jwt.RegisteredClaims        // 内置声明（过期时间、签发时间等）
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint64, address, role string) (string, error) {
	// 从全局配置获取过期时间和密钥
	expireHours := GlobalConfig.JWT.ExpireHours
	secret := GlobalConfig.JWT.Secret

	// 构造声明
	claims := Claims{
		UserID:  userID,
		Address: address,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                             // 签发时间
			Issuer:    "credit-dapp",                                                              // 签发者
		},
	}

	// 生成token（HS256算法）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析并验证JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("不支持的签名算法")
		}
		// 返回密钥
		return []byte(GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token并提取声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token无效")
}
