import { useState, useEffect } from 'react'
import { Pin, Plus, Pencil, Trash2 } from 'lucide-react'
import { useAnnouncements, useCreateAnnouncement, useUpdateAnnouncement, useDeleteAnnouncement } from '@/hooks/useCommunication'
import { Card, Badge, Button, Spinner, EmptyState, Modal } from '@/components/ui'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/features/auth/AuthContext'
import { isLeaderRole, isAdmin } from '@/lib/permissions'
import { toast } from 'sonner'
import { format } from 'date-fns'
import { arLB as ar } from '@/lib/arLB'
import type { Announcement } from '@/types'

const schema = z.object({
  title:    z.string().min(2, 'العنوان مطلوب'),
  body:     z.string().min(5, 'المحتوى مطلوب'),
  audience: z.enum(['all', 'unit', 'leaders']),
  pinned:   z.boolean().optional(),
})
type Inputs = z.infer<typeof schema>

const audienceLabel: Record<string, string> = { all: 'الجميع', unit: 'الطليعة', leaders: 'القادة' }

export function CommunicationPage() {
  const { user } = useAuth()
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Announcement | undefined>(undefined)
  const isLeader = isLeaderRole(user)

  const { data: announcements, isLoading } = useAnnouncements()
  const createAnnouncement = useCreateAnnouncement()
  const updateAnnouncement = useUpdateAnnouncement(editing?.id ?? 0)
  const deleteAnnouncement = useDeleteAnnouncement()

  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<Inputs>({
    resolver: zodResolver(schema),
    defaultValues: { audience: 'all', pinned: false },
  })

  useEffect(() => {
    if (editing) reset({ title: editing.title, body: editing.body, audience: editing.audience, pinned: editing.pinned })
    else reset({ title: '', body: '', audience: 'all', pinned: false })
  }, [editing, showForm, reset])

  const onSubmit = async (data: Inputs) => {
    if (editing) { await updateAnnouncement.mutateAsync(data); toast.success('تم تحديث الإعلان') }
    else { await createAnnouncement.mutateAsync(data); toast.success('تم نشر الإعلان') }
    setShowForm(false)
  }

  const handleDelete = async (a: Announcement) => {
    if (!confirm('حذف هذا الإعلان؟')) return
    await deleteAnnouncement.mutateAsync(a.id)
    toast.success('تم حذف الإعلان')
  }

  const canManageItem = (a: Announcement) => isAdmin(user) || a.author_id === user?.id
  const openCreate = () => { setEditing(undefined); setShowForm(true) }
  const openEdit = (a: Announcement) => { setEditing(a); setShowForm(true) }

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">الإعلانات والتواصل</h1>
        {isLeader && <Button onClick={openCreate}><Plus size={16} />إعلان جديد</Button>}
      </div>

      {isLoading ? <Spinner className="h-48" /> : !announcements?.length ? (
        <EmptyState icon={Pin} title="لا توجد إعلانات" />
      ) : (
        <div className="space-y-3">
          {announcements.map((a) => (
            <Card key={a.id}>
              <div className="flex items-start gap-3">
                {a.pinned && <Pin size={16} className="text-secondary flex-shrink-0 mt-0.5" />}
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-semibold text-gray-900 dark:text-white">{a.title}</h3>
                    <Badge variant="blue">{audienceLabel[a.audience]}</Badge>
                    {a.pinned && <Badge variant="gold">مثبت</Badge>}
                    {canManageItem(a) && (
                      <div className="flex items-center gap-1 mr-auto">
                        <button onClick={() => openEdit(a)} title="تعديل"
                          className="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                          <Pencil size={14} />
                        </button>
                        <button onClick={() => handleDelete(a)} title="حذف"
                          className="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                          <Trash2 size={14} />
                        </button>
                      </div>
                    )}
                  </div>
                  <p className="text-sm text-gray-600 dark:text-slate-300 whitespace-pre-wrap">{a.body}</p>
                  <div className="flex items-center gap-3 mt-3 text-xs text-gray-400">
                    <span>{a.author?.full_name}</span>
                    <span>·</span>
                    <span>{a.published_at ? format(new Date(a.published_at), 'PPp', { locale: ar }) : ''}</span>
                  </div>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      <Modal open={showForm} onClose={() => setShowForm(false)} title={editing ? 'تعديل الإعلان' : 'إعلان جديد'} size="lg">
        <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4" dir="rtl">
          <div>
            <label className="label">العنوان *</label>
            <input {...register('title')} className="input" placeholder="عنوان الإعلان" />
            {errors.title && <p className="text-xs text-red-600 mt-1">{errors.title.message}</p>}
          </div>
          <div>
            <label className="label">المحتوى *</label>
            <textarea {...register('body')} rows={5} className="input resize-none" placeholder="نص الإعلان..." />
            {errors.body && <p className="text-xs text-red-600 mt-1">{errors.body.message}</p>}
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="label">الجمهور *</label>
              <select {...register('audience')} className="input">
                <option value="all">الجميع</option>
                <option value="unit">الطليعة</option>
                <option value="leaders">القادة فقط</option>
              </select>
            </div>
            <div className="flex items-center gap-2 mt-6">
              <input {...register('pinned')} type="checkbox" id="pinned" className="rounded" />
              <label htmlFor="pinned" className="text-sm text-gray-700 dark:text-slate-300">تثبيت الإعلان</label>
            </div>
          </div>
          <div className="flex gap-3 justify-end">
            <Button type="button" variant="ghost" onClick={() => setShowForm(false)}>إلغاء</Button>
            <Button type="submit" loading={isSubmitting}>{editing ? 'حفظ' : 'نشر'}</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
