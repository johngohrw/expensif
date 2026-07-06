import { useState, useEffect, useRef } from "react";

// Tracks an input element's value by id, debounced: the returned value only
// updates `delay` ms after the user stops typing, and is trimmed/lowercased
// for case-insensitive substring matching.
export function useDebouncedInputValue(elementId: string, delay: number): string {
  const [value, setValue] = useState("");
  const timerRef = useRef<number | null>(null);

  useEffect(() => {
    const input = document.getElementById(elementId) as HTMLInputElement | null;
    if (!input) return;

    const handleInput = () => {
      if (timerRef.current) window.clearTimeout(timerRef.current);
      const current = input.value;
      timerRef.current = window.setTimeout(() => {
        setValue(current.trim().toLowerCase());
      }, delay);
    };

    input.addEventListener("input", handleInput);
    return () => {
      if (timerRef.current) window.clearTimeout(timerRef.current);
      input.removeEventListener("input", handleInput);
    };
  }, [elementId, delay]);

  return value;
}
