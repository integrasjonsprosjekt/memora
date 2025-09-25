import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkRehype from "remark-rehype";
import rehypeSanitize from "rehype-sanitize";
import rehypeStringify from "rehype-stringify";
import rehypeHighlight from "rehype-highlight";
import styles from "./markdown.module.css";
import "highlight.js/styles/github-dark.css";

type MarkdownProps = {
  content: string;
};

export async function markdownToHtml(markdown: string): Promise<string> {
  const file = await unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkRehype)
    .use(rehypeSanitize)
    .use(rehypeHighlight)
    .use(rehypeStringify)
    .process(markdown);

  return String(file);
}

export default async function MarkdownRenderer({content}: MarkdownProps) {
  const html = await markdownToHtml(content);
  return <div className={styles.markdown} dangerouslySetInnerHTML={{ __html: html }} />;
}
