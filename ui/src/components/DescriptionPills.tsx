import { useState, useEffect, useCallback, useRef } from 'react';
import { PillSelect } from './PillSelect';
import { useDebouncedInputValue } from '../lib/useDebouncedInputValue';

interface DescriptionCount {
  description: string;
  count: number;
}

interface DescriptionPillsProps {
  initialCategory?: string;
  fadeColor?: string;
}

export function DescriptionPills({ initialCategory, fadeColor }: DescriptionPillsProps) {
  const [descriptions, setDescriptions] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<number | null>(null);
  const query = useDebouncedInputValue('description-input', 150);

  const fetchDescriptions = useCallback((cat: string) => {
    if (!cat) {
      setDescriptions([]);
      return;
    }
    setLoading(true);
    fetch(`/api/expenses/descriptions?category=${encodeURIComponent(cat)}`)
      .then((r) => {
        if (!r.ok) throw new Error();
        return r.json();
      })
      .then((json) => {
        const data: DescriptionCount[] = json.data || [];
        setDescriptions(data.map((d) => d.description));
      })
      .catch(() => setDescriptions([]))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (initialCategory) {
      fetchDescriptions(initialCategory);
    }
  }, [initialCategory, fetchDescriptions]);

  useEffect(() => {
    const input = document.getElementById('cat-input') as HTMLInputElement | null;
    if (!input) return;

    const handleChange = () => {
      fetchDescriptions(input.value.trim());
    };

    const handleInput = () => {
      if (debounceRef.current) {
        window.clearTimeout(debounceRef.current);
      }
      debounceRef.current = window.setTimeout(() => {
        fetchDescriptions(input.value.trim());
      }, 500);
    };

    input.addEventListener('change', handleChange);
    input.addEventListener('input', handleInput);
    return () => {
      if (debounceRef.current) {
        window.clearTimeout(debounceRef.current);
      }
      input.removeEventListener('change', handleChange);
      input.removeEventListener('input', handleInput);
    };
  }, [fetchDescriptions]);

  const handleSelect = (option: { label: string; value: string }) => {
    const input = document.getElementById(
      'description-input',
    ) as HTMLInputElement | null;
    if (input) input.value = option.value;
  };

  const filtered = query
    ? descriptions.filter((d) => d.toLowerCase().includes(query))
    : descriptions;

  const options = filtered.map((d) => ({ label: d, value: d }));

  return (
    <PillSelect
      options={options}
      onSelect={handleSelect}
      fadeColor={fadeColor}
      className="mt-2"
      loading={loading}
    />
  );
}
