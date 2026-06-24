import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { Modal, Button } from '@/components/ui'
import { useUpsertMedical } from '@/hooks/useMembers'
import { toast } from 'sonner'
import type { MemberMedical } from '@/types'

interface Inputs {
  blood_type: string
  allergies: string
  chronic_conditions: string
  medications: string
  emergency_notes: string
}

interface Props {
  open: boolean
  onClose: () => void
  memberId: number
  medical?: MemberMedical | null
}

export function MedicalFormModal({ open, onClose, memberId, medical }: Props) {
  const upsert = useUpsertMedical(memberId)
  const { register, handleSubmit, reset, formState: { isSubmitting } } = useForm<Inputs>()

  useEffect(() => {
    reset({
      blood_type: medical?.blood_type ?? '',
      allergies: medical?.allergies ?? '',
      chronic_conditions: medical?.chronic_conditions ?? '',
      medications: medical?.medications ?? '',
      emergency_notes: medical?.emergency_notes ?? '',
    })
  }, [medical, open, reset])

  const onSubmit = async (data: Inputs) => {
    await upsert.mutateAsync(data)
    toast.success('تم حفظ المعلومات الطبية')
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title="المعلومات الطبية" size="lg">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" dir="rtl">
        <div>
          <label className="label">فصيلة الدم</label>
          <input {...register('blood_type')} className="input" placeholder="مثال: O+" dir="ltr" />
        </div>
        <div>
          <label className="label">الحساسية</label>
          <textarea {...register('allergies')} rows={2} className="input resize-none" placeholder="حساسية من أدوية أو أطعمة..." />
        </div>
        <div>
          <label className="label">أمراض مزمنة</label>
          <textarea {...register('chronic_conditions')} rows={2} className="input resize-none" placeholder="ربو، سكري..." />
        </div>
        <div>
          <label className="label">الأدوية</label>
          <textarea {...register('medications')} rows={2} className="input resize-none" placeholder="أدوية يتناولها بانتظام..." />
        </div>
        <div>
          <label className="label">ملاحظات للطوارئ</label>
          <textarea {...register('emergency_notes')} rows={2} className="input resize-none" placeholder="معلومات مهمة في حالات الطوارئ..." />
        </div>
        <div className="flex gap-3 justify-end">
          <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
          <Button type="submit" loading={isSubmitting}>حفظ</Button>
        </div>
      </form>
    </Modal>
  )
}
