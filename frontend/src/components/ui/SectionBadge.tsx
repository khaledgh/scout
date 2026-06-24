import { Badge } from './Badge'
import type { Section } from '@/types'

const labels: Record<Section, string> = {
  ashbal:  'أشبال',
  kashaf:  'كشاف',
  jawala:  'جوالة',
  mukashe: 'مكاشفة',
}
const colors: Record<Section, 'green' | 'blue' | 'gold' | 'purple'> = {
  ashbal:  'green',
  kashaf:  'blue',
  jawala:  'gold',
  mukashe: 'purple',
}

export function SectionBadge({ section }: { section: Section }) {
  return <Badge variant={colors[section]}>{labels[section]}</Badge>
}
