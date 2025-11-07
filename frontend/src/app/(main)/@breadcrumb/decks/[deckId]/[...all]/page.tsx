import { BreadcrumbSeparator } from '@/components/ui/breadcrumb';
import { BreadcrumbPath } from '@/lib/breadcrumbs';

export default async function Page(props: { params: Promise<{ deckId: string; all: string[] }> }) {
  const { deckId, all } = await props.params;

  // Shift through the array until we find deckId, then shift once more
  const segments = [...all];
  while (segments.length > 0 && segments[0] !== deckId) {
    segments.shift();
  }
  // Remove the deckId itself
  segments.shift();

  return (
    <>
      <BreadcrumbSeparator />
      <BreadcrumbPath segments={segments} />
    </>
  );
}
