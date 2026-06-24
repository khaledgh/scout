import { useForm } from 'react-hook-form'
import { Modal, Button } from '@/components/ui'
import { useCreateEvaluation } from '@/hooks/useMembers'
import { toast } from 'sonner'

interface Inputs {
  period: string
  discipline: number
  participation: number
  leadership: number
  skill: number
  overall: number
  notes: string
}

interface Props {
  open: boolean
  onClose: () => void
  memberId: number
}

const sliders: { name: keyof Inputs; label: string }[] = [
  { name: 'discipline', label: 'الانضباط' },
  { name: 'participation', label: 'المشاركة' },
  { name: 'leadership', label: 'القيادة' },
  { name: 'skill', label: 'المهارة' },
  { name: 'overall', label: 'التقييم العام' },
]

export function EvaluationFormModal({ open, onClose, memberId }: Props) {
  const createEval = useCreateEvaluation(memberId)
  const { register, handleSubmit, reset, formState: { isSubmitting } } = useForm<Inputs>({
    defaultValues: {
      period: `${new Date().getFullYear()}-Q${Math.floor(new Date().getMonth() / 3) + 1}`,
      discipline: 7, participation: 7, leadership: 7, skill: 7, overall: 7, notes: '',
    },
  })

  const onSubmit = async (data: Inputs) => {
    await createEval.mutateAsync({
      ...data,
      discipline: Number(data.discipline),
      participation: Number(data.participation),
      leadership: Number(data.leadership),
      skill: Number(data.skill),
      overall: Number(data.overall),
    })
    toast.success('تم حفظ التقييم')
    reset()
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title="تقييم جديد" size="md">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" dir="rtl">
        <div>
          <label className="label">الفترة</label>
          <input {...register('period')} className="input" placeholder="2026-Q1" dir="ltr" />
        </div>
        {sliders.map(({ name, label }) => (
          <div key={name}>
            <label className="label">{label} (0–10)</label>
            <input {...register(name)} type="number" min={0} max={10} className="input" />
          </div>
        ))}
        <div>
          <label className="label">ملاحظات</label>
          <textarea {...register('notes')} rows={3} className="input resize-none" placeholder="ملاحظات حول أداء العضو..." />
        </div>
        <div className="flex gap-3 justify-end">
          <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
          <Button type="submit" loading={isSubmitting}>حفظ التقييم</Button>
        </div>
      </form>
    </Modal>
  )
}
