'use client';

import { UseFormReturn } from 'react-hook-form';
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form';
import { Textarea } from '@/components/ui/textarea';
import z from 'zod';
import { cardInputSchemas } from '@/lib/cardSchemas';
import { JSX } from 'react';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

type FrontBack = z.infer<typeof cardInputSchemas.front_back>;

type Props = { form: UseFormReturn<FrontBack> };

const MarkdownLogo = (): JSX.Element => {
  return (
    <div className="absolute top-2 right-2">
      <Tooltip>
        <TooltipTrigger asChild>
          <div
            className="bg-muted-foreground h-4 w-4"
            style={{
              maskImage: 'url(/markdown.svg)',
              maskSize: 'contain',
              maskRepeat: 'no-repeat',
              maskPosition: 'center',
              WebkitMaskImage: 'url(/markdown.svg)',
              WebkitMaskSize: 'contain',
              WebkitMaskRepeat: 'no-repeat',
              WebkitMaskPosition: 'center',
            }}
          />
        </TooltipTrigger>
        <TooltipContent>
          <p>Write content using Markdown</p>
        </TooltipContent>
      </Tooltip>
    </div>
  );
};

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
              <div className="relative">
                <Textarea
                  placeholder="Front text"
                  value={field.value ?? ''}
                  onChange={field.onChange}
                  onBlur={field.onBlur}
                  className="pr-10"
                />
                <MarkdownLogo />
              </div>
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
              <div className="relative">
                <Textarea
                  placeholder="Back text"
                  value={field.value ?? ''}
                  onChange={field.onChange}
                  onBlur={field.onBlur}
                  className="pr-10"
                />
                <MarkdownLogo />
              </div>
            </FormControl>
          </FormItem>
        )}
      />
    </>
  );
};
