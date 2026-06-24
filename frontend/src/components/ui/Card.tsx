import type { ReactNode, MouseEventHandler } from 'react'

interface Props {
  children: ReactNode
  className?: string
  padding?: 'none' | 'sm' | 'md' | 'lg'
  onClick?: MouseEventHandler<HTMLDivElement>
}

const paddings = { none: '', sm: 'p-3', md: 'p-5', lg: 'p-8' }

export function Card({ children, className = '', padding = 'md', onClick }: Props) {
  return (
    <div
      onClick={onClick}
      className={`bg-white dark:bg-slate-800 rounded-2xl shadow-card border border-gray-100 dark:border-slate-700 ${paddings[padding]} ${className}`}
    >
      {children}
    </div>
  )
}
