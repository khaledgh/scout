import { useRef, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useMember, useMemberTimeline, useMemberMedical, useMemberQR, useUploadMemberPhoto } from '@/hooks/useMembers'
import { useMemberBadges, useRevokeBadge } from '@/hooks/useBadges'
import { useUnits, useAssignMemberToUnit } from '@/hooks/useUnits'
import { Avatar, Badge, Card, Spinner, SectionBadge, Button, Modal, Select } from '@/components/ui'
import {
  Award, Activity, Star, ClipboardList, Heart, Camera, Pencil, Plus,
  Trash2, ArrowRight, MapPin, Phone, User, Calendar, ShieldCheck, Shield,
} from 'lucide-react'
import { format } from 'date-fns'
import { arLB as ar } from '@/lib/arLB'
import { QRCodeSVG } from 'qrcode.react'
import { useAuth } from '@/features/auth/AuthContext'
import { canManageMembers, isLeaderRole } from '@/lib/permissions'
import { toast } from 'sonner'
import { MemberFormModal } from './MemberFormModal'
import { MedicalFormModal } from './MedicalFormModal'
import { EvaluationFormModal } from './EvaluationFormModal'
import { AwardBadgeModal } from './AwardBadgeModal'

const tabs = [
  { id: 'overview',    label: 'نظرة عامة',  icon: Star },
  { id: 'badges',      label: 'الشارات',     icon: Award },
  { id: 'activities',  label: 'الأنشطة',     icon: Activity },
  { id: 'evaluations', label: 'التقييمات',   icon: ClipboardList },
  { id: 'medical',     label: 'الطبي',       icon: Heart, restricted: true },
]

export function MemberProfilePage() {
  const { id } = useParams<{ id: string }>()
  const memberId   = Number(id)
  const navigate   = useNavigate()
  const { user }   = useAuth()
  const canManage  = canManageMembers(user)
  const [activeTab, setActiveTab] = useState('overview')
  const [showEdit,   setShowEdit]   = useState(false)
  const [showMedical, setShowMedical] = useState(false)
  const [showEval,   setShowEval]   = useState(false)
  const [showAward,  setShowAward]  = useState(false)
  const [showAssignUnit, setShowAssignUnit] = useState(false)
  const [selectedUnitId, setSelectedUnitId] = useState('')
  const fileInput = useRef<HTMLInputElement>(null)

  const { data: member, isLoading } = useMember(memberId)
  const { data: badges }            = useMemberBadges(memberId)
  const { data: timeline }          = useMemberTimeline(memberId)
  const { data: qr }                = useMemberQR(memberId)
  const { data: medical }           = useMemberMedical(canManage ? memberId : undefined)
  const { data: allUnits }          = useUnits()
  const uploadPhoto    = useUploadMemberPhoto(memberId)
  const revokeBadge    = useRevokeBadge(memberId)
  const assignToUnit   = useAssignMemberToUnit()
  const isLeader       = isLeaderRole(user)

  const onPhotoSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    await uploadPhoto.mutateAsync(file)
    toast.success('تم تحديث الصورة')
  }

  const handleRevoke = async (badgeId: number) => {
    if (!confirm('هل تريد سحب هذه الشارة؟')) return
    await revokeBadge.mutateAsync(badgeId)
    toast.success('تم سحب الشارة')
  }

  const handleAssignUnit = async () => {
    if (!selectedUnitId) return
    await assignToUnit.mutateAsync({ unitId: Number(selectedUnitId), memberId })
    toast.success('تم تعيين العضو للطليعة')
    setShowAssignUnit(false)
    setSelectedUnitId('')
  }

  if (isLoading) return <Spinner className="h-64" />
  if (!member)   return (
    <div dir="rtl" className="text-center py-20">
      <p className="text-gray-500 dark:text-slate-400 mb-4">العضو غير موجود</p>
      <Button variant="outline" onClick={() => navigate('/members')}>
        <ArrowRight size={16} /> العودة للأعضاء
      </Button>
    </div>
  )

  const visibleTabs = tabs.filter((t) => !t.restricted || canManage)
  const activitiesCount = timeline?.attendances?.length ?? 0
  const presentCount    = timeline?.attendances?.filter((a: { status: string }) => a.status === 'present').length ?? 0

  return (
    <div dir="rtl" className="space-y-5 max-w-5xl mx-auto">

      {/* ── Back button ────────────────────────────────────────────── */}
      <button
        onClick={() => navigate('/members')}
        className="flex items-center gap-2 text-sm text-gray-500 dark:text-slate-400 hover:text-primary transition-colors"
      >
        <ArrowRight size={16} /> الأعضاء
      </button>

      {/* ── Profile header card ──────────────────────────────────── */}
      <div className="bg-white dark:bg-slate-800 rounded-3xl shadow-card border border-gray-100 dark:border-slate-700 overflow-hidden">
        {/* Cover */}
        <div className="h-36 bg-gradient-to-l from-primary-800 via-primary to-accent relative">
          <div className="absolute inset-0 opacity-20" style={{ backgroundImage: 'radial-gradient(circle at 20% 50%, white 1px, transparent 1px)', backgroundSize: '24px 24px' }} />
        </div>

        {/* Body */}
        <div className="px-6 pb-6">
          {/* Avatar row */}
          <div className="flex flex-col sm:flex-row sm:items-end gap-4 -mt-14">
            <div className="relative flex-shrink-0">
              <Avatar
                name={member.full_name}
                url={member.photo_url}
                size="xl"
                className="border-4 border-white dark:border-slate-800 shadow-lg"
              />
              {canManage && (
                <button
                  onClick={() => fileInput.current?.click()}
                  className="absolute bottom-0 left-0 p-1.5 rounded-full bg-primary text-white shadow-md hover:bg-primary-600 transition-colors"
                  title="تغيير الصورة"
                >
                  <Camera size={13} />
                </button>
              )}
              <input ref={fileInput} type="file" accept="image/*" className="hidden" onChange={onPhotoSelect} />
            </div>

            <div className="flex-1 pt-16 sm:pt-0 sm:pb-1">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="text-2xl font-extrabold text-gray-900 dark:text-white">{member.full_name}</h1>
                <Badge variant={member.status === 'active' ? 'green' : 'gray'}>
                  {member.status === 'active' ? 'نشط' : 'غير نشط'}
                </Badge>
              </div>
              <div className="flex flex-wrap items-center gap-2 mt-1.5">
                <SectionBadge section={member.section} />
              </div>
            </div>
          </div>

          {/* ── Edit button row ────────────────────────────────────── */}
          {canManage && (
            <div className="flex justify-start mt-3">
              <Button variant="outline" size="sm" onClick={() => setShowEdit(true)}>
                <Pencil size={14} /> تعديل
              </Button>
            </div>
          )}

          {/* ── Stats row ──────────────────────────────────────────── */}
          <div className="grid grid-cols-2 sm:grid-cols-5 gap-3 mt-5 pt-5 border-t border-gray-100 dark:border-slate-700">
            {[
              { label: 'نقاط XP', value: member.xp_total,         color: 'text-primary' },
              { label: 'المستوى',  value: member.level,             color: 'text-accent' },
              { label: 'الشارات', value: badges?.length ?? 0,      color: 'text-secondary' },
              { label: 'الأنشطة', value: activitiesCount,           color: 'text-emerald-600' },
              { label: 'المرتبة',  value: member.rank_stage || '—', color: 'text-violet-600' },
            ].map(({ label, value, color }) => (
              <div key={label} className="text-center p-3 rounded-2xl bg-gray-50 dark:bg-slate-700/50">
                <p className={`text-2xl font-extrabold tabular-nums ${color}`}>{value}</p>
                <p className="text-xs text-gray-400 dark:text-slate-400 mt-0.5">{label}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* ── Tabs ───────────────────────────────────────────────────── */}
      <div className="flex gap-1 bg-white dark:bg-slate-800 border border-gray-100 dark:border-slate-700 p-1.5 rounded-2xl shadow-card overflow-x-auto">
        {visibleTabs.map(({ id, label, icon: Icon }) => (
          <button
            key={id}
            onClick={() => setActiveTab(id)}
            className={`flex items-center gap-1.5 px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-all flex-shrink-0 ${
              activeTab === id
                ? 'bg-primary text-white shadow-sm'
                : 'text-gray-500 dark:text-slate-400 hover:bg-gray-50 dark:hover:bg-slate-700 hover:text-gray-900 dark:hover:text-slate-100'
            }`}
          >
            <Icon size={15} />
            {label}
          </button>
        ))}
      </div>

      {/* ── Overview ───────────────────────────────────────────────── */}
      {activeTab === 'overview' && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-5">
          <Card className="lg:col-span-2">
            <h3 className="section-title">المعلومات الأساسية</h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {[
                { icon: Calendar,    label: 'تاريخ الميلاد',  value: member.birth_date ? format(new Date(member.birth_date), 'PPP', { locale: ar }) : '—' },
                { icon: Calendar,    label: 'تاريخ الانتساب', value: member.join_date ? format(new Date(member.join_date), 'PPP', { locale: ar }) : '—' },
                { icon: User,        label: 'اسم ولي الأمر', value: member.parent_name || '—' },
                { icon: Phone,       label: 'هاتف ولي الأمر', value: member.parent_phone || '—' },
                { icon: Phone,       label: 'هاتف إضافي',    value: member.secondary_phone || '—' },
                { icon: MapPin,      label: 'العنوان',        value: member.address || '—' },
                { icon: User,        label: 'الجنس',          value: member.gender === 'male' ? 'ذكر' : 'أنثى' },
                { icon: ShieldCheck, label: 'المرتبة',        value: member.rank_stage || '—' },
              ].map(({ icon: Icon, label, value }) => (
                <div key={label} className="flex items-start gap-3 p-3 rounded-xl bg-gray-50 dark:bg-slate-700/40">
                  <div className="w-8 h-8 rounded-lg bg-primary/10 dark:bg-primary/20 flex items-center justify-center flex-shrink-0">
                    <Icon size={15} className="text-primary" />
                  </div>
                  <div className="min-w-0">
                    <p className="text-xs text-gray-400 dark:text-slate-500">{label}</p>
                    <p className="text-sm font-semibold text-gray-900 dark:text-white mt-0.5 break-words">{value}</p>
                  </div>
                </div>
              ))}
            </div>
          </Card>

          <div className="space-y-5">
            {/* Unit card */}
            <Card>
              <div className="flex items-center justify-between mb-3">
                <h3 className="section-title mb-0 flex items-center gap-2">
                  <Shield size={16} className="text-primary" /> الطليعة
                </h3>
                {isLeader && (
                  <button
                    onClick={() => setShowAssignUnit(true)}
                    className="text-xs text-primary hover:underline"
                  >
                    تعيين
                  </button>
                )}
              </div>
              {member.units?.[0]?.unit ? (
                <div className="flex items-center gap-3 p-3 rounded-xl bg-primary/5 dark:bg-primary/10">
                  <div className="w-9 h-9 rounded-xl bg-primary/10 dark:bg-primary/20 flex items-center justify-center flex-shrink-0">
                    <Shield size={16} className="text-primary" />
                  </div>
                  <div className="min-w-0">
                    <p className="font-semibold text-sm text-gray-900 dark:text-white truncate">{member.units[0].unit.name}</p>
                    <SectionBadge section={member.units[0].unit.section} />
                  </div>
                </div>
              ) : (
                <p className="text-sm text-gray-400 dark:text-slate-500 text-center py-2">غير مُعيَّن</p>
              )}
            </Card>

            {/* QR card */}
            <Card>
              <h3 className="section-title">رمز QR</h3>
              {qr ? (
                <div className="flex flex-col items-center gap-3">
                  <div className="p-3 bg-white rounded-xl border border-gray-100 dark:border-slate-700">
                    <QRCodeSVG value={qr.token} size={130} />
                  </div>
                  <p className="text-xs text-gray-400 dark:text-slate-500 text-center">امسح للتحقق من الحضور</p>
                </div>
              ) : (
                <p className="text-sm text-gray-400">لا يوجد رمز QR</p>
              )}
            </Card>

            {/* Attendance summary */}
            {activitiesCount > 0 && (
              <Card>
                <h3 className="section-title">ملخص الحضور</h3>
                <div className="space-y-3">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-500 dark:text-slate-400">نسبة الحضور</span>
                    <span className="font-bold text-gray-900 dark:text-white">
                      {activitiesCount > 0 ? Math.round((presentCount / activitiesCount) * 100) : 0}%
                    </span>
                  </div>
                  <div className="h-2 bg-gray-100 dark:bg-slate-700 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-gradient-to-l from-emerald-400 to-green-500 rounded-full transition-all"
                      style={{ width: `${activitiesCount > 0 ? Math.round((presentCount / activitiesCount) * 100) : 0}%` }}
                    />
                  </div>
                  <p className="text-xs text-gray-400 dark:text-slate-500">{presentCount} حاضر من {activitiesCount} نشاط</p>
                </div>
              </Card>
            )}
          </div>
        </div>
      )}

      {/* ── Badges ─────────────────────────────────────────────────── */}
      {activeTab === 'badges' && (
        <Card>
          <div className="flex items-center justify-between mb-5">
            <h3 className="section-title mb-0">الشارات المكتسبة ({badges?.length ?? 0})</h3>
            {canManage && (
              <Button size="sm" onClick={() => setShowAward(true)}><Plus size={14} /> منح شارة</Button>
            )}
          </div>
          {!badges?.length ? (
            <div className="text-center py-10">
              <Award size={40} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
              <p className="text-sm text-gray-400 dark:text-slate-500">لم يكتسب هذا العضو أي شارة بعد</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
              {badges.map((mb) => (
                <div key={mb.id} className="relative flex flex-col items-center gap-2 p-4 rounded-2xl bg-gradient-to-b from-secondary/5 to-transparent border border-secondary/15 dark:border-secondary/20 text-center">
                  {canManage && (
                    <button
                      onClick={() => handleRevoke(mb.badge_id)}
                      className="absolute top-2 left-2 p-1 rounded-lg text-gray-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                      title="سحب الشارة"
                    >
                      <Trash2 size={13} />
                    </button>
                  )}
                  <div className="w-14 h-14 rounded-full bg-secondary/15 dark:bg-secondary/25 flex items-center justify-center text-2xl shadow-sm">🏅</div>
                  <p className="text-sm font-semibold text-gray-900 dark:text-white leading-tight">{mb.badge?.name}</p>
                  <p className="text-xs text-gray-400">{mb.badge?.category}</p>
                  {(mb.badge?.xp_reward ?? 0) > 0 && (
                    <span className="badge bg-secondary/10 text-secondary text-xs">+{mb.badge?.xp_reward} XP</span>
                  )}
                </div>
              ))}
            </div>
          )}
        </Card>
      )}

      {/* ── Activities ─────────────────────────────────────────────── */}
      {activeTab === 'activities' && (
        <Card>
          <h3 className="section-title">سجل الأنشطة</h3>
          {!timeline?.attendances?.length ? (
            <p className="text-sm text-gray-400 dark:text-slate-500">لا سجلات حضور</p>
          ) : (
            <div className="divide-y divide-gray-50 dark:divide-slate-700">
              {timeline.attendances.map((a: { id: number; activity?: { title?: string; starts_at?: string }; status: string }) => (
                <div key={a.id} className="flex items-center justify-between py-3">
                  <div className="flex items-center gap-3 min-w-0">
                    <div className={`w-2 h-2 rounded-full flex-shrink-0 ${a.status === 'present' ? 'bg-emerald-500' : a.status === 'absent' ? 'bg-red-500' : 'bg-amber-500'}`} />
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-gray-900 dark:text-white truncate">{a.activity?.title}</p>
                      {a.activity?.starts_at && (
                        <p className="text-xs text-gray-400 dark:text-slate-500">
                          {format(new Date(a.activity.starts_at), 'PPP', { locale: ar })}
                        </p>
                      )}
                    </div>
                  </div>
                  <Badge variant={a.status === 'present' ? 'green' : a.status === 'absent' ? 'red' : 'gray'}>
                    {a.status === 'present' ? 'حاضر' : a.status === 'absent' ? 'غائب' : a.status === 'excused' ? 'معذور' : 'متأخر'}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}

      {/* ── Evaluations ────────────────────────────────────────────── */}
      {activeTab === 'evaluations' && (
        <Card>
          <div className="flex items-center justify-between mb-5">
            <h3 className="section-title mb-0">التقييمات</h3>
            {canManage && (
              <Button size="sm" onClick={() => setShowEval(true)}><Plus size={14} /> تقييم جديد</Button>
            )}
          </div>
          {!timeline?.evaluations?.length ? (
            <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد تقييمات</p>
          ) : (
            <div className="space-y-4">
              {timeline.evaluations.map((ev: { id: number; period: string; overall: number; discipline: number; participation: number; leadership: number; skill: number; notes?: string }) => (
                <div key={ev.id} className="p-4 rounded-2xl border border-gray-100 dark:border-slate-700 bg-gray-50/50 dark:bg-slate-700/30">
                  <div className="flex items-center justify-between mb-4">
                    <Badge variant="blue">{ev.period}</Badge>
                    <div className="flex items-center gap-1">
                      <span className="text-2xl font-extrabold text-primary tabular-nums">{ev.overall}</span>
                      <span className="text-gray-400 text-sm">/10</span>
                    </div>
                  </div>
                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                    {[['الانضباط', ev.discipline], ['المشاركة', ev.participation], ['القيادة', ev.leadership], ['المهارة', ev.skill]].map(([label, val]) => (
                      <div key={label as string} className="text-center">
                        <div className="text-lg font-bold text-gray-900 dark:text-white tabular-nums">{val}<span className="text-xs text-gray-400 font-normal">/10</span></div>
                        <p className="text-xs text-gray-400 dark:text-slate-500 mt-0.5">{label}</p>
                        <div className="mt-1.5 h-1.5 bg-gray-200 dark:bg-slate-600 rounded-full overflow-hidden">
                          <div className="h-full bg-primary/60 rounded-full" style={{ width: `${Number(val) * 10}%` }} />
                        </div>
                      </div>
                    ))}
                  </div>
                  {ev.notes && <p className="text-xs text-gray-500 dark:text-slate-400 mt-3 italic border-t border-gray-100 dark:border-slate-700 pt-3">{ev.notes}</p>}
                </div>
              ))}
            </div>
          )}
        </Card>
      )}

      {/* ── Medical ────────────────────────────────────────────────── */}
      {activeTab === 'medical' && canManage && (
        <Card>
          <div className="flex items-center justify-between mb-5">
            <h3 className="section-title mb-0 flex items-center gap-2">
              <Heart size={18} className="text-red-500" /> المعلومات الطبية
            </h3>
            <Button size="sm" variant="outline" onClick={() => setShowMedical(true)}>
              <Pencil size={14} /> {medical ? 'تعديل' : 'إضافة'}
            </Button>
          </div>
          {!medical ? (
            <div className="text-center py-10">
              <Heart size={40} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
              <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد معلومات طبية مسجلة</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {[
                { label: 'فصيلة الدم',    value: medical.blood_type },
                { label: 'الحساسية',       value: medical.allergies },
                { label: 'أمراض مزمنة',   value: medical.chronic_conditions },
                { label: 'الأدوية',        value: medical.medications },
                { label: 'ملاحظات الطوارئ', value: medical.emergency_notes },
              ].map(({ label, value }) => (
                <div key={label} className="p-3 rounded-xl bg-red-50/50 dark:bg-red-900/10 border border-red-100 dark:border-red-900/20">
                  <dt className="text-xs text-red-400 dark:text-red-400/70">{label}</dt>
                  <dd className="font-semibold text-gray-900 dark:text-white mt-0.5 text-sm">{value || '—'}</dd>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}

      {/* ── Modals ─────────────────────────────────────────────────── */}
      <MemberFormModal   open={showEdit}    onClose={() => setShowEdit(false)}    member={member} />
      <MedicalFormModal  open={showMedical} onClose={() => setShowMedical(false)} memberId={memberId} medical={medical} />
      <EvaluationFormModal open={showEval}  onClose={() => setShowEval(false)}    memberId={memberId} />
      <AwardBadgeModal   open={showAward}   onClose={() => setShowAward(false)}   memberId={memberId} earned={badges} />

      <Modal open={showAssignUnit} onClose={() => setShowAssignUnit(false)} title="تعيين لطليعة" size="sm">
        <div dir="rtl" className="space-y-4">
          <div>
            <label className="label">اختر الطليعة</label>
            <Select
              value={selectedUnitId}
              onChange={setSelectedUnitId}
              placeholder="— اختر طليعة —"
              options={(allUnits ?? []).map(u => ({ value: String(u.id), label: u.name }))}
            />
          </div>
          <div className="flex gap-3 justify-end">
            <Button variant="ghost" onClick={() => setShowAssignUnit(false)}>إلغاء</Button>
            <Button onClick={handleAssignUnit} loading={assignToUnit.isPending} disabled={!selectedUnitId}>
              تعيين
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
