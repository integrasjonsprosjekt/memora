import { BreadcrumbItem, BreadcrumbPage, BreadcrumbSeparator } from '@/components/ui/breadcrumb';
import { SquareGanttChart } from 'lucide-react';

export default async function Page() {
  return (
    <>
      <BreadcrumbSeparator />
      <BreadcrumbItem>
        <BreadcrumbPage>
          <SquareGanttChart size={16} />
        </BreadcrumbPage>
      </BreadcrumbItem>
    </>
  );
}
