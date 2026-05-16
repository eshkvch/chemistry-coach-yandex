import * as React from 'react'
import { cn } from '@/lib/utils'

interface ScoreBarProps extends React.ComponentProps<'div'> {
  value: number
  max?: number
  label: string
  showValue?: boolean
  size?: 'default' | 'compact'
}

function ScoreBar({
  value,
  max = 10,
  label,
  showValue = true,
  size = 'default',
  className,
  ...props
}: ScoreBarProps) {
  const percentage = Math.min(100, Math.max(0, (value / max) * 100))
  const isCompact = size === 'compact'
  
  return (
    <div
      data-slot="score-bar"
      className={cn(
        'flex items-center gap-3',
        isCompact && 'gap-2',
        className
      )}
      {...props}
    >
      <span className={cn(
        'text-foreground-secondary',
        isCompact ? 'text-[11px] min-w-[52px]' : 'text-small min-w-[80px]'
      )}>
        {label}
      </span>
      <div className={cn(
        'relative flex-1 bg-surface-elevated rounded-full overflow-hidden',
        isCompact ? 'h-1' : 'h-2'
      )}>
        <div
          className="absolute inset-y-0 left-0 bg-accent rounded-full transition-all duration-300"
          style={{ width: `${percentage}%` }}
        />
      </div>
      {showValue && (
        <span className={cn(
          'text-foreground font-semibold text-right',
          isCompact ? 'text-[11px] min-w-[16px]' : 'text-small min-w-[24px]'
        )}>
          {value}
        </span>
      )}
    </div>
  )
}

export { ScoreBar }
