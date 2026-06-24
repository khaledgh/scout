import type { LucideIcon } from 'lucide-react'
import { TrendingUp, TrendingDown } from 'lucide-react'

interface Props {
  label: string
  value: string | number
  subtitle?: string
  icon?: LucideIcon
  trend?: number
  color?: 'purple' | 'green' | 'teal' | 'gold' | 'blue' | 'red'
}

const gradients: Record<NonNullable<Props['color']>, string> = {
  purple: 'from-violet-500 to-purple-600',
  green:  'from-emerald-400 to-green-500',
  teal:   'from-cyan-400 to-teal-500',
  gold:   'from-amber-400 to-orange-500',
  blue:   'from-blue-400 to-indigo-500',
  red:    'from-rose-400 to-red-500',
}

export function StatCard({ label, value, subtitle, icon: Icon, trend, color = 'purple' }: Props) {
  return (
    <div className="bg-white dark:bg-slate-800 rounded-2xl shadow-card border border-gray-100 dark:border-slate-700 p-5">
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 min-w-0">
          <p className="text-xs font-medium text-gray-500 dark:text-slate-400 uppercase tracking-wide">{label}</p>
          <p className="text-3xl font-extrabold text-gray-900 dark:text-white mt-2 tabular-nums leading-none">{value}</p>
          {subtitle && (
            <p className="text-xs text-gray-400 dark:text-slate-500 mt-1">{subtitle}</p>
          )}
          {trend !== undefined && (
            <div className={`flex items-center gap-1 mt-2 text-xs font-semibold ${trend >= 0 ? 'text-emerald-600' : 'text-red-500'}`}>
              {trend >= 0 ? <TrendingUp size={13} /> : <TrendingDown size={13} />}
              <span>{Math.abs(trend).toFixed(1)}%</span>
            </div>
          )}
        </div>
        {Icon && (
          <div className={`flex-shrink-0 w-12 h-12 rounded-2xl bg-gradient-to-br ${gradients[color]} flex items-center justify-center shadow-sm`}>
            <Icon size={22} className="text-white" />
          </div>
        )}
      </div>
    </div>
  )
}
