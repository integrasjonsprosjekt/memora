'use client';

import { Brain, WifiSync, Sparkles } from 'lucide-react';
import './word-carousel.css';
import { MarkdownRenderer } from '@/components/markdown';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { useEffect, useMemo } from 'react';

const ROTATING_WORDS = ['languages', 'subjects', 'vocabulary', 'history', 'science', 'anything'];

// Generate keyframes dynamically based on word count
function generateKeyframes(wordCount: number, animationName: string) {
  const showDuration = 85; // Percentage of time showing each word
  const stepPercentage = 100 / wordCount;

  let keyframes = `@keyframes ${animationName} {\n  0% { top: 0; }\n`;

  for (let i = 0; i < wordCount; i++) {
    const showUntil = i * stepPercentage + (stepPercentage * showDuration) / 100;
    const transitionTo = (i + 1) * stepPercentage;
    const position = -1.5 * i;

    keyframes += `  ${showUntil.toFixed(2)}% { top: ${position}em; }\n`;
    keyframes += `  ${transitionTo.toFixed(2)}% { top: ${-1.5 * (i + 1)}em; }\n`;
  }

  keyframes += `}`;
  return keyframes;
}

export default function Page() {
  const animationName = useMemo(() => `word-slide-${ROTATING_WORDS.length}`, []);
  const animationDuration = ROTATING_WORDS.length * 2; // 2 seconds per word

  useEffect(() => {
    // Inject dynamic keyframes
    const styleId = 'word-carousel-animation';
    let styleElement = document.getElementById(styleId) as HTMLStyleElement;

    if (!styleElement) {
      styleElement = document.createElement('style');
      styleElement.id = styleId;
      document.head.appendChild(styleElement);
    }

    styleElement.textContent = generateKeyframes(ROTATING_WORDS.length, animationName);
  }, [animationName]);
  return (
    <div className="flex flex-1 flex-col p-4 pt-0">
      <div className="flex flex-1 flex-col gap-8">
        {/* Hero Section */}
        <div className="from-primary/10 via-primary/5 to-background rounded-xl bg-gradient-to-br p-8 md:p-12">
          <h1 className="from-primary to-primary/60 mb-4 bg-gradient-to-r bg-clip-text text-4xl font-bold text-transparent md:text-5xl">
            Welcome to Memora
          </h1>
          <p className="text-muted-foreground mb-6 text-lg md:text-xl">
            A flashcard app that uses spaced repetition to help you memorize{' '}
            <span className="word-carousel-container">
              <span
                className="word-carousel"
                style={{
                  animation: `${animationName} ${animationDuration}s ease infinite`,
                }}
              >
                {ROTATING_WORDS.map((word, index) => (
                  <span key={index} className="word-carousel-item">
                    {word}
                  </span>
                ))}
                {/* Duplicate first word for seamless loop */}
                <span className="word-carousel-item">{ROTATING_WORDS[0]}</span>
              </span>
            </span>
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid gap-4 md:grid-cols-3">
          <div className="bg-muted/50 hover:bg-muted/70 rounded-xl p-6 transition-colors">
            <div className="bg-primary/10 mb-4 flex h-12 w-12 items-center justify-center rounded-lg">
              <Brain className="text-primary h-6 w-6" />
            </div>
            <h3 className="mb-2 text-xl font-semibold">Spaced Repetition</h3>
            <p className="text-muted-foreground">
              Cards show up when you&apos;re about to forget them, helping you remember long-term.
            </p>
          </div>

          <div className="bg-muted/50 hover:bg-muted/70 rounded-xl p-6 transition-colors">
            <div className="bg-primary/10 mb-4 flex h-12 w-12 items-center justify-center rounded-lg">
              <WifiSync className="text-primary h-6 w-6" />
            </div>
            <h3 className="mb-2 text-xl font-semibold">Always Synchronized</h3>
            <p className="text-muted-foreground">Your data is automatically synchronized across all your devices.</p>
          </div>

          <div className="bg-muted/50 hover:bg-muted/70 rounded-xl p-6 transition-colors">
            <div className="bg-primary/10 mb-4 flex h-12 w-12 items-center justify-center rounded-lg">
              <Sparkles className="text-primary h-6 w-6" />
            </div>
            <h3 className="mb-2 text-xl font-semibold">Bring Your Own Content</h3>
            <p className="text-muted-foreground">
              Create your own cards with full Markdown and{' '}
              <Tooltip>
                <TooltipTrigger>
                  <MarkdownRenderer skeleton="LaTeX" inline>
                    {'$\\LaTeX{ }$'}
                  </MarkdownRenderer>
                </TooltipTrigger>
                <TooltipContent>
                  Uses{' '}
                  <MarkdownRenderer skeleton="KaTeX" inline>
                    {'$\\KaTeX{ }$'}
                  </MarkdownRenderer>{' '}
                  behind the scenes 😉
                </TooltipContent>
              </Tooltip>{' '}
              support, in a simple, distraction-free interface. <i>Import from other providers coming soon.</i>
            </p>
          </div>
        </div>
      </div>

      {/* Call to Action */}
      <div className="from-primary to-primary/80 text-primary-foreground mt-8 rounded-xl bg-gradient-to-r p-8 text-center">
        <h2 className="mb-3 text-2xl font-bold md:text-3xl">Get Started</h2>
        <p className="text-primary-foreground/90 mx-auto mb-6 max-w-2xl">Create your first deck to start learning.</p>
        <div className="text-primary-foreground/80 text-sm">
          Go to <span className="font-semibold">Decks</span> in the sidebar
        </div>
      </div>
    </div>
  );
}
