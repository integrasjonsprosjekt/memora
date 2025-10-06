package models

import "cloud.google.com/go/firestore"

type CreateDeck struct {
	Title        string   `json:"title" validate:"required" firestore:"title"`
	OwnerID      string   `json:"owner_id" validate:"required" firestore:"owner_id"`
	SharedEmails []string `json:"shared_emails" validate:"omitempty,dive,email" firestore:"shared_emails"` 
}

type DeckResponse struct {
	ID           string   `json:"id" firestore:"-"`
	Title        string   `json:"title" firestore:"title"`
	OwnerID      string   `json:"owner_id" firestore:"owner_id"`
	SharedEmails []string `json:"shared_emails" firestore:"shared_emails"`
	Cards        []Card   `json:"cards" firestore:"cards"`
}

type Deck struct {
	OwnerID      string                   `firestore:"owner_id"`
	Title        string                   `firestore:"title"`
	SharedEmails []string                 `firestore:"shared_emails"`
	Cards        []*firestore.DocumentRef `firestore:"cards"`
}

type UpdateDeck struct {
	Title     string   `json:"title,omitempty" firestore:"title"`
	OppCards  string   `json:"opp_cards,omitempty" validate:"omitempty,oneof=add remove"`
	Cards     []string `json:"cards,omitempty" firestore:"cards"`
	OppEmails string   `json:"opp_emails,omitempty" validate:"omitempty,oneof=add remove"`
	Emails    []string `json:"emails,omitempty" firestore:"shared_emails"`
}
