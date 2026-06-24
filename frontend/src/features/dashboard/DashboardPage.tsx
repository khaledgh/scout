import { Users, CalendarDays, TrendingUp, Trophy, AlertTriangle, Award, Compass, Star, Activity } from 'lucide-react'
import { useDashboard } from '@/hooks/useDashboard'
import { useAuth } from '@/features/auth/AuthContext'
import { StatCard, Card, Avatar, Spinner, SectionBadge, XPRing } from '@/components/ui'
import { format } from 'date-fns'
import { arLB as ar } from '@/lib/arLB'
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts'

const activityTypeLabel: Record<string, string> = {
  camp: 'مخيم', hike: 'مسير', training: 'تدريب', meeting: 'اجتماع', service: 'خدمة',
}
const sectionLabel: Record<string, string> = {
  ashbal: 'أشبال', kashaf: 'كشاف', jawala: 'جوالة', mukashe: 'مكاشفة',
}
const sectionColors: Record<string, string> = {
  ashbal: '#7C3AED', kashaf: '#4F46E5', jawala: '#D97706', mukashe: '#DB2777',
}

export function DashboardPage() {
  const { data, isLoading } = useDashboard()
  const { user } = useAuth()

  if (isLoading) return <Spinner className="h-64" />

  const sectionData = (data?.members_by_section ?? []).map((s) => ({
    name: sectionLabel[s.section] ?? s.section,
    value: Number(s.count),
    key: s.section,
  }))

  return (
    <div dir="rtl" className="space-y-6">

      {/* ── Hero ─────────────────────────────────────────────────────── */}
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-l from-primary-900 via-primary-700 to-accent text-white p-6 sm:p-8 shadow-card-hover">
        <div className="absolute -left-10 -bottom-10 opacity-10 pointer-events-none">
          <Compass size={220} />
        </div>
        <div className="relative z-10">
          <p className="text-sm font-medium text-white/60 tracking-wide uppercase">إدارة الفوج الكشفي</p>
          <h1 className="text-2xl sm:text-3xl font-extrabold mt-1">
            أهلاً، {user?.full_name?.split(' ')[0] ?? 'القائد'} 👋
          </h1>
          <p className="text-white/70 mt-2 text-sm max-w-md">
            نظرة شاملة على نشاط الفوج، الحضور، وإنجازات الأعضاء.
          </p>
        </div>
      </div>

      {/* ── KPI cards ─────────────────────────────────────────────────── */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          label="إجمالي الأعضاء"
          value={data?.member_count ?? 0}
          subtitle={`${data?.active_members ?? 0} نشط`}
          icon={Users}
          color="purple"
        />
        <StatCard
          label="الأعضاء النشطون"
          value={data?.active_members ?? 0}
          icon={TrendingUp}
          color="green"
        />
        <StatCard
          label="نسبة الحضور"
          value={`${(data?.attendance_rate ?? 0).toFixed(1)}%`}
          icon={CalendarDays}
          color="teal"
        />
        <StatCard
          label="أفضل طليعة"
          value={data?.top_unit?.name ?? '—'}
          subtitle={data?.top_unit ? `${data.top_unit.score_total} نقطة` : undefined}
          icon={Trophy}
          color="gold"
        />
      </div>

      {/* ── Charts row ────────────────────────────────────────────────── */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

        {/* Section donut */}
        <Card>
          <h2 className="section-title">توزّع الأعضاء</h2>
          {!sectionData.length ? (
            <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد بيانات</p>
          ) : (
            <ResponsiveContainer width="100%" height={220}>
              <PieChart>
                <Pie data={sectionData} dataKey="value" nameKey="name" cx="50%" cy="50%"
                  innerRadius={52} outerRadius={80} paddingAngle={4} strokeWidth={0}>
                  {sectionData.map((entry) => (
                    <Cell key={entry.key} fill={sectionColors[entry.key] ?? '#8B5CF6'} />
                  ))}
                </Pie>
                <Tooltip formatter={(v) => [`${v} عضو`]} />
                <Legend wrapperStyle={{ fontSize: 12 }} />
              </PieChart>
            </ResponsiveContainer>
          )}
        </Card>

        {/* Top XP leaderboard */}
        <Card className="lg:col-span-2">
          <h2 className="section-title flex items-center gap-2">
            <Star size={18} className="text-secondary" /> أبطال نقاط الخبرة
          </h2>
          {!data?.top_members?.length ? (
            <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد بيانات</p>
          ) : (
            <div className="space-y-1.5">
              {data.top_members.map((m, i) => (
                <div key={m.id} className="flex items-center gap-3 p-2.5 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                  <div className="w-8 h-8 flex items-center justify-center text-base flex-shrink-0">
                    {i === 0 ? '🥇' : i === 1 ? '🥈' : i === 2 ? '🥉' : (
                      <span className="text-xs font-bold text-gray-400 bg-gray-100 dark:bg-slate-700 rounded-full w-7 h-7 flex items-center justify-center">{i + 1}</span>
                    )}
                  </div>
                  <Avatar name={m.full_name} url={m.photo_url} size="sm" />
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold text-sm text-gray-900 dark:text-white truncate">{m.full_name}</p>
                    <SectionBadge section={m.section} />
                  </div>
                  <XPRing xp={m.xp_total} level={m.level} size={42} />
                  <div className="text-left w-14 flex-shrink-0">
                    <p className="font-extrabold text-primary tabular-nums text-sm">{m.xp_total}</p>
                    <p className="text-[10px] text-gray-400">XP</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>

      {/* ── Activities + At-risk ──────────────────────────────────────── */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">

        <Card>
          <h2 className="section-title flex items-center gap-2">
            <Activity size={18} className="text-accent" /> الأنشطة القادمة
          </h2>
          {!data?.upcoming_activities?.length ? (
            <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد أنشطة قادمة</p>
          ) : (
            <div className="space-y-2">
              {data.upcoming_activities.map((a) => (
                <div key={a.id} className="flex items-center gap-3 p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                  <div className="w-10 h-10 rounded-xl bg-accent/10 dark:bg-accent/20 flex items-center justify-center flex-shrink-0">
                    <CalendarDays size={18} className="text-accent" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm text-gray-900 dark:text-white truncate">{a.title}</p>
                    <p className="text-xs text-gray-400 dark:text-slate-500 mt-0.5">
                      {format(new Date(a.starts_at), 'PPP', { locale: ar })} · {activityTypeLabel[a.type] ?? a.type}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>

        <Card>
          <h2 className="section-title flex items-center gap-2">
            <AlertTriangle size={18} className="text-red-500" /> أعضاء يحتاجون متابعة
          </h2>
          {!data?.at_risk_members?.length ? (
            <p className="text-sm text-gray-400 dark:text-slate-500">لا يوجد أعضاء بحاجة لمتابعة</p>
          ) : (
            <div className="space-y-2">
              {data.at_risk_members.map((m) => (
                <div key={m.id} className="flex items-center gap-3 p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                  <Avatar name={m.full_name} url={m.photo_url} size="sm" />
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm text-gray-900 dark:text-white truncate">{m.full_name}</p>
                    <SectionBadge section={m.section} />
                  </div>
                  <span className="badge-red text-xs shrink-0">غياب متكرر</span>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>

      {/* ── Recent badges ────────────────────────────────────────────── */}
      {(data?.recent_badges?.length ?? 0) > 0 && (
        <Card>
          <h2 className="section-title flex items-center gap-2">
            <Award size={18} className="text-secondary" /> آخر الشارات الممنوحة
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            {data!.recent_badges.map((mb) => (
              <div key={mb.id} className="flex items-center gap-3 p-3.5 rounded-xl bg-gradient-to-l from-secondary/5 to-transparent border border-secondary/10 dark:border-secondary/20">
                <div className="w-11 h-11 rounded-full bg-secondary/15 dark:bg-secondary/25 flex items-center justify-center text-xl flex-shrink-0">🏅</div>
                <div className="min-w-0">
                  <p className="font-semibold text-sm text-gray-900 dark:text-white truncate">{mb.badge?.name}</p>
                  <p className="text-xs text-gray-400 dark:text-slate-500 truncate">{mb.member?.full_name}</p>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  )
}
