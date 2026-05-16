'use client'

import * as React from 'react'
import { cn } from '@/lib/utils'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

type StatusDotStatus = 'success' | 'warning' | 'danger'

interface StatusDotProps extends React.ComponentProps<'span'> {
  status: StatusDotStatus
  label?: string
}

const statusColors: Record<StatusDotStatus, string> = {
  success: 'bg-success',
  warning: 'bg-warning',
  danger: 'bg-danger',
}

function StatusDot({
  status,
  label,
  className,
  ...props
}: StatusDotProps) {
  const dot = (
    <span
      data-slot="status-dot"
      data-status={status}
      className={cn(
        'inline-block size-2.5 rounded-full shrink-0',
        statusColors[status],
        className,
      )}
      {...props}
    />
  )

  if (label) {
    return (
      <Tooltip>
        <TooltipTrigger asChild>
          {dot}
        </TooltipTrigger>
        <TooltipContent>
          {label}
        </TooltipContent>
      </Tooltip>
    )
  }

  return dot
}

export { StatusDot }
