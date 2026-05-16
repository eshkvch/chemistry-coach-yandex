import * as React from 'react'
import { cn } from '@/lib/utils'

interface ProgressLineProps extends React.ComponentProps<'div'> {
  steps: number
  currentStep: number
  labels?: string[]
}

function ProgressLine({
  steps,
  currentStep,
  labels,
  className,
  ...props
}: ProgressLineProps) {
  return (
    <div
      data-slot="progress-line"
      className={cn('flex flex-col gap-2', className)}
      {...props}
    >
      {/* Progress bars */}
      <div className="flex items-center gap-1.5">
        {Array.from({ length: steps }, (_, i) => {
          const stepNumber = i + 1
          const isCompleted = stepNumber < currentStep
          const isCurrent = stepNumber === currentStep
          
          return (
            <div
              key={i}
              className={cn(
                'h-1 flex-1 rounded-full transition-colors duration-200',
                isCompleted && 'bg-accent',
                isCurrent && 'bg-accent',
                !isCompleted && !isCurrent && 'bg-surface-elevated',
              )}
            />
          )
        })}
      </div>
      
      {/* Labels */}
      {labels && labels.length > 0 && (
        <div className="flex items-center justify-between">
          {labels.map((label, i) => {
            const stepNumber = i + 1
            const isCompleted = stepNumber < currentStep
            const isCurrent = stepNumber === currentStep
            
            return (
              <span
                key={i}
                className={cn(
                  'text-[11px] font-medium transition-colors',
                  isCompleted && 'text-foreground-secondary',
                  isCurrent && 'text-accent',
                  !isCompleted && !isCurrent && 'text-foreground-muted',
                )}
              >
                {label}
              </span>
            )
          })}
        </div>
      )}
    </div>
  )
}

export { ProgressLine }
