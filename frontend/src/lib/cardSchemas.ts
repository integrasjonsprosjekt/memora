import { z } from 'zod';

const frontBackSchema = z.object({
  type: z.literal('front_back'),
  front: z.string().min(1, 'Front required'),
  back: z.string().min(1, 'Back required'),
});

const frontBackSchemaPayload = frontBackSchema;

const blanksSchema = z.object({
  type: z.literal('blanks'),
  question: z.string().min(1, 'Question required'),
  answers: z.string().min(1, 'Answer required'),
});

const blanksSchemaPayload = blanksSchema.transform(({ type, question, answers }) => ({
  type,
  question,
  answers: answers
    .split(',')
    .map((s: string) => s.trim())
    .filter(Boolean),
}));

const multipleChoiceSchema = z.object({
  type: z.literal('multiple_choice'),
  question: z.string().min(1, 'Question required'),
  options: z.string().min(1, 'Choice required'),
  answer: z.string().min(1, 'Answer required'),
});

const multipleChoiceSchemaPayload = multipleChoiceSchema.transform(({ type, question, options, answer }) => {
  const optionsArray = options
    .split(',')
    .map((s: string) => s.trim())
    .filter(Boolean);

  const answerArray = answer
    .split(',')
    .map((s: string) => s.trim())
    .filter(Boolean);

  const mapped = optionsArray.reduce((acc: Record<string, boolean>, opt: string) => {
    acc[opt] = answerArray.includes(opt);
    return acc;
  }, {});
  return { type, question, options: mapped };
});

const orderedSchema = z.object({
  type: z.literal('ordered'),
  question: z.string().min(1, 'Question required'),
  options: z.string().min(1, 'Choice required'),
});

const orderedSchemaPayload = orderedSchema.transform(({ type, question, options }) => ({
  type,
  question,
  options: options
    .split(',')
    .map((s: string) => s.trim())
    .filter(Boolean),
}));

export const cardInputSchemas = {
  front_back: frontBackSchema,
  blanks: blanksSchema,
  multiple_choice: multipleChoiceSchema,
  ordered: orderedSchema,
};

export const cardPayloadSchemas = {
  front_back: frontBackSchemaPayload,
  blanks: blanksSchemaPayload,
  multiple_choice: multipleChoiceSchemaPayload,
  ordered: orderedSchemaPayload,
};

// Types
export type CardInput =
  | z.infer<typeof frontBackSchema>
  | z.infer<typeof blanksSchema>
  | z.infer<typeof multipleChoiceSchema>
  | z.infer<typeof orderedSchema>;

export type CardPayload =
  | z.infer<typeof frontBackSchemaPayload>
  | z.infer<typeof blanksSchemaPayload>
  | z.infer<typeof multipleChoiceSchemaPayload>
  | z.infer<typeof orderedSchemaPayload>;
