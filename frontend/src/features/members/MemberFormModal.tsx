import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Modal, Button } from '@/components/ui'
import { useCreateMember, useUpdateMember } from '@/hooks/useMembers'
import { toast } from 'sonner'
import type { Member } from '@/types'

const schema = z.object({
  full_name:       z.string().min(2, 'الاسم مطلوب'),
  birth_date:      z.string().min(1, 'تاريخ الميلاد مطلوب'),
  gender:          z.enum(['male', 'female']),
  section:         z.enum(['ashbal', 'kashaf', 'jawala', 'mukashe']),
  rank_stage:      z.string().optional(),
  join_date:       z.string().min(1, 'تاريخ الانتساب مطلوب'),
  parent_name:     z.string().optional(),
  parent_phone:    z.string().min(8, 'رقم هاتف ولي الأمر مطلوب'),
  secondary_phone: z.string().optional(),
  address:         z.string().optional(),
  status:          z.enum(['active', 'inactive']),
})

type Inputs = z.infer<typeof schema>

interface Props {
  open: boolean
  onClose: () => void
  member?: Member
}

export function MemberFormModal({ open, onClose, member }: Props) {
  const isEdit = !!member
  const create = useCreateMember()
  const update = useUpdateMember(member?.id ?? 0)

  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<Inputs>({
    resolver: zodResolver(schema),
    defaultValues: {
      gender: 'male', section: 'kashaf', status: 'active',
      join_date: new Date().toISOString().slice(0, 10),
    },
  })

  useEffect(() => {
    if (member) {
      reset({
        full_name: member.full_name,
        birth_date: member.birth_date?.slice(0, 10),
        gender: member.gender,
        section: member.section,
        rank_stage: member.rank_stage ?? '',
        join_date: member.join_date?.slice(0, 10),
        parent_name: member.parent_name ?? '',
        parent_phone: member.parent_phone ?? '',
        secondary_phone: member.secondary_phone ?? '',
        address: member.address ?? '',
        status: member.status,
      })
    } else {
      reset({
        gender: 'male', section: 'kashaf', status: 'active',
        join_date: new Date().toISOString().slice(0, 10),
        full_name: '', birth_date: '', rank_stage: '', parent_name: '',
        parent_phone: '', secondary_phone: '', address: '',
      })
    }
  }, [member, open, reset])

  const onSubmit = async (data: Inputs) => {
    if (isEdit) {
      await update.mutateAsync(data)
      toast.success('تم تحديث بيانات العضو')
    } else {
      await create.mutateAsync(data)
      toast.success('تم إضافة العضو بنجاح')
    }
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title={isEdit ? 'تعديل بيانات العضو' : 'إضافة عضو جديد'} size="lg">
      <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4" dir="rtl">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="label">الاسم الكامل *</label>
            <input {...register('full_name')} className="input" placeholder="اسم العضو" />
            {errors.full_name && <p className="text-xs text-red-600 mt-1">{errors.full_name.message}</p>}
          </div>
          <div>
            <label className="label">الشعبة *</label>
            <select {...register('section')} className="input">
              <option value="ashbal">أشبال</option>
              <option value="kashaf">كشاف</option>
              <option value="jawala">جوالة</option>
              <option value="mukashe">مكاشفة</option>
            </select>
          </div>
          <div>
            <label className="label">تاريخ الميلاد *</label>
            <input {...register('birth_date')} type="date" className="input" />
            {errors.birth_date && <p className="text-xs text-red-600 mt-1">{errors.birth_date.message}</p>}
          </div>
          <div>
            <label className="label">الجنس *</label>
            <select {...register('gender')} className="input">
              <option value="male">ذكر</option>
              <option value="female">أنثى</option>
            </select>
          </div>
          <div>
            <label className="label">تاريخ الانتساب *</label>
            <input {...register('join_date')} type="date" className="input" />
            {errors.join_date && <p className="text-xs text-red-600 mt-1">{errors.join_date.message}</p>}
          </div>
          <div>
            <label className="label">المرتبة / المرحلة</label>
            <input {...register('rank_stage')} className="input" placeholder="مثال: طلائعي أول" />
          </div>
          <div>
            <label className="label">هاتف ولي الأمر *</label>
            <input {...register('parent_phone')} type="tel" className="input" dir="ltr" placeholder="03XXXXXXX" />
            {errors.parent_phone && <p className="text-xs text-red-600 mt-1">{errors.parent_phone.message}</p>}
          </div>
          <div>
            <label className="label">هاتف إضافي</label>
            <input {...register('secondary_phone')} type="tel" className="input" dir="ltr" placeholder="03XXXXXXX" />
          </div>
          <div>
            <label className="label">اسم ولي الأمر</label>
            <input {...register('parent_name')} className="input" placeholder="اسم ولي الأمر" />
          </div>
          <div>
            <label className="label">الحالة</label>
            <select {...register('status')} className="input">
              <option value="active">نشط</option>
              <option value="inactive">غير نشط</option>
            </select>
          </div>
          <div className="sm:col-span-2">
            <label className="label">العنوان</label>
            <input {...register('address')} className="input" placeholder="المنطقة / القضاء" />
          </div>
        </div>
        <div className="flex gap-3 justify-end pt-2">
          <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
          <Button type="submit" loading={isSubmitting}>{isEdit ? 'حفظ التعديلات' : 'حفظ'}</Button>
        </div>
      </form>
    </Modal>
  )
}
