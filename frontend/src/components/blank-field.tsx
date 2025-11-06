'use client';

import { UseFormReturn } from 'react-hook-form';
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { cardInputSchemas } from '@/lib/cardSchemas';
import z from 'zod';

type Blanks = z.infer<typeof cardInputSchemas.blanks>;

type Props = { form: UseFormReturn<Blanks> };

export const BlankFields = ({ form }: Props) => {
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
                placeholder="Use '{}' for blank fields"
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
        name="answers"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Answer (comma separated)</FormLabel>
            <FormControl>
              <Input
                placeholder="e.g. apple, banana"
                value={field.value ?? ''}
                onChange={field.onChange}
                onBlur={field.onBlur}
              />
            </FormControl>
          </FormItem>
        )}
      />
    </>
  );
};
