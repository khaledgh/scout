import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMemberLeaderboard } from '@/hooks/useBadges'
import { useUnitLeaderboard } from '@/hooks/useUnits'
import { Card, Spinner, Avatar, SectionBadge, Select } from '@/components/ui'
import { assetUrl } from '@/lib/assetUrl'
import { Trophy, Users, Shield } from 'lucide-react'
import type { Section } from '@/types'

const sectionOptions: { value: string; label: string }[] = [
  { value: '',        label: 'جميع الشعب' },
  { value: 'ashbal',  label: 'أشبال' },
  { value: 'kashaf',  label: 'كشاف' },
  { value: 'jawala',  label: 'جوالة' },
  { value: 'mukashe', label: 'مكاشفة' },
]

const rankStyle = (i: number) =>
  i === 0 ? 'bg-yellow-400 text-yellow-900' :
  i === 1 ? 'bg-gray-300 text-gray-700' :
  i === 2 ? 'bg-amber-700 text-white' :
             'bg-gray-100 dark:bg-slate-700 text-gray-600 dark:text-slate-300'

export function LeaderboardPage() {
  const navigate = useNavigate()
  const [tab, setTab] = useState<'members' | 'units'>('members')
  const [section, setSection] = useState('')

  const { data: members, isLoading: membersLoading } = useMemberLeaderboard(section || undefined)
  const { data: units,   isLoading: unitsLoading   } = useUnitLeaderboard()

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header flex items-center gap-2">
          <Trophy size={24} className="text-secondary" /> لوحة المتصدرين
        </h1>
      </div>

      {/* Tabs */}
      <div className="flex gap-2">
        {[
          { key: 'members', label: 'أبطال الخبرة',  icon: Users },
          { key: 'units',   label: 'أفضل الطلائع',  icon: Shield },
        ].map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            onClick={() => setTab(key as 'members' | 'units')}
            className={`flex items-center gap-1.5 px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
              tab === key
                ? 'bg-primary text-white'
                : 'bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-slate-400'
            }`}
          >
            <Icon size={15} /> {label}
          </button>
        ))}
      </div>

      {/* Members leaderboard */}
      {tab === 'members' && (
        <Card>
          <div className="flex items-center justify-between mb-4">
            <h2 className="section-title mb-0">أبطال الخبرة</h2>
            <Select
              options={sectionOptions}
              value={section}
              onChange={setSection}
              className="w-36"
            />
          </div>
          {membersLoading ? (
            <Spinner className="h-32" />
          ) : !members?.length ? (
            <div className="text-center py-12">
              <Trophy size={40} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
              <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد بيانات</p>
            </div>
          ) : (
            <div className="space-y-2">
              {members.map((m: { id: number; full_name: string; photo_url?: string; xp_total: number; level: number; section: Section }, i: number) => (
                <div
                  key={m.id}
                  className="flex items-center gap-4 p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors cursor-pointer"
                  onClick={() => navigate(`/members/${m.id}`)}
                >
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm flex-shrink-0 ${rankStyle(i)}`}>
                    {i + 1}
                  </div>
                  <Avatar name={m.full_name} url={m.photo_url} size="sm" />
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold text-gray-900 dark:text-white truncate">{m.full_name}</p>
                    <SectionBadge section={m.section} />
                  </div>
                  <div className="text-left flex-shrink-0">
                    <p className="font-bold text-lg tabular-nums text-primary">{m.xp_total}</p>
                    <p className="text-xs text-gray-400">XP · مستوى {m.level}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}

      {/* Units leaderboard */}
      {tab === 'units' && (
        <Card>
          <h2 className="section-title mb-4">أفضل الطلائع</h2>
          {unitsLoading ? (
            <Spinner className="h-32" />
          ) : !units?.length ? (
            <div className="text-center py-12">
              <Shield size={40} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
              <p className="text-sm text-gray-400 dark:text-slate-500">لا توجد بيانات</p>
            </div>
          ) : (
            <div className="space-y-2">
              {units.map((u, i) => (
                <div
                  key={u.id}
                  className="flex items-center gap-4 p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors cursor-pointer"
                  onClick={() => navigate(`/units/${u.id}`)}
                >
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm flex-shrink-0 ${rankStyle(i)}`}>
                    {i + 1}
                  </div>
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary to-accent flex items-center justify-center text-white font-bold text-lg flex-shrink-0">
                    {u.emblem_url ? <img src={assetUrl(u.emblem_url)} className="w-full h-full object-cover rounded-xl" /> : u.name[0]}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold text-gray-900 dark:text-white truncate">{u.name}</p>
                    <SectionBadge section={u.section} />
                  </div>
                  <div className="text-left flex-shrink-0">
                    <p className="font-bold text-lg tabular-nums text-primary">{u.score_total}</p>
                    <p className="text-xs text-gray-400">نقطة · مستوى {u.level}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}
    </div>
  )
}
