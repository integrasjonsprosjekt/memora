import * as React from "react"

import { cn } from "@/lib/utils"

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        // Size, spacing, layout
        "w-full min-w-0 min-h-10 px-3 py-2",
        "flex rounded-md bg-transparent text-base md:text-sm min-w-0",

        // File input styles
        "file:text-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium",

        // Placeholder, selection
        "placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground",

        // Background, shadow, transition
        "dark:bg-input/30 shadow-xs transition-[color,box-shadow]",

        // Border
        "border-input border",

        // Focus/active states
        "outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",

        // Validation states
        "aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",

        // Disabled states
        "disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50",

        className
      )}
      {...props}
    />
  )
}

export { Input }
