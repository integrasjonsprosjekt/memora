'use client';

import { UseFormReturn } from 'react-hook-form';
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import z from 'zod';
import { cardInputSchemas } from '@/lib/cardSchemas';

type FrontBack = z.infer<typeof cardInputSchemas.front_back>;

type Props = { form: UseFormReturn<FrontBack> };

export const FrontBackField = ({ form }: Props) => {
  return (
    <>
      <FormField
        control={form.control}
        name="front"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Front</FormLabel>
            <FormControl>
              <Input
                placeholder="Front text"
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
        name="back"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Back</FormLabel>
            <FormControl>
              <Input
                placeholder="Back text"
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
