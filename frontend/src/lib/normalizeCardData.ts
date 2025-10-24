import { Card } from '@/types/card';

export default function normalizeCardData(card: Card) {
  switch (card.type) {
    case 'multiple_choice': {
      const options = Object.keys(card.options || {}).join(', ');
      const answer = Object.entries(card.options || {})
        .filter(([_, v]) => v === true)
        .map(([key]) => key)
        .join(', ');
      return {
        question: card.question ?? '',
        options,
        answer,
      };
    }

    case 'front_back':
      return {
        front: card.front ?? '',
        back: card.back ?? '',
      };

    case 'blanks':
      return {
        question: card.question ?? '',
        answers: Array.isArray(card.answers) ? card.answers.join(', ') : '',
      };

    case 'ordered':
      return {
        question: card.question ?? '',
        options: Array.isArray(card.options) ? card.options.join(', ') : '',
      };

    default:
      return {};
  }
}
