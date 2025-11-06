import * as React from 'react';

import { useState, useRef } from 'react';
import { X } from 'lucide-react';

interface EmailInputProps {
  value?: string[];
  onChange?: (emails: string[]) => void;
}

function EmailInput({ value = [], onChange }: EmailInputProps) {
  const [inputValue, setInputValue] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  function isValidEmail(email: string) {
    // Simple email validation
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  }

  function handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    setInputValue(e.target.value);
  }

  function handleInputKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === ' ' || e.key === 'Enter') {
      const trimmed = inputValue.trim();
      if (trimmed && isValidEmail(trimmed) && !value.includes(trimmed)) {
        onChange?.([...value, trimmed]);
        setInputValue('');
      }
      e.preventDefault();
    } else if (e.key === 'Backspace' && !inputValue && value.length) {
      // Remove last email if input is empty
      onChange?.(value.slice(0, -1));
    }
  }

  function removeEmail(idx: number) {
    const newEmails = value.filter((_, i) => i !== idx);
    onChange?.(newEmails);
    // Refocus input after removing
    setTimeout(() => {
      inputRef.current?.focus();
    }, 0);
  }

  return (
    <div className="border-input dark:bg-input/30 focus-within:border-ring focus-within:ring-ring/50 flex min-h-10 w-full min-w-0 flex-wrap items-center gap-2 rounded-md border px-3 py-1 shadow-xs focus-within:ring-[3px]">
      {/* Render email chips */}
      {value.map((email, idx) => (
        <span key={email} className="bg-muted mr-1 flex items-center rounded border px-2 py-1 text-sm">
          {email}
          <button
            type="button"
            className="text-destructive hover:text-destructive/80 focus:outline-ring ml-2 rounded focus:outline-2"
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
        className="placeholder:text-muted-foreground min-w-[120px] flex-1 border-none bg-transparent text-base outline-none"
        placeholder="Add email and press space"
        aria-label="Add email"
      />
    </div>
  );
}

export { EmailInput };
