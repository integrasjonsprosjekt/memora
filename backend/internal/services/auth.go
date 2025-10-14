package services

import (
	"context"
	"memora/internal/firebase"
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

func (s AuthService) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	return s.repo.VerifyIDToken(ctx, idToken)
}
