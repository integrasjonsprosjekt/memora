'use client';

import { UseFormReturn } from 'react-hook-form';
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { cardInputSchemas } from '@/lib/cardSchemas';
import z from 'zod';

type MultipleChoice = z.infer<typeof cardInputSchemas.multiple_choice>;

type Props = { form: UseFormReturn<MultipleChoice> };

export const MultipleChoiceFields = ({ form }: Props) => {
  return (
    <>
      <FormField
        control={form.control}
        name="question"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Question</FormLabel>
            <FormControl>
              <Input
                placeholder="Enter the question"
                value={field.value ?? ''}
                onChange={field.onChange}
                onBlur={field.onBlur}
              />
            </FormControl>
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name="options"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Options (comma-separated)</FormLabel>
            <FormControl>
              <Input
                placeholder="e.g. A, B, C, D"
                value={field.value ?? ''}
                onChange={field.onChange}
                onBlur={field.onBlur}
              />
            </FormControl>
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name="answer"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Correct answer</FormLabel>
            <FormControl>
              <Input placeholder="e.g. B" value={field.value ?? ''} onChange={field.onChange} onBlur={field.onBlur} />
            </FormControl>
          </FormItem>
        )}
      />
    </>
  );
};
