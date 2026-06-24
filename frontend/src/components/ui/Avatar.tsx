import { assetUrl } from '@/lib/assetUrl'

interface Props {
  name: string
  url?: string | null
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  className?: string
}

const sizes = {
  xs: 'w-6 h-6 text-xs',
  sm: 'w-8 h-8 text-sm',
  md: 'w-10 h-10 text-sm',
  lg: 'w-14 h-14 text-lg',
  xl: 'w-20 h-20 text-2xl',
}

const colors = [
  'bg-primary-500', 'bg-accent-500', 'bg-secondary-500',
  'bg-purple-500', 'bg-pink-500', 'bg-teal-500',
]

function colorFromName(name: string) {
  let sum = 0
  for (let i = 0; i < name.length; i++) sum += name.charCodeAt(i)
  return colors[sum % colors.length]
}

export function Avatar({ name, url, size = 'md', className = '' }: Props) {
  const initials = name.split(' ').map((w) => w[0]).slice(0, 2).join('')
  if (url) {
    return (
      <img
        src={assetUrl(url)}
        alt={name}
        className={`rounded-full object-cover flex-shrink-0 ${sizes[size]} ${className}`}
      />
    )
  }
  return (
    <div className={`rounded-full flex items-center justify-center text-white font-semibold flex-shrink-0 ${colorFromName(name)} ${sizes[size]} ${className}`}>
      {initials}
    </div>
  )
}
