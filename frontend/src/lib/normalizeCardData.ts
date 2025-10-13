import { Card as CardType } from '@/types/cards';

export default function normalizeCardData(card: CardType) {
  switch (card.type) {
    case "multiple_choice": {
      const options = Object.keys(card.options || {}).join(", ");
      const answer = Object.entries(card.options || {})
        .find(([_, v]) => v === true)?.[0] ?? "";
      return {
        question: card.question ?? "",
        options,
        answer,
      };
    }

    case "front_back":
      return {
        front: card.front ?? "",
        back: card.back ?? "",
      };

    case "blanks":
      return {
        question: card.question ?? "",
        answers: Array.isArray(card.answers) ? card.answers.join(", ") : "",
      };

    case "ordered":
      return {
        question: card.question ?? "",
        options: Array.isArray(card.options) ? card.options.join(", ") : "",
      };

    default:
      return {};
  }
}
