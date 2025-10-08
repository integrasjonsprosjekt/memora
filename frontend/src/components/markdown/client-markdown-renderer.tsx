'use client';

import { useEffect, useState } from 'react';
import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkGfm from 'remark-gfm';
import remarkRehype from 'remark-rehype';
import rehypeSanitize from 'rehype-sanitize';
import rehypeStringify from 'rehype-stringify';
import rehypeHighlight from 'rehype-highlight';
import styles from './markdown.module.css';

type ClientMarkdownRendererProps = {
  children: string;
};

export function ClientMarkdownRenderer({ children }: ClientMarkdownRendererProps) {
  const [html, setHtml] = useState<string>('');
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const processMarkdown = async () => {
      try {
        const file = await unified()
          .use(remarkParse)
          .use(remarkGfm)
          .use(remarkRehype)
          .use(rehypeSanitize)
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
    return <div className={styles.markdown}>Loading...</div>;
  }

  return <div className={styles.markdown} dangerouslySetInnerHTML={{ __html: html }} />;
}
