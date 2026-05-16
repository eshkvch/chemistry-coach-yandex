import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'

import { cn } from '@/lib/utils'
import { Spinner } from '@/components/ui/spinner'

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-xl text-[15px] font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:ring-ring/50 focus-visible:ring-[3px]",
  {
    variants: {
      variant: {
        default:
          'bg-accent text-accent-foreground hover:bg-accent-hover active:scale-[0.98] shadow-sm',
        secondary:
          'border border-border bg-transparent text-foreground hover:bg-surface-elevated hover:border-accent/50 active:scale-[0.98]',
        outline:
          'border border-border bg-transparent text-foreground hover:bg-surface-elevated hover:border-accent/50 active:scale-[0.98]',
        ghost:
          'text-foreground-secondary hover:text-foreground hover:bg-surface-elevated active:scale-[0.98]',
        destructive:
          'bg-danger text-white hover:bg-danger/90 active:scale-[0.98] shadow-sm',
        link: 'text-accent underline-offset-4 hover:underline hover:text-accent-hover',
      },
      size: {
        default: 'h-10 px-5 py-2.5',
        sm: 'h-9 rounded-xl gap-1.5 px-4 text-[13px]',
        lg: 'h-12 rounded-xl px-6 text-[15px]',
        icon: 'size-10',
        'icon-sm': 'size-9',
        'icon-lg': 'size-12',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

export interface ButtonProps
  extends React.ComponentProps<'button'>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
  loading?: boolean
}

function Button({
  className,
  variant,
  size,
  asChild = false,
  loading = false,
  children,
  disabled,
  ...props
}: ButtonProps) {
  const Comp = asChild ? Slot : 'button'

  return (
    <Comp
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <>
          <Spinner className="size-4" />
          <span>{children}</span>
        </>
      ) : (
        children
      )}
    </Comp>
  )
}

export { Button, buttonVariants }
