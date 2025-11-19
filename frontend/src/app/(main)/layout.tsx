'use client';

import { usePathname, useRouter } from 'next/navigation';
import { AppSidebar } from '@/components/app-sidebar';
import { ModeToggle } from '@/components/theme-toggle';
import { Separator } from '@/components/ui/separator';
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useEffect, useState } from 'react';
import { cn } from '@/lib/utils';

export default function RootLayout({
  breadcrumb,
  children,
}: Readonly<{
  breadcrumb: React.ReactNode;
  children: React.ReactNode;
}>) {
  const pathname = usePathname();
  const router = useRouter();
  const [defaultOpen, setDefaultOpen] = useState<boolean | null>(null);

  useEffect(() => {
    const sidebarState = document.cookie
      .split('; ')
      .find((row) => row.startsWith('sidebar_state='))
      ?.split('=')[1];
    setDefaultOpen(sidebarState !== 'false');
  }, []);

  // Don't render until we've read the cookie
  if (defaultOpen === null) {
    return null;
  }

  const isRoot = pathname === '/';

  const handleBack = () => {
    router.back();
  };

  return (
    <SidebarProvider defaultOpen={defaultOpen}>
      <AppSidebar />
      <SidebarInset className="border-border flex flex-col border px-4">
        <header className="flex h-16 shrink-0 items-center gap-2">
          <div className="-ml-1 flex items-center gap-1">
            <SidebarTrigger />
            {!isRoot && (
              <Button
                variant="ghost"
                size="icon"
                className={cn('size-7', 'animate-in fade-in slide-in-from-left-2 rounded-full duration-200')}
                onClick={handleBack}
              >
                <ArrowLeft />
                <span className="sr-only">Back</span>
              </Button>
            )}
            <Separator orientation="vertical" className="mr-2 data-[orientation=vertical]:h-4" />
            {breadcrumb}
          </div>
          <div className="ml-auto">
            <ModeToggle />
          </div>
        </header>
        {children}
      </SidebarInset>
    </SidebarProvider>
  );
}
