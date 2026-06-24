import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, CalendarDays, MapPin, Pencil, Trash2 } from 'lucide-react'
import { useActivities, useCreateActivity, useUpdateActivity, useDeleteActivity } from '@/hooks/useActivities'
import { Button, Card, Badge, Spinner, EmptyState, Modal, Select, DateTimeInput } from '@/components/ui'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/features/auth/AuthContext'
import { isLeaderRole } from '@/lib/permissions'
import { toast } from 'sonner'
import { format } from 'date-fns'
import { arLB as ar } from '@/lib/arLB'
import type { Activity } from '@/types'

const schema = z.object({
  title:       z.string().min(2, 'العنوان مطلوب'),
  type:        z.enum(['camp', 'hike', 'training', 'meeting', 'service']),
  status:      z.enum(['planned', 'ongoing', 'completed', 'cancelled']),
  location:    z.string().optional(),
  starts_at:   z.string().min(1, 'وقت البداية مطلوب'),
  ends_at:     z.string().min(1, 'وقت النهاية مطلوب'),
  description: z.string().optional(),
})
type Inputs = z.infer<typeof schema>

const typeOptions = [
  { value: 'camp',     label: 'مخيم' },
  { value: 'hike',     label: 'مسير' },
  { value: 'training', label: 'تدريب' },
  { value: 'meeting',  label: 'اجتماع' },
  { value: 'service',  label: 'خدمة' },
]
const statusOptions = [
  { value: 'planned',   label: 'مجدول' },
  { value: 'ongoing',   label: 'جارٍ' },
  { value: 'completed', label: 'مكتمل' },
  { value: 'cancelled', label: 'ملغى' },
]
const typeLabel: Record<string, string>  = Object.fromEntries(typeOptions.map(o => [o.value, o.label]))
const typeColor: Record<string, 'green' | 'blue' | 'gold' | 'purple' | 'gray'> = {
  camp: 'green', hike: 'blue', training: 'gold', meeting: 'gray', service: 'purple',
}
const statusLabel: Record<string, string> = Object.fromEntries(statusOptions.map(o => [o.value, o.label]))
const statusColor: Record<string, 'blue' | 'green' | 'gray' | 'red'> = {
  planned: 'blue', ongoing: 'green', completed: 'gray', cancelled: 'red',
}

export function ActivitiesPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Activity | undefined>(undefined)
  const [typeFilter, setTypeFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)

  const isLeader = isLeaderRole(user)

  const { data, isLoading } = useActivities({ type: typeFilter || undefined, status: statusFilter || undefined, page })
  const createActivity = useCreateActivity()
  const updateActivity = useUpdateActivity(editing?.id ?? 0)
  const deleteActivity = useDeleteActivity()

  const { register, handleSubmit, reset, control, formState: { errors, isSubmitting } } = useForm<Inputs>({
    resolver: zodResolver(schema),
    defaultValues: { type: 'meeting', status: 'planned' },
  })

  useEffect(() => {
    if (editing) {
      reset({
        title: editing.title, type: editing.type, status: editing.status,
        location: editing.location ?? '',
        starts_at: editing.starts_at?.slice(0, 16),
        ends_at: editing.ends_at?.slice(0, 16),
        description: editing.description ?? '',
      })
    } else {
      reset({ type: 'meeting', status: 'planned', title: '', location: '', starts_at: '', ends_at: '', description: '' })
    }
  }, [editing, showForm, reset])

  const onSubmit = async (data: Inputs) => {
    if (editing) { await updateActivity.mutateAsync(data); toast.success('تم تحديث النشاط') }
    else { await createActivity.mutateAsync(data); toast.success('تم إنشاء النشاط') }
    setShowForm(false)
  }

  const handleDelete = async (a: Activity) => {
    if (!confirm(`حذف ${a.title}؟`)) return
    await deleteActivity.mutateAsync(a.id)
    toast.success('تم حذف النشاط')
  }

  const openCreate = () => { setEditing(undefined); setShowForm(true) }
  const openEdit = (a: Activity) => { setEditing(a); setShowForm(true) }

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">الأنشطة</h1>
        {isLeader && <Button onClick={openCreate}><Plus size={16} />نشاط جديد</Button>}
      </div>

      <div className="flex gap-3 flex-wrap">
        <Select
          value={typeFilter}
          onChange={setTypeFilter}
          className="w-40"
          placeholder="كل الأنواع"
          options={[{ value: '', label: 'كل الأنواع' }, ...typeOptions]}
        />
        <Select
          value={statusFilter}
          onChange={setStatusFilter}
          className="w-40"
          placeholder="كل الحالات"
          options={[{ value: '', label: 'كل الحالات' }, ...statusOptions]}
        />
      </div>

      {isLoading ? <Spinner className="h-48" /> : !data?.data?.length ? (
        <EmptyState icon={CalendarDays} title="لا توجد أنشطة" />
      ) : (
        <>
          <div className="space-y-3">
            {data.data.map((a) => (
              <Card key={a.id} className="cursor-pointer hover:shadow-card-hover transition-shadow"
                onClick={() => navigate(`/activities/${a.id}`)}>
                <div className="flex items-start gap-4">
                  <div className="w-12 h-12 rounded-xl bg-primary/10 dark:bg-primary/20 flex items-center justify-center flex-shrink-0">
                    <CalendarDays size={22} className="text-primary" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <h3 className="font-semibold text-gray-900 dark:text-white">{a.title}</h3>
                      <div className="flex gap-2 flex-shrink-0 items-center">
                        <Badge variant={typeColor[a.type] ?? 'gray'}>{typeLabel[a.type]}</Badge>
                        <Badge variant={statusColor[a.status] ?? 'gray'}>{statusLabel[a.status]}</Badge>
                        {isLeader && (
                          <div className="flex items-center gap-1" onClick={(e) => e.stopPropagation()}>
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
                    </div>
                    <div className="flex flex-wrap gap-4 mt-2 text-sm text-gray-500 dark:text-slate-400">
                      <span className="flex items-center gap-1">
                        <CalendarDays size={14} />
                        {format(new Date(a.starts_at), 'PPp', { locale: ar })}
                      </span>
                      {a.location && (
                        <span className="flex items-center gap-1">
                          <MapPin size={14} />{a.location}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </Card>
            ))}
          </div>
          {data.meta && data.meta.total_pages > 1 && (
            <div className="flex gap-2 justify-center">
              <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(p => p - 1)}>السابق</Button>
              <Button variant="outline" size="sm" disabled={page >= data.meta.total_pages} onClick={() => setPage(p => p + 1)}>التالي</Button>
            </div>
          )}
        </>
      )}

      <Modal open={showForm} onClose={() => setShowForm(false)} title={editing ? 'تعديل النشاط' : 'نشاط جديد'} size="lg">
        <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4" dir="rtl">
          <div>
            <label className="label">عنوان النشاط *</label>
            <input {...register('title')} className="input" placeholder="عنوان النشاط" />
            {errors.title && <p className="text-xs text-red-600 mt-1">{errors.title.message}</p>}
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="label">النوع *</label>
              <Controller name="type" control={control} render={({ field }) => (
                <Select options={typeOptions} value={field.value} onChange={field.onChange} onBlur={field.onBlur} />
              )} />
            </div>
            <div>
              <label className="label">الموقع</label>
              <input {...register('location')} className="input" placeholder="المكان" />
            </div>
            {editing && (
              <div>
                <label className="label">الحالة</label>
                <Controller name="status" control={control} render={({ field }) => (
                  <Select options={statusOptions} value={field.value} onChange={field.onChange} onBlur={field.onBlur} />
                )} />
              </div>
            )}
            <div>
              <DateTimeInput
                {...register('starts_at')}
                label="وقت البداية *"
                error={errors.starts_at?.message}
              />
            </div>
            <div>
              <DateTimeInput
                {...register('ends_at')}
                label="وقت النهاية *"
                error={errors.ends_at?.message}
              />
            </div>
          </div>
          <div>
            <label className="label">وصف</label>
            <textarea {...register('description')} rows={3} className="input resize-none" placeholder="تفاصيل النشاط" />
          </div>
          <div className="flex gap-3 justify-end">
            <Button type="button" variant="ghost" onClick={() => setShowForm(false)}>إلغاء</Button>
            <Button type="submit" loading={isSubmitting}>{editing ? 'حفظ' : 'إنشاء'}</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
