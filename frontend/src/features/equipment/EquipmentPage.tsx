import { useState, useEffect } from 'react'
import { Plus, Package, Pencil, Trash2 } from 'lucide-react'
import { useEquipment, useCreateEquipment, useUpdateEquipment, useDeleteEquipment } from '@/hooks/useEquipment'
import { Button, Badge, Spinner, EmptyState, Modal } from '@/components/ui'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/features/auth/AuthContext'
import { isLeaderRole } from '@/lib/permissions'
import { toast } from 'sonner'
import type { Equipment } from '@/types'

const schema = z.object({
  name:               z.string().min(2, 'الاسم مطلوب'),
  category:           z.string().optional(),
  quantity_total:     z.number().min(0),
  quantity_available: z.number().min(0),
  condition:          z.string().optional(),
  notes:              z.string().optional(),
})
type Inputs = z.infer<typeof schema>

function EquipmentFormModal({ open, onClose, item }: { open: boolean; onClose: () => void; item?: Equipment }) {
  const isEdit = !!item
  const create = useCreateEquipment()
  const update = useUpdateEquipment(item?.id ?? 0)
  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<Inputs>({
    resolver: zodResolver(schema),
    defaultValues: { quantity_total: 1, quantity_available: 1 },
  })

  useEffect(() => {
    if (item) {
      reset({
        name: item.name, category: item.category, quantity_total: item.quantity_total,
        quantity_available: item.quantity_available, condition: item.condition, notes: item.notes,
      })
    } else {
      reset({ name: '', category: '', quantity_total: 1, quantity_available: 1, condition: '', notes: '' })
    }
  }, [item, open, reset])

  const onSubmit = async (data: Inputs) => {
    if (isEdit) { await update.mutateAsync(data); toast.success('تم تحديث المعدة') }
    else { await create.mutateAsync(data); toast.success('تمت إضافة المعدة') }
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title={isEdit ? 'تعديل معدة' : 'معدة جديدة'} size="lg">
      <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4" dir="rtl">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="label">الاسم *</label>
            <input {...register('name')} className="input" placeholder="مثال: خيمة 6 أشخاص" />
            {errors.name && <p className="text-xs text-red-600 mt-1">{errors.name.message}</p>}
          </div>
          <div>
            <label className="label">الفئة</label>
            <input {...register('category')} className="input" placeholder="تخييم / طبخ / إسعاف..." />
          </div>
          <div>
            <label className="label">الكمية الإجمالية</label>
            <input {...register('quantity_total', { valueAsNumber: true })} type="number" min={0} className="input" />
          </div>
          <div>
            <label className="label">المتاح</label>
            <input {...register('quantity_available', { valueAsNumber: true })} type="number" min={0} className="input" />
          </div>
          <div>
            <label className="label">الحالة</label>
            <input {...register('condition')} className="input" placeholder="جيدة / تحتاج صيانة..." />
          </div>
          <div className="sm:col-span-2">
            <label className="label">ملاحظات</label>
            <textarea {...register('notes')} rows={2} className="input resize-none" />
          </div>
        </div>
        <div className="flex gap-3 justify-end">
          <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
          <Button type="submit" loading={isSubmitting}>{isEdit ? 'حفظ' : 'إضافة'}</Button>
        </div>
      </form>
    </Modal>
  )
}

export function EquipmentPage() {
  const { user } = useAuth()
  const canManage = isLeaderRole(user)
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Equipment | undefined>(undefined)

  const { data: items, isLoading } = useEquipment()
  const del = useDeleteEquipment()

  const handleDelete = async (item: Equipment) => {
    if (!confirm(`حذف ${item.name}؟`)) return
    await del.mutateAsync(item.id)
    toast.success('تم الحذف')
  }

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">المعدات</h1>
        {canManage && <Button onClick={() => { setEditing(undefined); setShowForm(true) }}><Plus size={16} />معدة جديدة</Button>}
      </div>

      {isLoading ? <Spinner className="h-48" /> : !items?.length ? (
        <EmptyState icon={Package} title="لا توجد معدات" description="أضف أول معدة للمخزون" />
      ) : (
        <div className="card overflow-hidden" style={{ padding: 0 }}>
          <table className="w-full text-sm">
            <thead className="bg-gray-50 dark:bg-slate-900/50">
              <tr>
                {['المعدة', 'الفئة', 'المتاح / الإجمالي', 'الحالة', ''].map((h) => (
                  <th key={h} className="px-4 py-3 text-right font-medium text-gray-500 dark:text-slate-400 text-xs">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-slate-700">
              {items.map((item) => (
                <tr key={item.id} className="hover:bg-gray-50 dark:hover:bg-slate-700/30 transition-colors">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <div className="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center">
                        <Package size={16} className="text-primary" />
                      </div>
                      <span className="font-medium text-gray-900 dark:text-white">{item.name}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-gray-500">{item.category || '—'}</td>
                  <td className="px-4 py-3">
                    <Badge variant={item.quantity_available > 0 ? 'green' : 'red'}>
                      {item.quantity_available} / {item.quantity_total}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-gray-500">{item.condition || '—'}</td>
                  <td className="px-4 py-3 text-left">
                    {canManage && (
                      <div className="flex items-center gap-1 justify-end">
                        <button onClick={() => { setEditing(item); setShowForm(true) }} title="تعديل"
                          className="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                          <Pencil size={15} />
                        </button>
                        <button onClick={() => handleDelete(item)} title="حذف"
                          className="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                          <Trash2 size={15} />
                        </button>
                      </div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <EquipmentFormModal open={showForm} onClose={() => setShowForm(false)} item={editing} />
    </div>
  )
}
