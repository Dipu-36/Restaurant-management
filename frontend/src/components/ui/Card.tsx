import { cn } from '../../lib/utils';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  variant?: 'default' | 'interactive' | 'flat';
  onClick?: () => void;
}

export function Card({
  children,
  className,
  variant = 'default',
  onClick,
}: CardProps) {
  const isInteractive = variant === 'interactive';

  return (
    <div
      onClick={onClick}
      role={isInteractive ? 'button' : undefined}
      tabIndex={isInteractive ? 0 : undefined}
      onKeyDown={
        isInteractive && onClick
          ? (e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                onClick();
              }
            }
          : undefined
      }
      className={cn(
        'rounded-2xl bg-surface',
        variant === 'default' && 'border border-border shadow-soft',
        variant === 'flat' && 'border border-border',
        isInteractive &&
          'cursor-pointer transition-all duration-200 hover:shadow-medium hover:-translate-y-0.5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2',
        className
      )}
    >
      {children}
    </div>
  );
}