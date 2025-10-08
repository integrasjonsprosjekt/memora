import { z } from "zod";
//import { CardType } from '@/types/cards';

export const cardSchemas = {
  front_back: z.object({
    front: z.string().min(1, "Front requierd"),
    back: z.string().min(1, "Front requierd"),
  }),
  blanks: z.object({
    question: z.string().min(1, "Question required"),
    answer: z.string().min(1, "Answer required"),
  }),
  multiple_choice: z.object({
    question: z.string().min(1, "Question required"),
    options: z.string().min(1, "Choice required"),
    answer: z.string().min(1, "Answer required"),
  }),
  ordered: z.object({
    question: z.string().min(1, "Question required"),
    options: z.string().min(1, "Choice required"),
  }),
};
export type CardSchema = typeof cardSchemas; //<T extends CardType> = z.infer<(typeof cardSchemas)[T]>;
