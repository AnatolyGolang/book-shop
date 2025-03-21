package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	r "github.com/AnatolyGolang/book-shop/internal/app/repositories"
	sm "github.com/AnatolyGolang/book-shop/internal/app/services/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your_secret_key")

type Claims struct {
	Id      int    `json:"user_id"`
	Email   string `json:"role"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type JWTServiceImpl struct {
	repository r.TokenRepository
}

func NewJWTService(repo r.TokenRepository) *JWTServiceImpl {
	return &JWTServiceImpl{
		repository: repo,
	}
}

func (s *JWTServiceImpl) GenerateJWT(ctx context.Context, user sm.DomainUser) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		Id:      user.Id,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	err = s.repository.SaveToken(ctx, user.Id, signedToken, expirationTime)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *JWTServiceImpl) GetUser(ctx context.Context, token string) (sm.DomainUser, error) {
	var userClaims Claims
	parsedJwt, err := jwt.ParseWithClaims(token, &userClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid singning method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return sm.DomainUser{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !parsedJwt.Valid {
		return sm.DomainUser{}, errors.New("invalid token")
	}

	return claimsToUser(userClaims), nil
}

func (s *JWTServiceImpl) RevokeToken(ctx context.Context, token string) error {
	return s.repository.DeleteToken(ctx, token)
}

func (s *JWTServiceImpl) StartTokenCleanupScheduler() {
	go func() {
		for {
			err := s.repository.CleanupExpiredTokens(context.Background())
			if err != nil {
				log.Printf("failed to cleanup expired tokens: %v", err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
}

func claimsToUser(claims Claims) sm.DomainUser {
	return sm.DomainUser{
		Id:      claims.Id,
		Email:   claims.Email,
		IsAdmin: claims.IsAdmin,
	}
}
