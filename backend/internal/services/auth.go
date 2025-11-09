package services

import (
	"context"
	"memora/internal/firebase"

	firebaseauth "firebase.google.com/go/auth"
)

type AuthService struct {
	repo firebase.FirebaseAuth
}

func NewAuthService(
	deps *ServiceDeps,
) *AuthService {
	return &AuthService{
		repo: deps.AuthRepo,
	}
}

func (s AuthService) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*firebaseauth.Token, error) {
	return s.repo.VerifyIDToken(ctx, idToken)
}
