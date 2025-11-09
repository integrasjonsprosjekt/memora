import { Card, FillBlanksCard, MultipleChoiceCard, OrderedCard } from '@/types/card';
import { UserAnswer } from '@/components/card/types';

/**
 * Validates a user's answer against the correct answer for a card.
 * Returns true if the answer is correct, false otherwise.
 */
export function validateAnswer(card: Card, userAnswer: UserAnswer | null): boolean {
  if (!userAnswer && card.type !== 'front_back') return false;

  switch (card.type) {
    case 'front_back':
      return validateFrontBackAnswer(userAnswer as boolean);
    case 'blanks':
      return validateFillBlanksAnswer(card, userAnswer as string[]);
    case 'multiple_choice':
      return validateMultipleChoiceAnswer(card, userAnswer as string | Record<string, boolean>);
    case 'ordered':
      return validateOrderedAnswer(card, userAnswer as string[]);
    default:
      return false;
  }
}

function validateFrontBackAnswer(_userAnswer: boolean): boolean {
  // Front-back cards are always considered correct (flip is triggered on check)
  // We still accept the userAnswer parameter for consistency

  return true;
}

function validateFillBlanksAnswer(card: FillBlanksCard, userAnswer: string[]): boolean {
  if (!Array.isArray(userAnswer)) return false;
  if (userAnswer.length !== card.answers.length) return false;

  // Check if all blanks are filled correctly (case-insensitive, trimmed)
  return card.answers.every((correctAnswer, index) => {
    const userAns = userAnswer[index]?.trim().toLowerCase();
    const correctAns = correctAnswer.trim().toLowerCase();
    return userAns === correctAns;
  });
}

function validateMultipleChoiceAnswer(card: MultipleChoiceCard, userAnswer: string | Record<string, boolean>): boolean {
  const keys = Object.keys(card.options);
  const correctAnswers = keys.filter((key) => card.options[key]);

  // Single choice (radio button)
  if (correctAnswers.length === 1) {
    if (typeof userAnswer !== 'string') return false;
    return userAnswer === correctAnswers[0];
  }

  // Multiple choice (checkboxes)
  if (typeof userAnswer === 'string' || Array.isArray(userAnswer)) return false;

  // Check that all options match the correct state
  return keys.every((key) => {
    const isCorrect = card.options[key];
    const isSelected = userAnswer[key] || false;
    return isCorrect === isSelected;
  });
}

function validateOrderedAnswer(card: OrderedCard, userAnswer: string[]): boolean {
  if (!Array.isArray(userAnswer)) return false;
  if (userAnswer.length !== card.options.length) return false;

  // Check if the order matches exactly
  return card.options.every((item, index) => item === userAnswer[index]);
}

/**
 * Checks if a user has provided any answer (even if incorrect).
 * Useful for determining if the "Check" button should be enabled.
 */
export function hasAnswer(card: Card, userAnswer: UserAnswer | null): boolean {
  if (!userAnswer) return false;

  switch (card.type) {
    case 'front_back':
      // Front-back cards always have an "answer" (they just need to flip)
      return true;
    case 'blanks':
      if (!Array.isArray(userAnswer)) return false;
      return userAnswer.some((ans) => ans.trim().length > 0);
    case 'multiple_choice':
      if (typeof userAnswer === 'string') {
        return userAnswer.length > 0;
      }
      if (typeof userAnswer === 'object' && !Array.isArray(userAnswer)) {
        return Object.values(userAnswer).some((val) => val === true);
      }
      return false;
    case 'ordered':
      return Array.isArray(userAnswer) && userAnswer.length > 0;
    default:
      return false;
  }
}
