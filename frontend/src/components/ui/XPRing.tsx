interface Props {
  xp: number
  level: number
  levelBaseXP?: number
  size?: number
}

export function XPRing({ xp, level, levelBaseXP = 100, size = 80 }: Props) {
  const xpForCurrentLevel = Math.pow(level - 1, 2) * levelBaseXP
  const xpForNextLevel = Math.pow(level, 2) * levelBaseXP
  const progress = xpForNextLevel > xpForCurrentLevel
    ? (xp - xpForCurrentLevel) / (xpForNextLevel - xpForCurrentLevel)
    : 1

  const r = (size - 12) / 2
  const circumference = 2 * Math.PI * r
  const dash = progress * circumference

  return (
    <div className="relative inline-flex items-center justify-center" style={{ width: size, height: size }}>
      <svg width={size} height={size} className="-rotate-90">
        <circle cx={size / 2} cy={size / 2} r={r} fill="none" stroke="currentColor" className="text-gray-200 dark:text-slate-700" strokeWidth={6} />
        <circle
          cx={size / 2} cy={size / 2} r={r}
          fill="none" stroke="currentColor" className="text-secondary"
          strokeWidth={6}
          strokeDasharray={`${dash} ${circumference}`}
          strokeLinecap="round"
          style={{ transition: 'stroke-dasharray 0.5s ease' }}
        />
      </svg>
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <span className="text-sm font-bold text-gray-900 dark:text-white tabular-nums">{level}</span>
        <span className="text-[10px] text-gray-500 dark:text-slate-400">مستوى</span>
      </div>
    </div>
  )
}
