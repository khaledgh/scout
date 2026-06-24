import { useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useCreateUnit, useUpdateUnit } from '@/hooks/useUnits'
import { Button, Modal, Select } from '@/components/ui'
import { toast } from 'sonner'
import type { Unit } from '@/types'

const schema = z.object({
  name:    z.string().min(2, 'اسم الطليعة مطلوب'),
  section: z.enum(['ashbal', 'kashaf', 'jawala', 'mukashe']),
  motto:   z.string().optional(),
})
type Inputs = z.infer<typeof schema>

const sectionOptions = [
  { value: 'ashbal',  label: 'أشبال' },
  { value: 'kashaf',  label: 'كشاف' },
  { value: 'jawala',  label: 'جوالة' },
  { value: 'mukashe', label: 'مكاشفة' },
]

interface Props {
  open: boolean
  onClose: () => void
  unit?: Unit
}

export function UnitFormModal({ open, onClose, unit }: Props) {
  const isEdit = !!unit
  const createUnit = useCreateUnit()
  const updateUnit = useUpdateUnit(unit?.id ?? 0)

  const { register, handleSubmit, reset, control, formState: { errors, isSubmitting } } = useForm<Inputs>({
    resolver: zodResolver(schema),
    defaultValues: { section: 'kashaf' },
  })

  useEffect(() => {
    if (unit) reset({ name: unit.name, section: unit.section, motto: unit.motto ?? '' })
    else reset({ name: '', section: 'kashaf', motto: '' })
  }, [unit, open, reset])

  const onSubmit = async (data: Inputs) => {
    if (isEdit) { await updateUnit.mutateAsync(data); toast.success('تم تحديث الطليعة') }
    else { await createUnit.mutateAsync(data); toast.success('تم إنشاء الطليعة') }
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title={isEdit ? 'تعديل الطليعة' : 'إنشاء طليعة جديدة'}>
      <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4" dir="rtl">
        <div>
          <label className="label">اسم الطليعة *</label>
          <input {...register('name')} className="input" placeholder="مثال: طليعة الأرز" />
          {errors.name && <p className="text-xs text-red-600 mt-1">{errors.name.message}</p>}
        </div>
        <div>
          <label className="label">الشعبة *</label>
          <Controller name="section" control={control} render={({ field }) => (
            <Select options={sectionOptions} value={field.value} onChange={field.onChange} onBlur={field.onBlur} />
          )} />
        </div>
        <div>
          <label className="label">الشعار</label>
          <input {...register('motto')} className="input" placeholder="شعار الطليعة" />
        </div>
        <div className="flex gap-3 justify-end">
          <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
          <Button type="submit" loading={isSubmitting}>{isEdit ? 'حفظ' : 'إنشاء'}</Button>
        </div>
      </form>
    </Modal>
  )
}
