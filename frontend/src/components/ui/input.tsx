import * as React from "react"

import { cn } from "@/lib/utils"
import { useState, useRef } from "react"
import { X } from "lucide-react"

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
function EmailInput() {
  const [emails, setEmails] = useState<string[]>([])
  const [inputValue, setInputValue] = useState("")
  const inputRef = useRef<HTMLInputElement>(null)

  function isValidEmail(email: string) {
    // Simple email validation
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
  }

  function handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    setInputValue(e.target.value)
  }

  function handleInputKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === " " || e.key === "Enter") {
      const value = inputValue.trim()
      if (value && isValidEmail(value) && !emails.includes(value)) {
        setEmails([value, ...emails])
        setInputValue("")
      }
      e.preventDefault()
    } else if (e.key === "Backspace" && !inputValue && emails.length) {
      // Remove last email if input is empty
      setEmails(emails.slice(1))
    }
  }

  function removeEmail(idx: number) {
    setEmails(emails.filter((_, i) => i !== idx))
    // Refocus input after removing
    setTimeout(() => {
      inputRef.current?.focus()
    }, 0)
  }

  return (
    <div
      className={cn(
        "w-full min-h-10 min-w-0 px-3 py-1",
        "flex flex-wrap items-center gap-2",
        "border-input border rounded-md shadow-xs",
        "dark:bg-input/30",
        "focus-within:border-ring focus-within:ring-ring/50 focus-within:ring-[3px]",
      )}
    >
      {/* Render email chips */}
      {emails.map((email, idx) => (
        <span
          key={email}
          className="flex items-center border bg-muted px-2 py-1 rounded text-sm mr-1"
        >
          {email}
          <button
            type="button"
            className="ml-2 text-destructive hover:text-destructive/80 rounded focus:outline-2 focus:outline-ring"
            onClick={() => removeEmail(idx)}
            aria-label={`Remove ${email}`}
          >
            <X className="h-4 w-4" />
          </button>
        </span>
      ))}
      {/* Input field */}
      <input
        ref={inputRef}
        type="text"
        value={inputValue}
        onChange={handleInputChange}
        onKeyDown={handleInputKeyDown}
        className="flex-1 min-w-[120px] placeholder:text-muted-foreground text-base bg-transparent outline-none border-none"
        placeholder="Add email and press space"
        aria-label="Add email"
      />
    </div>
  )
}

export { Input, EmailInput }
