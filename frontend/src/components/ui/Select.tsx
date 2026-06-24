import { useState, useRef, useEffect } from 'react'
import { ChevronDown, Check } from 'lucide-react'

export interface SelectOption {
  value: string
  label: string
}

interface Props {
  options: SelectOption[]
  value?: string
  onChange?: (value: string) => void
  onBlur?: () => void
  name?: string
  placeholder?: string
  className?: string
}

export function Select({ options, value, onChange, onBlur, name, placeholder = 'اختر', className = '' }: Props) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  const selected = options.find(o => o.value === value)

  useEffect(() => {
    if (!open) return
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false)
        onBlur?.()
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [open, onBlur])

  return (
    <div ref={ref} dir="rtl" className={`relative ${className}`} onKeyDown={(e) => e.key === 'Escape' && setOpen(false)}>
      <button
        type="button"
        name={name}
        onClick={() => setOpen(o => !o)}
        className={`w-full flex items-center justify-between gap-2 px-3 py-2 text-sm rounded-xl border transition-all duration-150
          ${open
            ? 'border-primary shadow-[0_0_0_3px_rgb(var(--color-primary)/0.15)] bg-white dark:bg-slate-800'
            : 'border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-800 hover:border-gray-400 dark:hover:border-slate-500'
          } text-gray-900 dark:text-slate-100 focus:outline-none`}
      >
        <span className={selected ? 'font-medium' : 'text-gray-400 dark:text-slate-500 font-normal'}>
          {selected ? selected.label : placeholder}
        </span>
        <ChevronDown
          size={15}
          className={`flex-shrink-0 transition-transform duration-200 ${open ? 'rotate-180 text-primary' : 'text-gray-400'}`}
        />
      </button>

      {open && (
        <div className="absolute top-[calc(100%+6px)] right-0 left-0 z-50 bg-white dark:bg-slate-800 border border-gray-200 dark:border-slate-700 rounded-2xl shadow-xl py-1.5 overflow-hidden">
          {options.map(option => {
            const isSelected = option.value === value
            return (
              <button
                key={option.value}
                type="button"
                onClick={() => { onChange?.(option.value); setOpen(false); onBlur?.() }}
                className={`w-full flex items-center justify-between gap-3 px-4 py-2.5 text-sm text-right transition-colors
                  ${isSelected
                    ? 'bg-primary/10 dark:bg-primary/20 text-primary font-semibold'
                    : 'text-gray-700 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-slate-700/60'
                  }`}
              >
                <span>{option.label}</span>
                {isSelected && <Check size={14} className="text-primary flex-shrink-0" />}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
