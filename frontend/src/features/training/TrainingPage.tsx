import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { BookOpen, CheckCircle, Clock, Plus, Pencil, Trash2 } from 'lucide-react'
import { useTrainingLessons, useMyTrainingProgress, useCreateLesson, useUpdateLesson, useDeleteLesson } from '@/hooks/useTraining'
import { Card, Badge, Spinner, EmptyState, Button, Modal } from '@/components/ui'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/features/auth/AuthContext'
import { isLeaderRole } from '@/lib/permissions'
import { toast } from 'sonner'
import type { TrainingLesson } from '@/types'

const lessonSchema = z.object({
  title:        z.string().min(2, 'العنوان مطلوب'),
  category:     z.string().min(1, 'الفئة مطلوبة'),
  content:      z.string().optional(),
  order_index:  z.number().min(0),
  is_published: z.boolean().optional(),
})
type LessonInputs = z.infer<typeof lessonSchema>

function LessonFormModal({ open, onClose, lesson }: { open: boolean; onClose: () => void; lesson?: TrainingLesson }) {
  const isEdit = !!lesson
  const create = useCreateLesson()
  const update = useUpdateLesson(lesson?.id ?? 0)
  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<LessonInputs>({
    resolver: zodResolver(lessonSchema),
    defaultValues: { order_index: 1, is_published: true },
  })

  useEffect(() => {
    if (lesson) reset({ title: lesson.title, category: lesson.category, content: lesson.content, order_index: lesson.order_index, is_published: lesson.is_published })
    else reset({ title: '', category: '', content: '', order_index: 1, is_published: true })
  }, [lesson, open, reset])

  const onSubmit = async (data: LessonInputs) => {
    if (isEdit) { await update.mutateAsync(data); toast.success('تم تحديث الدرس') }
    else { await create.mutateAsync(data); toast.success('تمت إضافة الدرس') }
    onClose()
  }

  return (
    <Modal open={open} onClose={onClose} title={isEdit ? 'تعديل الدرس' : 'درس جديد'} size="lg">
      <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4" dir="rtl">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="label">العنوان *</label>
            <input {...register('title')} className="input" placeholder="عنوان الدرس" />
            {errors.title && <p className="text-xs text-red-600 mt-1">{errors.title.message}</p>}
          </div>
          <div>
            <label className="label">الفئة *</label>
            <input {...register('category')} className="input" placeholder="مهارات كشفية / صحة..." />
            {errors.category && <p className="text-xs text-red-600 mt-1">{errors.category.message}</p>}
          </div>
          <div>
            <label className="label">الترتيب</label>
            <input {...register('order_index', { valueAsNumber: true })} type="number" min={0} className="input" />
          </div>
          <div className="flex items-center gap-2 mt-6">
            <input {...register('is_published')} type="checkbox" id="published" className="rounded" />
            <label htmlFor="published" className="text-sm text-gray-700 dark:text-slate-300">منشور</label>
          </div>
        </div>
        <div>
          <label className="label">المحتوى</label>
          <textarea {...register('content')} rows={6} className="input resize-none" placeholder="محتوى الدرس..." />
        </div>
        <div className="flex gap-3 justify-end">
          <Button type="button" variant="ghost" onClick={onClose}>إلغاء</Button>
          <Button type="submit" loading={isSubmitting}>{isEdit ? 'حفظ' : 'إضافة'}</Button>
        </div>
      </form>
    </Modal>
  )
}

export function TrainingPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const canManage = isLeaderRole(user)
  const [categoryFilter, setCategoryFilter] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<TrainingLesson | undefined>(undefined)

  const { data: lessons, isLoading } = useTrainingLessons()
  const { data: progress } = useMyTrainingProgress()
  const del = useDeleteLesson()

  const categories = [...new Set(lessons?.map(l => l.category) ?? [])]
  const filtered = categoryFilter ? lessons?.filter(l => l.category === categoryFilter) : lessons

  const getProgress = (lessonId: number) =>
    progress?.find(p => p.lesson_id === lessonId)

  const handleDelete = async (lesson: TrainingLesson) => {
    if (!confirm(`حذف ${lesson.title}؟`)) return
    await del.mutateAsync(lesson.id)
    toast.success('تم حذف الدرس')
  }

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">مركز التدريب</h1>
        {canManage && <Button onClick={() => { setEditing(undefined); setShowForm(true) }}><Plus size={16} />درس جديد</Button>}
      </div>

      {categories.length > 0 && (
        <div className="flex gap-2 flex-wrap">
          <button onClick={() => setCategoryFilter('')}
            className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
              !categoryFilter ? 'bg-primary text-white' : 'bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-slate-400'
            }`}>الكل</button>
          {categories.map(cat => (
            <button key={cat} onClick={() => setCategoryFilter(cat)}
              className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
                categoryFilter === cat ? 'bg-primary text-white' : 'bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-slate-400'
              }`}>{cat}</button>
          ))}
        </div>
      )}

      {isLoading ? <Spinner className="h-48" /> : !filtered?.length ? (
        <EmptyState icon={BookOpen} title="لا توجد دروس" />
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((lesson) => {
            const prog = getProgress(lesson.id)
            return (
              <Card key={lesson.id} className="cursor-pointer hover:shadow-card-hover transition-shadow relative"
                onClick={() => navigate(`/training/${lesson.id}`)}>
                {canManage && (
                  <div className="absolute top-2 left-2 z-10 flex items-center gap-0.5" onClick={(e) => e.stopPropagation()}>
                    <button onClick={() => { setEditing(lesson); setShowForm(true) }} title="تعديل"
                      className="p-1.5 rounded-lg bg-white/90 dark:bg-slate-800/90 text-gray-500 hover:text-primary shadow-sm transition-colors">
                      <Pencil size={14} />
                    </button>
                    <button onClick={() => handleDelete(lesson)} title="حذف"
                      className="p-1.5 rounded-lg bg-white/90 dark:bg-slate-800/90 text-gray-500 hover:text-red-500 shadow-sm transition-colors">
                      <Trash2 size={14} />
                    </button>
                  </div>
                )}
                {lesson.cover_url ? (
                  <img src={lesson.cover_url} alt={lesson.title}
                    className="w-full h-32 object-cover rounded-xl mb-4 -mx-0 -mt-0" />
                ) : (
                  <div className="w-full h-32 bg-gradient-to-br from-accent/20 to-primary/20 rounded-xl mb-4 flex items-center justify-center">
                    <BookOpen size={40} className="text-primary/40" />
                  </div>
                )}
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-900 dark:text-white truncate">{lesson.title}</h3>
                    <Badge variant="blue" className="mt-1">{lesson.category}</Badge>
                  </div>
                  {prog ? (
                    prog.passed
                      ? <CheckCircle size={20} className="text-green-500 flex-shrink-0" />
                      : prog.best_score > 0
                        ? <Clock size={20} className="text-yellow-500 flex-shrink-0" />
                        : null
                  ) : null}
                </div>
                {prog && (
                  <div className="mt-3 pt-3 border-t border-gray-100 dark:border-slate-700">
                    <div className="flex items-center justify-between text-xs text-gray-500">
                      <span>{prog.passed ? 'اجتزت الاختبار ✓' : `أفضل نتيجة: ${prog.best_score}%`}</span>
                    </div>
                    <div className="mt-1 h-1.5 bg-gray-100 dark:bg-slate-700 rounded-full overflow-hidden">
                      <div className={`h-full rounded-full transition-all ${prog.passed ? 'bg-green-500' : 'bg-secondary'}`}
                        style={{ width: `${Math.min(prog.best_score, 100)}%` }} />
                    </div>
                  </div>
                )}
              </Card>
            )
          })}
        </div>
      )}

      <LessonFormModal open={showForm} onClose={() => setShowForm(false)} lesson={editing} />
    </div>
  )
}
