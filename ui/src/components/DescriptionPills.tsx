import { PillSelect } from './PillSelect';

const MOCK_DESCRIPTIONS = [
  'Lunch',
  'Dinner',
  'Groceries',
  'Coffee',
  'Gas',
  'Transport',
  'Utilities',
  'Shopping',
];

interface DescriptionPillsProps {
  fadeColor?: string;
}

export function DescriptionPills({ fadeColor }: DescriptionPillsProps) {
  const handleSelect = (option: { label: string; value: string }) => {
    const input = document.getElementById(
      'description-input',
    ) as HTMLInputElement | null;
    if (input) input.value = option.value;
  };

  const options = MOCK_DESCRIPTIONS.map((d) => ({ label: d, value: d }));

  return (
    <PillSelect
      options={options}
      onSelect={handleSelect}
      fadeColor={fadeColor}
      className="mt-2"
    />
  );
}
