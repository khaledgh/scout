import type { ReactNode } from 'react'

type Variant = 'green' | 'gold' | 'blue' | 'red' | 'gray' | 'purple'

const variants: Record<Variant, string> = {
  green:  'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  gold:   'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
  blue:   'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  red:    'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
  gray:   'bg-gray-100 text-gray-700 dark:bg-slate-700 dark:text-slate-300',
  purple: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
}

interface Props {
  children: ReactNode
  variant?: Variant
  className?: string
}

export function Badge({ children, variant = 'gray', className = '' }: Props) {
  return (
    <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium ${variants[variant]} ${className}`}>
      {children}
    </span>
  )
}
