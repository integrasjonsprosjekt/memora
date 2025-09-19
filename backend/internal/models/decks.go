package models

import "cloud.google.com/go/firestore"

type CreateDeck struct {
	Title        string   `json:"title"`
	OwnerID      string   `json:"owner_id"`
	SharedEmails []string `json:"shared_email"`
}

type DeckResponse struct {
	ID           string   `json:"id"`
	OwnerID      string   `json:"owner_id"`
	SharedEmails []string `json:"shared_emails" firestore:"shared_emails"`
	Cards        []Card   `json:"cards"`
}

type Deck struct {
	OwnerID      string                   `firestore:"owner_id"`
	SharedEmails []string                 `firestore:"shared_emails"`
	Cards        []*firestore.DocumentRef `firestore:"cards"`
}
