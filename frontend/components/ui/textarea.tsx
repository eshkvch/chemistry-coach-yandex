import * as React from 'react'

import { cn } from '@/lib/utils'

function Textarea({ className, ...props }: React.ComponentProps<'textarea'>) {
  return (
    <textarea
      data-slot="textarea"
      className={cn(
        'flex min-h-[100px] w-full rounded-xl border border-border bg-surface px-4 py-3 text-[15px] text-foreground placeholder:text-foreground-muted shadow-sm transition-all outline-none resize-none',
        'focus:border-accent focus:ring-2 focus:ring-accent/20',
        'disabled:cursor-not-allowed disabled:opacity-50',
        className,
      )}
      {...props}
    />
  )
}

export { Textarea }
