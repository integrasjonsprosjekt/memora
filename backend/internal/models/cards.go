package models

import "memora/internal/utils"

type Card interface {
	GetType() string
}

type MutlipleChoiceCard struct {
	ID      string          `json:"id,omitempty" firestore:"-"`
	Type    string          `json:"type" validate:"required" firestore:"type"`
	Options map[string]bool `json:"options" validate:"required" firestore:"options"`
}

func (m MutlipleChoiceCard) GetType() string { return utils.MULTIPLE_CHOICE_CARD }

type FrontBackCard struct {
	ID    string `json:"id,omitempty" firestore:"-"`
	Type  string `json:"type" validate:"required" firestore:"type"`
	Front string `json:"front" validate:"required" firestore:"front"`
	Back  string `json:"back" validate:"required" firestore:"back"`
}

func (f FrontBackCard) GetType() string { return utils.FRONT_BACK_CARD }

type OrderedCard struct {
	ID      string   `json:"id,omitempty" firestore:"-"`
	Type    string   `json:"type" validate:"required"  firestore:"type"`
	Options []string `json:"options" validate:"required" firestore:"options"`
}

func (o OrderedCard) GetType() string { return utils.ORDERED_CARD }

type BlanksCard struct {
	ID       string   `json:"id,omitempty" firestore:"-"`
	Type     string   `json:"type" validate:"required" firestore:"type"`
	Question string   `json:"question" validate:"required" firestore:"question"`
	Answers  []string `json:"answers" validate:"required" firestore:"answers"`
}

func (b BlanksCard) GetType() string { return utils.BLANKS_CARD }

type CardType struct {
	Type string `json:"type"`
}
