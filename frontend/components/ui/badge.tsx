import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'

import { cn } from '@/lib/utils'

const badgeVariants = cva(
  'inline-flex items-center justify-center rounded-lg px-2.5 py-1 text-[13px] font-medium whitespace-nowrap shrink-0 [&>svg]:size-3.5 gap-1.5 [&>svg]:pointer-events-none transition-colors',
  {
    variants: {
      variant: {
        default:
          'bg-surface-elevated text-foreground-secondary border border-border',
        success:
          'bg-success/15 text-success border border-success/20',
        warning:
          'bg-warning/15 text-warning border border-warning/20',
        danger:
          'bg-danger/15 text-danger border border-danger/20',
        accent:
          'bg-accent/15 text-accent border border-accent/20',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  },
)

function Badge({
  className,
  variant,
  asChild = false,
  ...props
}: React.ComponentProps<'span'> &
  VariantProps<typeof badgeVariants> & { asChild?: boolean }) {
  const Comp = asChild ? Slot : 'span'

  return (
    <Comp
      data-slot="badge"
      className={cn(badgeVariants({ variant }), className)}
      {...props}
    />
  )
}

export { Badge, badgeVariants }
