import { useState, useRef } from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';

interface RippleButtonProps {
  children: React.ReactNode;
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
  className?: string;
  disabled?: boolean;
  color?: 'blue' | 'red' | 'green' | 'yellow';
}

interface RippleStyle {
  left: number;
  top: number;
  width: number;
  height: number;
  transform: string;
}

const buttonVariants = cva(
  cn(
    'overflow-hidden relative inline-flex items-center justify-center gap-2 cursor-pointer',
    'disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed',
    'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
    "[&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none",
    'whitespace-nowrap rounded-md text-sm font-medium transition-all aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive'
  ),
  {
    variants: {
      variant: {
        default:
          'bg-primary text-primary-foreground shadow-xs hover:bg-primary/90',
        destructive:
          'bg-destructive text-white shadow-xs hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60',
        outline:
          'bg-background text-primary shadow-xs hover:bg-accent dark:bg-input/30 dark:border-input dark:hover:bg-input/50',
        secondary:
          'bg-secondary text-secondary-foreground shadow-xs hover:bg-secondary/80',
        ghost:
          'hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50',
        link: 'text-primary underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-9 px-4 py-2 has-[>svg]:px-3',
        sm: 'h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5',
        lg: 'h-10 rounded-md px-6 has-[>svg]:px-4',
        icon: 'size-9',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

const rippleColors = {
  default: 'bg-primary',
  destructive: 'bg-destructive',
  outline: 'bg-secondary',
  secondary: 'bg-secondary',
  ghost: 'bg-accent',
  link: 'bg-primary',
};

function Button({
  className,
  variant,
  size,
  asChild = false,
  loading,
  disabled,
  children,
  ...props
}: React.ComponentProps<'button'> &
  VariantProps<typeof buttonVariants> & {
    loading?: boolean;
    asChild?: boolean;
  }) {
  const [ripples, setRipples] = useState<RippleStyle[]>([]);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const Comp = asChild ? Slot : 'button';

  const handleMouseDown = (event: React.MouseEvent<HTMLButtonElement>) => {
    console.log('handleMouseDown');
    if (disabled) return;

    const button = buttonRef.current;
    if (!button) return;

    // Get button dimensions and position
    const rect = button.getBoundingClientRect();

    // Calculate ripple size (should be bigger than the button for full effect)
    const size = Math.max(rect.width, rect.height) * 2;

    // Calculate ripple position
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    // Create new ripple
    const newRipple: RippleStyle = {
      left: x,
      top: y,
      width: size,
      height: size,
      transform: 'scale(0)',
    };

    setRipples([...ripples, newRipple]);
  };

  const handleLeave = () => {
    setRipples([]);
  };

  return (
    <Comp
      className={cn(buttonVariants({ variant, size, className }))}
      disabled={loading || disabled}
      onMouseDown={handleMouseDown}
      onMouseUp={handleLeave}
      onMouseLeave={handleLeave}
      ref={buttonRef}
      {...props}
    >
      {ripples.map((style, i) => (
        <span
          key={i}
          className={cn(
            'absolute rounded-full pointer-events-none animate-ripple',
            rippleColors[variant || 'default']
          )}
          style={{
            left: style.left - style.width / 2,
            top: style.top - style.height / 2,
            width: style.width,
            height: style.height,
            zIndex: 0,
          }}
        />
      ))}
      {loading ? (
        <Loader2 className='m-2 h-4 w-4 animate-spin' />
      ) : (
        <div className='z-10'>{children}</div>
      )}
    </Comp>
  );
}

export { Button, buttonVariants };
