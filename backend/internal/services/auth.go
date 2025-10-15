package services

import (
	"context"
	"memora/internal/firebase"

	"firebase.google.com/go/auth"
)

type AuthService struct {
	repo firebase.FirebaseAuth
}

func NewAuthService(
	repo firebase.FirebaseAuth,
) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s AuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return s.repo.VerifyIDToken(ctx, idToken)
}
