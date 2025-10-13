import { CardType } from "@/types/cards";

export default function buildCardPayload(cardType: CardType, formData: Record<string, any>) {
  switch (cardType) {
    case "front_back":
      return {
        front: formData.front,
        back: formData.back,
      };
    case "blanks":
      return {
        question: formData.question,
        answers: formData.answers
          .split(",")
          .map((s: string) => s.trim())
          .filter(Boolean),
      };
    case "multiple_choice": {
      const optionsArray = formData.options
        .split(",")
        .map((s: string) => s.trim())
        .filter(Boolean);
      const options = optionsArray.reduce((acc: Record<string, boolean>, opt: string) => {
        acc[opt] = opt === formData.answer.trim();
        return acc;
      }, {});
      return {
        question: formData.question,
        options,
      };
    }
    case "ordered":
      return {
        question: formData.question,
        options: formData.options
          .split(",")
          .map((s: string) => s.trim())
          .filter(Boolean),
      };
    default:
      throw new Error(`Unknown card type: ${cardType}`);
  }
}
