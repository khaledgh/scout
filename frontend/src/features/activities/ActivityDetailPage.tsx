import { useState, useRef } from 'react'
import { useParams } from 'react-router-dom'
import { useActivity, useActivityAttendance, useRecordAttendance, useUploadActivityMedia } from '@/hooks/useActivities'
import { useMembers } from '@/hooks/useMembers'
import { useAuth } from '@/features/auth/AuthContext'
import { isLeaderRole } from '@/lib/permissions'
import { Card, Badge, Spinner, Button, Avatar } from '@/components/ui'
import { CalendarDays, MapPin, Users, CheckCircle, ImagePlus, Images, Play } from 'lucide-react'
import { format } from 'date-fns'
import { arLB as ar } from '@/lib/arLB'
import { assetUrl } from '@/lib/assetUrl'
import { toast } from 'sonner'

const statusLabel: Record<string, string> = {
  present: 'حاضر', absent: 'غائب', excused: 'معذور', late: 'متأخر',
}

export function ActivityDetailPage() {
  const { id } = useParams<{ id: string }>()
  const activityId = Number(id)
  const { user } = useAuth()
  const isLeader = isLeaderRole(user)
  const mediaInput = useRef<HTMLInputElement>(null)
  const { data: activity, isLoading } = useActivity(activityId)
  const { data: attendance } = useActivityAttendance(activityId)
  const { data: members } = useMembers({ page_size: 200 })
  const recordAttendance = useRecordAttendance(activityId)
  const uploadMedia = useUploadActivityMedia(activityId)

  const [localStatus, setLocalStatus] = useState<Record<number, string>>({})

  const onMediaSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    await uploadMedia.mutateAsync(file)
    toast.success('تم رفع الملف')
    e.target.value = ''
  }

  const handleStatusChange = (memberId: number, status: string) => {
    setLocalStatus(prev => ({ ...prev, [memberId]: status }))
  }

  const handleSaveAttendance = async () => {
    const records = Object.entries(localStatus).map(([mid, status]) => ({
      member_id: Number(mid),
      status,
    }))
    if (!records.length) { toast.error('لم يتم تعديل أي حضور'); return }
    await recordAttendance.mutateAsync(records)
    toast.success('تم تسجيل الحضور')
    setLocalStatus({})
  }

  if (isLoading) return <Spinner className="h-64" />
  if (!activity) return <div className="text-center py-16 text-gray-500">النشاط غير موجود</div>

  const presentCount = attendance?.filter(a => a.status === 'present').length ?? 0

  return (
    <div dir="rtl" className="space-y-6">
      {/* Header */}
      <Card>
        <div className="flex items-start gap-4">
          <div className="w-14 h-14 rounded-2xl bg-primary/10 flex items-center justify-center flex-shrink-0">
            <CalendarDays size={26} className="text-primary" />
          </div>
          <div className="flex-1">
            <h1 className="text-xl font-bold text-gray-900 dark:text-white">{activity.title}</h1>
            <div className="flex flex-wrap gap-3 mt-2 text-sm text-gray-500 dark:text-slate-400">
              <span className="flex items-center gap-1.5">
                <CalendarDays size={14} />
                {format(new Date(activity.starts_at), 'PPp', { locale: ar })}
              </span>
              {activity.location && (
                <span className="flex items-center gap-1.5">
                  <MapPin size={14} />{activity.location}
                </span>
              )}
              <span className="flex items-center gap-1.5">
                <Users size={14} />{presentCount} حاضر
              </span>
            </div>
            {activity.description && (
              <p className="text-sm text-gray-600 dark:text-slate-300 mt-2">{activity.description}</p>
            )}
          </div>
        </div>
      </Card>

      {/* Attendance */}
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h2 className="section-title mb-0">تسجيل الحضور</h2>
          <div className="flex gap-2">
            <Badge variant="green">{presentCount} حاضر</Badge>
            <Badge variant="red">{attendance?.filter(a => a.status === 'absent').length ?? 0} غائب</Badge>
          </div>
        </div>

        {/* Quick attendance from members list */}
        <div className="space-y-2">
          {members?.data?.map((m) => {
            const existing = attendance?.find(a => a.member_id === m.id)
            const current = localStatus[m.id] ?? existing?.status ?? ''
            return (
              <div key={m.id} className="flex items-center gap-3 p-2 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/30 transition-colors">
                <Avatar name={m.full_name} url={m.photo_url} size="sm" />
                <span className="flex-1 text-sm font-medium text-gray-900 dark:text-white">{m.full_name}</span>
                <div className="flex gap-1">
                  {['present', 'absent', 'excused', 'late'].map((s) => (
                    <button key={s} onClick={() => handleStatusChange(m.id, s)}
                      className={`px-2 py-1 rounded-lg text-xs font-medium transition-colors ${
                        current === s
                          ? s === 'present' ? 'bg-green-500 text-white' : s === 'absent' ? 'bg-red-500 text-white' : 'bg-blue-500 text-white'
                          : 'bg-gray-100 dark:bg-slate-700 text-gray-500 dark:text-slate-400 hover:bg-gray-200'
                      }`}>
                      {statusLabel[s]}
                    </button>
                  ))}
                </div>
              </div>
            )
          })}
        </div>

        {Object.keys(localStatus).length > 0 && (
          <div className="mt-4 flex justify-end">
            <Button onClick={handleSaveAttendance} loading={recordAttendance.isPending}>
              <CheckCircle size={16} /> حفظ الحضور
            </Button>
          </div>
        )}
      </Card>

      {/* Media */}
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
              <input ref={mediaInput} type="file" accept="image/*,video/mp4" className="hidden" onChange={onMediaSelect} />
            </>
          )}
        </div>
        {!activity.media?.length ? (
          <div className="text-center py-10">
            <Images size={36} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
            <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد صور أو وسائط</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
            {activity.media.map((m) =>
              m.media_type === 'image' ? (
                <a key={m.id} href={assetUrl(m.url)} target="_blank" rel="noopener noreferrer">
                  <img src={assetUrl(m.url)} alt={m.caption} className="w-full h-32 object-cover rounded-xl hover:opacity-90 transition-opacity" />
                </a>
              ) : (
                <div key={m.id} className="w-full h-32 rounded-xl bg-gray-100 dark:bg-slate-700 flex items-center justify-center">
                  <Play size={32} className="text-gray-400" />
                </div>
              )
            )}
          </div>
        )}
      </Card>
    </div>
  )
}
