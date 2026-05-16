import * as React from 'react'
import { cn } from '@/lib/utils'

interface ChatBubbleProps extends React.ComponentProps<'div'> {
  variant: 'user' | 'persona'
  children: React.ReactNode
}

function ChatBubble({
  variant,
  children,
  className,
  ...props
}: ChatBubbleProps) {
  return (
    <div
      data-slot="chat-bubble"
      data-variant={variant}
      className={cn(
        'max-w-[85%] rounded-2xl px-4 py-3 text-body',
        variant === 'user' && 'ml-auto bg-surface-elevated text-foreground rounded-br-md',
        variant === 'persona' && 'mr-auto bg-surface text-foreground border border-border rounded-bl-md',
        className,
      )}
      {...props}
    >
      {children}
    </div>
  )
}

export { ChatBubble }
