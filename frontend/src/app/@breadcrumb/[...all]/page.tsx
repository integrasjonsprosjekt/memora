import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import React from 'react';
import type { ReactElement } from 'react';

// List of routes that should not be linked
const NON_LINKABLE_ROUTES = new Set(['/decks', '/cards']);

function isValidRoute(path: string): boolean {
  return !NON_LINKABLE_ROUTES.has(path);
}

export default async function BreadcrumbSlot(props: { params: Promise<{ all: string[] }> }) {
  const params = await props.params;
  const breadcrumbItems: ReactElement[] = [];
  let breadcrumbPage: ReactElement = <></>;

  for (let i = 0; i < params.all.length; i++) {
    const route = params.all[i];
    const href = `/${params.all.slice(0, i + 1).join('/')}`;
    const isLast = i === params.all.length - 1;
    const shouldLink = isValidRoute(href);

    if (isLast) {
      breadcrumbPage = (
        <BreadcrumbItem>
          <BreadcrumbPage className="capitalize">{route}</BreadcrumbPage>
        </BreadcrumbItem>
      );
    } else {
      breadcrumbItems.push(
        <React.Fragment key={href}>
          <BreadcrumbItem>
            {shouldLink ? (
              <BreadcrumbLink href={href} className="capitalize">
                {route}
              </BreadcrumbLink>
            ) : (
              <BreadcrumbPage className="capitalize">{route}</BreadcrumbPage>
            )}
          </BreadcrumbItem>
          <BreadcrumbSeparator />
        </React.Fragment>,
      );
    }
  }

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {breadcrumbItems}
        {breadcrumbPage}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
