package models

import "cloud.google.com/go/firestore"

type CreateDeck struct {
	Title        string   `json:"title" validate:"required"`
	OwnerID      string   `json:"owner_id" validate:"required"`
	SharedEmails []string `json:"shared_email" validate:"omitempty,dive,email"`
}

type DeckResponse struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	OwnerID      string   `json:"owner_id"`
	SharedEmails []string `json:"shared_emails"`
	Cards        []Card   `json:"cards"`
}

type Deck struct {
	OwnerID      string                   `firestore:"owner_id"`
	Title        string                   `firestore:"title"`
	SharedEmails []string                 `firestore:"shared_emails"`
	Cards        []*firestore.DocumentRef `firestore:"cards"`
}

type UpdateDeck struct {
	Title     string   `json:"title,omitempty"`
	Operation string   `json:"operation,omitempty" validate:"omitempty,oneof=add remove"`
	Cards     []string `json:"cards,omitempty"`
}

type UpdateCards struct {
	Operation string   `json:"operation,omitempty" validate:"omitempty,oneof=add remove"`
	Cards     []string `json:"cards,omitempty"`
}
