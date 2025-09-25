import {markdownToHtml} from "@/components/Markdown/markdown";
import styles from "./markdown.module.css";

type MarkdownProps = {
  content: string;
};
export default async function MarkdownRenderer({content}: MarkdownProps) {
  const html = await markdownToHtml(content);
  return <div className={styles.markdown} dangerouslySetInnerHTML={{ __html: html }} />;
}
