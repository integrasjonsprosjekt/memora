'use client';

import { ReactNode, useEffect, useState } from 'react';
import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import remarkRehype from 'remark-rehype';
import rehypeSanitize, { defaultSchema } from 'rehype-sanitize';
import rehypeStringify from 'rehype-stringify';
import rehypeHighlight from 'rehype-highlight';
import rehypeKatex from 'rehype-katex';
import styles from './markdown.module.css';
import { Skeleton } from '@/components/ui/skeleton';
import { cn } from '@/lib/utils';

export type MarkdownProps = {
  children: string;
  className?: string;
  skeleton?: ReactNode;
  inline?: boolean;
};

export function MarkdownRenderer({ children, className, skeleton, inline = false }: MarkdownProps) {
  const [html, setHtml] = useState<string>('');
  const [isLoading, setIsLoading] = useState(true);
  const classNames = cn(styles.markdown, inline ? styles.inline : '', className);

  useEffect(() => {
    const processMarkdown = async () => {
      try {
        // Create a custom sanitize schema that allows KaTeX elements
        const katexSchema = {
          ...defaultSchema,
          attributes: {
            ...defaultSchema.attributes,
            '*': [...(defaultSchema.attributes?.['*'] || []), 'className'],
            span: [...(defaultSchema.attributes?.span || []), 'className', 'style'],
            div: [...(defaultSchema.attributes?.div || []), 'className', 'style'],
          },
          tagNames: [
            ...(defaultSchema.tagNames || []),
            'math',
            'semantics',
            'mrow',
            'mi',
            'mo',
            'mn',
            'msup',
            'msub',
            'mfrac',
            'mtext',
            'annotation',
          ],
        };

        const file = await unified()
          .use(remarkParse)
          .use(remarkGfm)
          .use(remarkMath)
          .use(remarkRehype)
          .use(rehypeKatex)
          .use(rehypeSanitize, katexSchema)
          .use(rehypeHighlight)
          .use(rehypeStringify)
          .process(children);

        setHtml(String(file));
      } catch (error) {
        console.error('Error processing markdown:', error);
        // Fallback to plain text if markdown processing fails
        setHtml(`<p>${children}</p>`);
      } finally {
        setIsLoading(false);
      }
    };

    processMarkdown();
  }, [children]);

  if (isLoading) {
    if (skeleton) {
      return <div className={classNames}>{skeleton}</div>;
    }

    return (
      <div className={classNames}>
        <Skeleton className="h-4 w-[250px]" />
        <Skeleton className="h-4 w-[200px]" />
      </div>
    );
  }

  return <div className={classNames} dangerouslySetInnerHTML={{ __html: html }} />;
}
