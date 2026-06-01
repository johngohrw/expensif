import { useState, useEffect } from "react";
import { PillSelect } from "./PillSelect";

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

  const options = categories.map((cat) => ({ label: cat, value: cat }));

  return (
    <PillSelect
      options={options}
      onSelect={handleSelect}
      fadeColor={fadeColor}
      className="mt-2"
    />
  );
}
