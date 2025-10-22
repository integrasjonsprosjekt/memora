import { BreadcrumbPath } from '@/lib/breadcrumbs';

export default async function Page(props: { params: Promise<{ all: string[] }> }) {
  const params = await props.params;

  return <BreadcrumbPath segments={params.all} />;
}
