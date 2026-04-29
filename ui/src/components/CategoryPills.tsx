import { useState, useEffect } from "react";
import { Button } from "./Button";

interface CategoryPillsProps {
  initialCategories?: string[];
}

export function CategoryPills({ initialCategories }: CategoryPillsProps) {
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

  const setCategory = (cat: string) => {
    const input = document.getElementById(
      "cat-input",
    ) as HTMLInputElement | null;
    if (input) input.value = cat;
  };

  if (categories.length === 0) return null;

  return (
    <div className="flex flex-wrap gap-2 mt-2">
      {categories.map((cat) => (
        <Button
          key={cat}
          variant="pill"
          size="xs"
          onClick={() => setCategory(cat)}
        >
          {cat}
        </Button>
      ))}
    </div>
  );
}
