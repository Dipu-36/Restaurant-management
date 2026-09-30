import { cn } from '../../lib/utils';

type BadgeVariant =
  | 'neutral'
  | 'success'
  | 'warning'
  | 'danger'
  | 'info'
  | 'primary';

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
  size?: 'sm' | 'md';
  className?: string;
}

const variantStyles: Record<BadgeVariant, string> = {
  neutral: 'bg-muted text-text-secondary',
  success: 'bg-green-100 text-green-800',
  warning: 'bg-amber-100 text-amber-800',
  danger: 'bg-red-100 text-red-700',
  info: 'bg-blue-100 text-blue-800',
  primary: 'bg-primary/10 text-primary',
};

const sizeStyles = {
  sm: 'h-5 px-2 text-[10px]',
  md: 'h-6 px-2.5 text-xs',
};

export function Badge({
  children,
  variant = 'neutral',
  size = 'sm',
  className,
}: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center justify-center rounded-full font-semibold uppercase tracking-wide',
        variantStyles[variant],
        sizeStyles[size],
        className
      )}
    >
      {children}
    </span>
  );
}
