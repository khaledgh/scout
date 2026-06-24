import { useState } from 'react'
import { Modal, Button, Spinner } from '@/components/ui'
import { useBadges, useAwardBadge } from '@/hooks/useBadges'
import { toast } from 'sonner'
import type { MemberBadge } from '@/types'

interface Props {
  open: boolean
  onClose: () => void
  memberId: number
  earned?: MemberBadge[]
}

export function AwardBadgeModal({ open, onClose, memberId, earned }: Props) {
  const { data: badges, isLoading } = useBadges()
  const award = useAwardBadge(memberId)
  const [selected, setSelected] = useState<number | null>(null)

  const earnedIds = new Set((earned ?? []).map((mb) => mb.badge_id))

  const onAward = async () => {
    if (!selected) return
    await award.mutateAsync(selected)
    toast.success('تم منح الشارة')
    setSelected(null)
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title="منح شارة" size="lg">
      {isLoading ? <Spinner className="h-32" /> : (
        <div dir="rtl" className="space-y-4">
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 max-h-80 overflow-y-auto">
            {badges?.map((badge) => {
              const owned = earnedIds.has(badge.id)
              return (
                <button
                  key={badge.id}
                  type="button"
                  disabled={owned}
                  onClick={() => setSelected(badge.id)}
                  className={`p-3 rounded-xl border text-center transition-colors ${
                    owned ? 'opacity-40 cursor-not-allowed border-gray-200 dark:border-slate-700'
                      : selected === badge.id ? 'border-primary bg-primary/10'
                      : 'border-gray-200 dark:border-slate-700 hover:border-primary/50'
                  }`}
                >
                  <div className="text-2xl mb-1">🏅</div>
                  <p className="text-sm font-medium text-gray-900 dark:text-white leading-tight">{badge.name}</p>
                  <p className="text-xs text-secondary mt-1">+{badge.xp_reward} XP</p>
                  {owned && <p className="text-xs text-gray-400 mt-1">ممنوحة</p>}
                </button>
              )
            })}
          </div>
          <div className="flex gap-3 justify-end">
            <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
            <Button type="button" onClick={onAward} disabled={!selected} loading={award.isPending}>منح الشارة</Button>
          </div>
        </div>
      )}
    </Modal>
  )
}
