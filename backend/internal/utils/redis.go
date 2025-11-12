package utils

const (
	UserKeyPrefix = "user"
	DeckKeyPrefix = "deck"
	CardKeyPrefix = "card"
)

func UserKey(userID string) string {
	return UserKeyPrefix + ":" + userID
}

func DeckKey(deckID string) string {
	return DeckKeyPrefix + ":" + deckID
}

func DeckCardsKey(deckID string) string {
	return DeckKeyPrefix + ":" + deckID + ":cards"
}

func DeckCardKey(deckID, cardID string) string {
	return DeckKeyPrefix + ":" + deckID + ":" + CardKeyPrefix + ":" + cardID
}

func UserEmailDecksKey(email string) string {
	return "user:email:" + email + ":decks"
}

func UserKeyRateLimit(userID string) string {
	return "rate_limit:user:" + userID
}
