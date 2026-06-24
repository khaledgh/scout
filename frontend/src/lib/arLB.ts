import { ar } from 'date-fns/locale'
import type { Locale } from 'date-fns'

const months = [
  'كانون الثاني', 'شباط',       'آذار',
  'نيسان',        'أيار',        'حزيران',
  'تموز',         'آب',          'أيلول',
  'تشرين الأول',  'تشرين الثاني', 'كانون الأول',
]

export const arLB: Locale = {
  ...ar,
  localize: {
    ...ar.localize!,
    month: (index: number) => months[index] ?? months[0],
  },
}
