import { useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useTrainingLesson, useLessonQuiz, useUploadLessonCover, useUploadLessonMedia, useDeleteLessonMedia } from '@/hooks/useTraining'
import { useAuth } from '@/features/auth/AuthContext'
import { isLeaderRole } from '@/lib/permissions'
import { Card, Spinner, Button, Badge } from '@/components/ui'
import { ArrowRight, BookOpen, ImagePlus, CheckCircle, Images, Play, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { assetUrl } from '@/lib/assetUrl'

const categoryLabel: Record<string, string> = {
  scout: 'كشفي',
  first_aid: 'إسعاف أولي',
  nature: 'طبيعة',
  leadership: 'قيادة',
}

export function TrainingLessonDetailPage() {
  const { id } = useParams<{ id: string }>()
  const lessonId = Number(id)
  const navigate = useNavigate()
  const { user } = useAuth()
  const isLeader = isLeaderRole(user)
  const fileInput = useRef<HTMLInputElement>(null)
  const mediaInput = useRef<HTMLInputElement>(null)

  const { data: lesson, isLoading } = useTrainingLesson(lessonId)
  const { data: quiz } = useLessonQuiz(lessonId)
  const uploadCover = useUploadLessonCover(lessonId)
  const uploadMedia = useUploadLessonMedia(lessonId)
  const deleteMedia = useDeleteLessonMedia(lessonId)

  const onCoverSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    await uploadCover.mutateAsync(file)
    toast.success('تم رفع صورة الغلاف')
    e.target.value = ''
  }

  const onMediaSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    await uploadMedia.mutateAsync(file)
    toast.success('تم رفع الملف')
    e.target.value = ''
  }

  const handleDeleteMedia = async (mediaId: number) => {
    if (!confirm('حذف هذا الملف؟')) return
    await deleteMedia.mutateAsync(mediaId)
    toast.success('تم حذف الملف')
  }

  if (isLoading) return <Spinner className="h-64" />
  if (!lesson) return (
    <div dir="rtl" className="text-center py-20">
      <p className="text-gray-500 dark:text-slate-400 mb-4">الدرس غير موجود</p>
      <Button variant="outline" onClick={() => navigate('/training')}>
        <ArrowRight size={16} /> العودة للتدريب
      </Button>
    </div>
  )

  return (
    <div dir="rtl" className="space-y-5 max-w-4xl mx-auto">
      <button
        onClick={() => navigate('/training')}
        className="flex items-center gap-2 text-sm text-gray-500 dark:text-slate-400 hover:text-primary transition-colors"
      >
        <ArrowRight size={16} /> التدريب
      </button>

      {/* Cover / Header */}
      <div className="bg-white dark:bg-slate-800 rounded-3xl shadow-card border border-gray-100 dark:border-slate-700 overflow-hidden">
        {lesson.cover_url ? (
          <img src={assetUrl(lesson.cover_url)} alt={lesson.title} className="w-full h-52 object-cover" />
        ) : (
          <div className="h-40 bg-gradient-to-l from-primary-800 via-primary to-accent flex items-center justify-center">
            <BookOpen size={48} className="text-white/40" />
          </div>
        )}
        <div className="p-6">
          <div className="flex items-start gap-3 flex-wrap">
            <div className="flex-1 min-w-0">
              <h1 className="text-2xl font-extrabold text-gray-900 dark:text-white">{lesson.title}</h1>
              {lesson.category && (
                <Badge variant="blue" className="mt-2">
                  {categoryLabel[lesson.category] ?? lesson.category}
                </Badge>
              )}
            </div>
            {isLeader && (
              <>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => fileInput.current?.click()}
                  loading={uploadCover.isPending}
                >
                  <ImagePlus size={14} /> رفع صورة الغلاف
                </Button>
                <input
                  ref={fileInput}
                  type="file"
                  accept="image/*"
                  className="hidden"
                  onChange={onCoverSelect}
                />
              </>
            )}
          </div>
        </div>
      </div>

      {/* Content */}
      {lesson.content && (
        <Card>
          <h2 className="section-title flex items-center gap-2">
            <BookOpen size={18} className="text-primary" /> محتوى الدرس
          </h2>
          <div className="text-sm text-gray-700 dark:text-slate-300 whitespace-pre-wrap leading-relaxed mt-3">
            {lesson.content}
          </div>
        </Card>
      )}

      {/* Quiz */}
      <Card>
        <h2 className="section-title flex items-center gap-2 mb-4">
          <CheckCircle size={18} className="text-emerald-500" /> الاختبار
        </h2>
        {quiz ? (
          <div className="p-4 rounded-2xl bg-emerald-50 dark:bg-emerald-900/10 border border-emerald-100 dark:border-emerald-900/20">
            <p className="font-semibold text-gray-900 dark:text-white">{quiz.title}</p>
            <div className="flex flex-wrap gap-4 mt-2 text-sm text-gray-500 dark:text-slate-400">
              <span>{quiz.questions?.length ?? 0} سؤال</span>
              <span>درجة النجاح: {quiz.pass_score}%</span>
              {quiz.xp_reward > 0 && (
                <span className="font-semibold text-secondary">+{quiz.xp_reward} XP</span>
              )}
            </div>
            <Button className="mt-3" size="sm" disabled>
              ابدأ الاختبار (قريباً)
            </Button>
          </div>
        ) : (
          <div className="text-center py-8">
            <CheckCircle size={36} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
            <p className="text-sm text-gray-400 dark:text-slate-500">لا يوجد اختبار لهذا الدرس</p>
          </div>
        )}
      </Card>

      {/* Media Gallery */}
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h2 className="section-title mb-0 flex items-center gap-2">
            <Images size={18} className="text-primary" /> الصور والوسائط
          </h2>
          {isLeader && (
            <>
              <Button size="sm" variant="outline" onClick={() => mediaInput.current?.click()} loading={uploadMedia.isPending}>
                <ImagePlus size={14} /> إضافة صورة/فيديو
              </Button>
              <input
                ref={mediaInput}
                type="file"
                accept="image/*,video/*"
                className="hidden"
                onChange={onMediaSelect}
              />
            </>
          )}
        </div>
        {!lesson.media?.length ? (
          <div className="text-center py-10">
            <Images size={36} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
            <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد صور أو مقاطع مرفوعة</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
            {lesson.media.map((m) => (
              <div key={m.id} className="relative group rounded-2xl overflow-hidden aspect-video bg-gray-100 dark:bg-slate-700">
                {m.media_type === 'image' ? (
                  <a href={assetUrl(m.url)} target="_blank" rel="noopener noreferrer" className="block w-full h-full">
                    <img src={assetUrl(m.url)} alt={m.caption ?? ''} className="w-full h-full object-cover" />
                  </a>
                ) : (
                  <a href={assetUrl(m.url)} target="_blank" rel="noopener noreferrer"
                    className="flex items-center justify-center w-full h-full bg-gray-900">
                    <Play size={32} className="text-white/80" />
                  </a>
                )}
                {isLeader && (
                  <button
                    onClick={() => handleDeleteMedia(m.id)}
                    className="absolute top-1.5 left-1.5 p-1.5 rounded-lg bg-black/50 text-white opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-600"
                    title="حذف"
                  >
                    <Trash2 size={13} />
                  </button>
                )}
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  )
}
