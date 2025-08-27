"use client";

import * as React from "react";
import "./oiiai.scss";

interface OiiaiTriggerProps {
  children: React.ReactNode;
  targetId: string;
  audioUrl?: string;
  holdDuration?: number;
  animationDuration?: number;
}

export function OiiaiTrigger({
  children,
  targetId: targetId,
  audioUrl = "/oiiai.wav",
  holdDuration = 1500,
  animationDuration = 6800,
}: OiiaiTriggerProps) {
  const holdTimerRef = React.useRef<NodeJS.Timeout | null>(null);
  const animationTriggeredRef = React.useRef(false);

  const triggerAnimation = React.useCallback(() => {
    const audio = new Audio(audioUrl);
    const targetElement = document.getElementById(targetId);

    if (!targetElement) return;

    // Play audio and animation
    audio.play().catch(console.error);

    targetElement.classList.add("oiiai");

    // Remove the class after the animation completes
    setTimeout(() => {
      targetElement.classList.remove("oiiai");
    }, animationDuration);
  }, [targetId, audioUrl, animationDuration]);

  const handleMouseDown = React.useCallback(() => {
    animationTriggeredRef.current = false;

    // Set a timer for the hold duration
    holdTimerRef.current = setTimeout(() => {
      animationTriggeredRef.current = true;
      triggerAnimation();
    }, holdDuration);
  }, [holdDuration, triggerAnimation]);

  const handleMouseUp = React.useCallback(() => {
    // Clear the timer if mouse is released before hold duration
    if (holdTimerRef.current) {
      clearTimeout(holdTimerRef.current);
      holdTimerRef.current = null;
    }
  }, []);

  const handleMouseLeave = handleMouseUp;

  const handleClick = React.useCallback((e: React.MouseEvent) => {
    // Prevent default action if the animation was triggered by holding
    if (animationTriggeredRef.current) {
      e.preventDefault();
      animationTriggeredRef.current = false;
    }
  }, []);

  // Cleanup timer on unmount
  React.useEffect(() => {
    return () => {
      if (holdTimerRef.current) {
        clearTimeout(holdTimerRef.current);
      }
    };
  }, []);

  return (
    <div
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseLeave}
      onClick={handleClick}
    >
      {children}
    </div>
  );
}
