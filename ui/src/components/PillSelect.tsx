import { useState, useEffect, useRef, useCallback } from 'react';
import { Button } from './Button';

export interface PillOption {
  label: string;
  value: string;
  disabled?: boolean;
}

interface PillSelectProps {
  options: PillOption[];
  onSelect: (option: PillOption) => void;
  value?: string;
  fadeColor?: string;
  className?: string;
  loading?: boolean;
}

export function PillSelect({
  options,
  onSelect,
  value,
  fadeColor = '#ffffff',
  className = '',
  loading = false,
}: PillSelectProps) {
  const [scrolled, setScrolled] = useState(false);
  const [hasOverflow, setHasOverflow] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  const checkOverflow = useCallback(() => {
    const el = scrollRef.current;
    if (!el) return;
    setHasOverflow(el.scrollWidth > el.clientWidth);
    setScrolled(el.scrollLeft > 0);
  }, []);

  useEffect(() => {
    checkOverflow();
    const el = scrollRef.current;
    if (!el) return;

    const ro = new ResizeObserver(checkOverflow);
    ro.observe(el);
    return () => ro.disconnect();
  }, [options, checkOverflow]);

  const handleScroll = useCallback(() => {
    const el = scrollRef.current;
    if (!el) return;
    setScrolled(el.scrollLeft > 0);
  }, []);

  // Loading and "resolved but empty" both render the same px-3 py-1 text-xs
  // row shape as a real pill, so the row's height never changes across
  // loading -> pills or loading -> no-suggestions transitions.
  if (loading) {
    return (
      <div className={`relative ${className}`}>
        <div className="flex items-center pb-1">
          <span className="text-xs text-gray-400 px-3 py-1 animate-pulse">Loading…</span>
        </div>
      </div>
    );
  }

  if (options.length === 0) {
    return (
      <div className={`relative ${className}`}>
        <div className="flex items-center pb-1">
          <span className="text-xs text-gray-400 px-3 py-1">No suggestions</span>
        </div>
      </div>
    );
  }

  return (
    <div className={`relative ${className}`}>
      <div
        ref={scrollRef}
        onScroll={handleScroll}
        className="flex gap-2 overflow-x-auto pb-1 snap-x [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden"
      >
        {options.map((option) => {
          const isActive = value === option.value;
          const activeClass = isActive
            ? 'bg-blue-100 text-blue-800'
            : '';
          const disabledClass = option.disabled
            ? 'opacity-50 cursor-not-allowed'
            : '';

          return (
            <Button
              key={option.value}
              variant="pill"
              size="xs"
              className={`flex-shrink-0 snap-start ${activeClass} ${disabledClass}`}
              disabled={option.disabled}
              onClick={option.disabled ? undefined : () => onSelect(option)}
            >
              {option.label}
            </Button>
          );
        })}
      </div>

      {/* Right fade */}
      {hasOverflow && (
        <div
          className="absolute right-0 top-0 bottom-0 w-8 pointer-events-none"
          style={{ background: `linear-gradient(to right, transparent, ${fadeColor})` }}
        />
      )}

      {/* Left fade */}
      {scrolled && (
        <div
          className="absolute left-0 top-0 bottom-0 w-8 pointer-events-none"
          style={{ background: `linear-gradient(to left, transparent, ${fadeColor})` }}
        />
      )}
    </div>
  );
}
