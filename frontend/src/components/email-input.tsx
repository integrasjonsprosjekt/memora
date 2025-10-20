import * as React from "react"

import { useState, useRef } from "react"
import { X } from "lucide-react"

interface EmailInputProps {
  value?: string[];
  onChange?: (emails: string[]) => void;
}

function EmailInput({value = [], onChange}: EmailInputProps) {
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
      const trimmed = inputValue.trim()
      if (trimmed && isValidEmail(trimmed) && !value.includes(trimmed)) {
        onChange?.([trimmed, ...value])
        setInputValue("")
      }
      e.preventDefault()
    } else if (e.key === "Backspace" && !inputValue && value.length) {
      // Remove last email if input is empty
      onChange?.(value.slice(1))
    }
  }

  function removeEmail(idx: number) {
    const newEmails = value.filter((_, i) => i !== idx)
    onChange?.(newEmails)
    // Refocus input after removing
    setTimeout(() => {
      inputRef.current?.focus()
    }, 0)
  }

  return (
    <div
      className="w-full min-h-10 min-w-0 px-3 py-1
        flex flex-wrap items-center gap-2
        border-input border rounded-md shadow-xs
        dark:bg-input/30
        focus-within:border-ring focus-within:ring-ring/50 focus-within:ring-[3px]"
      >
      {/* Render email chips */}
      {value.map((email, idx) => (
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
        type="email"
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

export { EmailInput }
