import { useState, useEffect } from 'react'
import { Award, Plus, Pencil, Trash2 } from 'lucide-react'
import { useBadges, useMemberLeaderboard, useCreateBadge, useUpdateBadge, useDeleteBadge } from '@/hooks/useBadges'
import { Card, Badge, Spinner, EmptyState, Avatar, XPRing, Button, Modal } from '@/components/ui'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/features/auth/AuthContext'
import { isAdmin as isAdminRole, isLeaderRole } from '@/lib/permissions'
import { toast } from 'sonner'
import type { Badge as BadgeType } from '@/types'

const medalEmoji = ['🥇', '🥈', '🥉']

const badgeSchema = z.object({
  name:        z.string().min(2, 'الاسم مطلوب'),
  category:    z.string().optional(),
  description: z.string().optional(),
  xp_reward:   z.number().min(0),
})
type BadgeInputs = z.infer<typeof badgeSchema>

function BadgeFormModal({ open, onClose, badge }: { open: boolean; onClose: () => void; badge?: BadgeType }) {
  const isEdit = !!badge
  const create = useCreateBadge()
  const update = useUpdateBadge(badge?.id ?? 0)
  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<BadgeInputs>({
    resolver: zodResolver(badgeSchema),
    defaultValues: { xp_reward: 10 },
  })

  useEffect(() => {
    if (badge) reset({ name: badge.name, category: badge.category, description: badge.description, xp_reward: badge.xp_reward })
    else reset({ name: '', category: '', description: '', xp_reward: 10 })
  }, [badge, open, reset])

  const onSubmit = async (data: BadgeInputs) => {
    if (isEdit) { await update.mutateAsync(data); toast.success('تم تحديث الشارة') }
    else { await create.mutateAsync(data); toast.success('تمت إضافة الشارة') }
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title={isEdit ? 'تعديل الشارة' : 'شارة جديدة'} size="md">
      <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4" dir="rtl">
        <div>
          <label className="label">الاسم *</label>
          <input {...register('name')} className="input" placeholder="اسم الشارة" />
          {errors.name && <p className="text-xs text-red-600 mt-1">{errors.name.message}</p>}
        </div>
        <div>
          <label className="label">الفئة</label>
          <input {...register('category')} className="input" placeholder="مهارات كشفية / صحة..." />
        </div>
        <div>
          <label className="label">نقاط الخبرة (XP)</label>
          <input {...register('xp_reward', { valueAsNumber: true })} type="number" min={0} className="input" />
        </div>
        <div>
          <label className="label">الوصف</label>
          <textarea {...register('description')} rows={3} className="input resize-none" />
        </div>
        <div className="flex gap-3 justify-end">
          <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
          <Button type="submit" loading={isSubmitting}>{isEdit ? 'حفظ' : 'إضافة'}</Button>
        </div>
      </form>
    </Modal>
  )
}

export function BadgesPage() {
  const { user } = useAuth()
  const isAdmin = isAdminRole(user)
  const isLeader = isLeaderRole(user)
  const [tab, setTab] = useState<'catalog' | 'leaderboard'>('catalog')
  const [sectionFilter, setSectionFilter] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<BadgeType | undefined>(undefined)

  const { data: badges, isLoading } = useBadges()
  const { data: leaderboard } = useMemberLeaderboard(sectionFilter || undefined)
  const del = useDeleteBadge()

  const handleDelete = async (badge: BadgeType) => {
    if (!confirm(`حذف ${badge.name}؟`)) return
    await del.mutateAsync(badge.id)
    toast.success('تم حذف الشارة')
  }

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">الشارات والمتصدرون</h1>
        {isLeader && tab === 'catalog' && (
          <Button onClick={() => { setEditing(undefined); setShowForm(true) }}><Plus size={16} />شارة جديدة</Button>
        )}
      </div>

      <div className="flex gap-2">
        {[{ key: 'catalog', label: 'كتالوج الشارات' }, { key: 'leaderboard', label: 'لوحة المتصدرين' }].map(({ key, label }) => (
          <button key={key} onClick={() => setTab(key as 'catalog' | 'leaderboard')}
            className={`px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
              tab === key ? 'bg-primary text-white' : 'bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-slate-400'
            }`}>{label}</button>
        ))}
      </div>

      {tab === 'catalog' && (
        isLoading ? <Spinner className="h-48" /> : !badges?.length ? (
          <EmptyState icon={Award} title="لا توجد شارات" />
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {badges.map((badge) => (
              <Card key={badge.id} className="text-center relative" padding="sm">
                {isAdmin && (
                  <div className="absolute top-1.5 left-1.5 flex items-center gap-0.5">
                    <button onClick={() => { setEditing(badge); setShowForm(true) }} title="تعديل"
                      className="p-1 rounded-md text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                      <Pencil size={13} />
                    </button>
                    <button onClick={() => handleDelete(badge)} title="حذف"
                      className="p-1 rounded-md text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                      <Trash2 size={13} />
                    </button>
                  </div>
                )}
                <div className="w-16 h-16 rounded-full bg-secondary/20 flex items-center justify-center mx-auto mb-3 text-3xl">
                  {badge.icon_url ? <img src={badge.icon_url} className="w-full h-full object-cover rounded-full" /> : '🏅'}
                </div>
                <p className="font-semibold text-sm text-gray-900 dark:text-white leading-tight">{badge.name}</p>
                <p className="text-xs text-gray-500 dark:text-slate-400 mt-0.5">{badge.category}</p>
                {badge.xp_reward > 0 && (
                  <Badge variant="gold" className="mt-2">+{badge.xp_reward} XP</Badge>
                )}
              </Card>
            ))}
          </div>
        )
      )}

      {tab === 'leaderboard' && (
        <>
          <div className="flex gap-2">
            <select value={sectionFilter} onChange={(e) => setSectionFilter(e.target.value)} className="input w-40">
              <option value="">كل الشعب</option>
              <option value="ashbal">أشبال</option>
              <option value="kashaf">كشاف</option>
              <option value="jawala">جوالة</option>
              <option value="mukashe">مكاشفة</option>
            </select>
          </div>
          <Card>
            <div className="space-y-2">
              {leaderboard?.map((entry: { member_id: number; full_name: string; section: string; xp_total: number; level: number; photo_url?: string }, i: number) => (
                <div key={entry.member_id} className="flex items-center gap-4 p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm flex-shrink-0 ${
                    i < 3 ? 'bg-transparent text-2xl' : 'bg-gray-100 dark:bg-slate-700 text-gray-600'
                  }`}>
                    {i < 3 ? medalEmoji[i] : i + 1}
                  </div>
                  <Avatar name={entry.full_name} url={entry.photo_url} size="sm" />
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold text-gray-900 dark:text-white truncate">{entry.full_name}</p>
                    <p className="text-xs text-gray-500">{entry.section}</p>
                  </div>
                  <XPRing xp={entry.xp_total} level={entry.level} size={48} />
                  <div className="text-left">
                    <p className="font-bold text-lg text-primary tabular-nums">{entry.xp_total}</p>
                    <p className="text-xs text-gray-400">XP</p>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </>
      )}

      <BadgeFormModal open={showForm} onClose={() => setShowForm(false)} badge={editing} />
    </div>
  )
}
