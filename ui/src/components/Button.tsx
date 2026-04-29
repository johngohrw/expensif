import { type ButtonHTMLAttributes } from 'react';

// Keep in sync with templates/partials/button.html — both must use identical
// Tailwind class mappings for variant + size.
type ButtonVariant = 'primary' | 'secondary' | 'neutral' | 'ghost' | 'danger' | 'pill';
type ButtonSize = 'md' | 'sm' | 'xs';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
}

const variantClasses: Record<ButtonVariant, string> = {
  primary: 'bg-blue-600 text-white hover:bg-blue-700 rounded-lg',
  secondary: 'bg-blue-500 text-white hover:bg-blue-600 rounded-lg',
  neutral: 'bg-gray-100 text-gray-700 hover:bg-gray-200 rounded',
  ghost: 'bg-gray-200 text-gray-700 hover:bg-gray-300 rounded-lg',
  danger: 'bg-red-50 text-red-600 hover:bg-red-100 rounded',
  pill: 'bg-blue-50 text-blue-700 hover:bg-blue-100 rounded-full',
};

const sizeClasses: Record<ButtonSize, string> = {
  md: 'py-2.5',
  sm: 'text-xs px-3 py-1.5',
  xs: 'text-xs px-3 py-1',
};

export function Button({
  variant = 'primary',
  size = 'md',
  className = '',
  type = 'button',
  children,
  ...props
}: ButtonProps) {
  const classes = [
    'font-medium transition focus:outline-none',
    variantClasses[variant],
    sizeClasses[size],
    className,
  ]
    .join(' ')
    .trim();

  return (
    <button type={type} className={classes} {...props}>
      {children}
    </button>
  );
}
