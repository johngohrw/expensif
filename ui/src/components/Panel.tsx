import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react';

interface PanelProps {
  isOpen: boolean;
  onClose: () => void;
  position?: 'left' | 'right';
  width?: string;
  children: ReactNode;
}

export function Panel({
  isOpen,
  onClose,
  position = 'left',
  width = '280px',
  children,
}: PanelProps) {
  const [mounted, setMounted] = useState(false);
  const [visible, setVisible] = useState(false);
  const scrollCompensationRef = useRef('');

  const lockScroll = useCallback(() => {
    const scrollbarWidth =
      window.innerWidth - document.documentElement.clientWidth;
    scrollCompensationRef.current = document.body.style.paddingRight;
    document.body.style.overflow = 'hidden';
    if (scrollbarWidth > 0) {
      document.body.style.paddingRight = `${scrollbarWidth}px`;
    }
  }, []);

  const unlockScroll = useCallback(() => {
    document.body.style.overflow = '';
    document.body.style.paddingRight = scrollCompensationRef.current;
  }, []);

  useEffect(() => {
    if (isOpen) {
      setMounted(true);
      // Double rAF ensures the browser paints the initial hidden state
      // before we flip to visible — otherwise React batches the updates
      // and the CSS transition is skipped.
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          setVisible(true);
          lockScroll();
        });
      });
    } else {
      setVisible(false);
      const timer = setTimeout(() => {
        setMounted(false);
        unlockScroll();
      }, 150);
      return () => clearTimeout(timer);
    }
  }, [isOpen, lockScroll, unlockScroll]);

  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    document.addEventListener('keydown', handleKey);
    return () => document.removeEventListener('keydown', handleKey);
  }, [isOpen, onClose]);

  if (!mounted) return null;

  const translateClass =
    position === 'left'
      ? visible
        ? 'translate-x-0'
        : '-translate-x-full'
      : visible
        ? 'translate-x-0'
        : 'translate-x-full';

  return (
    <div className="fixed inset-0 z-50" aria-modal="true" role="dialog">
      {/* Backdrop */}
      <div
        className={`absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity duration-150 ${
          visible ? 'opacity-100' : 'opacity-0'
        }`}
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Panel */}
      <div
        className={`absolute top-0 h-full bg-white shadow-xl transition-transform duration-150 ease-out ${translateClass}`}
        style={{ width, [position]: 0 }}
      >
        {/* Close button */}
        <button
          type="button"
          onClick={onClose}
          className="absolute top-3 right-3 p-2 rounded-lg hover:bg-gray-100 text-gray-500 transition"
          aria-label="Close panel"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>

        <div className="pt-12 px-4 h-full overflow-y-auto">{children}</div>
      </div>
    </div>
  );
}
