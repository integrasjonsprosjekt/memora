"use client";

import * as React from "react";
import { CalendarClock, Boxes, PieChart, Send } from "lucide-react";

import Image from "next/image";
import Link from "next/link";
import { NavMain } from "@/components/nav-main";
import { NavActions } from "@/components/nav-actions";
import { NavSecondary } from "@/components/nav-secondary";
import { NavUser } from "@/components/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";

const data = {
  user: {
    name: "user",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  actions: [
    {
      name: "Dashboard",
      url: "#",
      icon: PieChart,
    },
    {
      name: "Today",
      url: "#",
      icon: CalendarClock,
    },
  ],
  decks: [
    {
      title: "Objektorientert Programmering",
      url: "#",
      isActive: true,
      icon: Boxes,
      items: [
        {
          title: "Overview",
          url: "#",
        },
        {
          title: "Today",
          url: "#",
        },
      ],
    },
  ],
  footer: [
    {
      title: "Support",
      url: "https://github.com/integrasjonsprosjekt/memora/issues",
      icon: Send,
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
                    <div className="bg-sidebar-accent text-xs border-border border-1 wrap-anywhere text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                      <Image
                        src="/logo.svg"
                        width={250}
                        height={250}
                        alt="Memora"
                      />
                    </div>
                    <div className="grid flex-1 text-left text-sm leading-tight pl-2">
                      <span className="truncate font-bold">Memora</span>
                      <span className="truncate text-xs">Flashcards</span>
                    </div>
                </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavActions projects={data.actions} />
        <NavMain items={data.decks} />
        <NavSecondary items={data.footer} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
    </Sidebar>
  );
}
