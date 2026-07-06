import { useState, useEffect } from "react";
import { PillSelect } from "./PillSelect";
import { useDebouncedInputValue } from "../lib/useDebouncedInputValue";

interface CategoryPillsProps {
  initialCategories?: string[];
  fadeColor?: string;
}

export function CategoryPills({
  initialCategories,
  fadeColor,
}: CategoryPillsProps) {
  const [categories, setCategories] = useState<string[]>(
    initialCategories || [],
  );
  const query = useDebouncedInputValue("cat-input", 150);

  useEffect(() => {
    if (initialCategories && initialCategories.length > 0) return;

    fetch("/api/categories")
      .then((r) => r.json())
      .then((json) => setCategories(json.data || []))
      .catch(() => {});
  }, [initialCategories]);

  const handleSelect = (option: { label: string; value: string }) => {
    const input = document.getElementById(
      "cat-input",
    ) as HTMLInputElement | null;
    if (input) {
      input.value = option.value;
      input.dispatchEvent(new Event("change", { bubbles: true }));
    }
  };

  const filtered = query
    ? categories.filter((cat) => cat.toLowerCase().includes(query))
    : categories;

  const options = filtered.map((cat) => ({ label: cat, value: cat }));

  return (
    <PillSelect
      options={options}
      onSelect={handleSelect}
      fadeColor={fadeColor}
      className="mt-2"
    />
  );
}
