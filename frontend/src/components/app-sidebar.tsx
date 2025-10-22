'use client';

import * as React from 'react';
import { CalendarClock, PieChart, CircleQuestionMark, BookText } from 'lucide-react';

import Image from 'next/image';
import Link from 'next/link';
import { NavMain } from '@/components/nav-main';
import { NavActions } from '@/components/nav-actions';
import { NavSecondary } from '@/components/nav-secondary';
import { NavUser } from '@/components/nav-user';
import { OiiaiTrigger } from '@/components/oiiai';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar';

const data = {
  user: {
    name: 'user',
    email: 'm@example.com',
    avatar: '/logo.svg',
  },
  actions: [
    {
      name: 'Dashboard',
      url: '/dashboard',
      icon: PieChart,
    },
    {
      name: 'Today',
      url: '/today',
      icon: CalendarClock,
    },
  ],
  footer: [
    {
      title: 'Help',
      url: 'https://github.com/integrasjonsprosjekt/memora/issues',
      icon: CircleQuestionMark,
    },
    {
      title: 'Documentation',
      url: 'https://github.com/integrasjonsprosjekt/memora/wiki',
      icon: BookText,
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar variant="inset" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link href="/">
                <OiiaiTrigger targetId="logo">
                  <div id="logo" className="flex cursor-pointer flex-row">
                    <div className="bg-sidebar-accent border-border text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg border-1 text-xs wrap-anywhere">
                      <Image src="/logo.svg" width={250} height={250} alt="Memora" />
                    </div>
                    <div className="grid flex-1 pl-2 text-left text-sm leading-tight">
                      <span className="truncate font-bold">Memora</span>
                      <span className="truncate text-xs">Flashcards</span>
                    </div>
                  </div>
                </OiiaiTrigger>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavActions projects={data.actions} />
        <NavMain />
        <NavSecondary items={data.footer} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
    </Sidebar>
  );
}
