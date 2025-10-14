package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type FirebaseAuth interface {
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
}

type FirebaseAuthRepo struct {
	client *auth.Client
}

func NewFirebaseAuth(app *firebase.App) (*FirebaseAuthRepo, error) {
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	return &FirebaseAuthRepo{client: client}, nil
}

func (f *FirebaseAuthRepo) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := f.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}
