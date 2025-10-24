'use client';

import { UseFormReturn } from 'react-hook-form';
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { cardInputSchemas } from '@/lib/cardSchemas';
import z from 'zod';

type Ordered = z.infer<typeof cardInputSchemas.ordered>;

type Props = { form: UseFormReturn<Ordered> };

export const OrderedFields = ({ form }: Props) => {
  return (
    <>
      <FormField
        control={form.control}
        name="question"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Question</FormLabel>
            <FormControl>
              <Input placeholder="Question" value={field.value ?? ''} onChange={field.onChange} onBlur={field.onBlur} />
            </FormControl>
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name="options"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Options in order (comma-separated)</FormLabel>
            <FormControl>
              <Input
                placeholder="e.g. A, B, C"
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
