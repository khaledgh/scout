import { Calendar, Clock } from 'lucide-react'
import type { InputHTMLAttributes } from 'react'

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export function DateTimeInput({ label, error, className = '', ...props }: Props) {
  return (
    <div className="space-y-1">
      {label && <label className="label">{label}</label>}
      <div className={`flex items-center gap-0 rounded-xl border transition-all duration-150 bg-white dark:bg-slate-800 overflow-hidden
        ${error
          ? 'border-red-400 shadow-[0_0_0_3px_rgb(239_68_68/0.12)]'
          : 'border-gray-300 dark:border-slate-600 focus-within:border-primary focus-within:shadow-[0_0_0_3px_rgb(var(--color-primary)/0.15)]'
        }`}>
        {/* Date part */}
        <div className="flex items-center gap-1.5 px-3 py-2 border-l border-gray-200 dark:border-slate-700 flex-shrink-0">
          <Calendar size={14} className="text-primary flex-shrink-0" />
        </div>
        <input
          type="datetime-local"
          {...props}
          className={`flex-1 min-w-0 px-2 py-2 text-sm bg-transparent text-gray-900 dark:text-slate-100 focus:outline-none datetime-input ${className}`}
        />
        <div className="flex items-center gap-1.5 px-3 py-2 border-r border-gray-200 dark:border-slate-700 flex-shrink-0">
          <Clock size={14} className="text-gray-400 flex-shrink-0" />
        </div>
      </div>
      {error && <p className="text-xs text-red-600 mt-1">{error}</p>}
    </div>
  )
}
