import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, Users, Trophy, Pencil, Trash2 } from 'lucide-react'
import { useUnits, useUnitLeaderboard, useDeleteUnit } from '@/hooks/useUnits'
import { assetUrl } from '@/lib/assetUrl'
import { Button, Card, Spinner, EmptyState, SectionBadge } from '@/components/ui'
import { useAuth } from '@/features/auth/AuthContext'
import { isAdmin as isAdminRole } from '@/lib/permissions'
import { toast } from 'sonner'
import type { Unit } from '@/types'
import { UnitFormModal } from './UnitFormModal'

export function UnitsPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Unit | undefined>(undefined)
  const [tab, setTab] = useState<'grid' | 'leaderboard'>('grid')

  const { data: units, isLoading } = useUnits()
  const { data: leaderboard } = useUnitLeaderboard()
  const deleteUnit = useDeleteUnit()

  const handleDelete = async (unit: Unit) => {
    if (!confirm(`حذف ${unit.name}؟`)) return
    await deleteUnit.mutateAsync(unit.id)
    toast.success('تم حذف الطليعة')
  }

  const openCreate = () => { setEditing(undefined); setShowForm(true) }
  const openEdit = (unit: Unit) => { setEditing(unit); setShowForm(true) }

  const isAdmin = isAdminRole(user)

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">الطلائع</h1>
        {isAdmin && <Button onClick={openCreate}><Plus size={16} />طليعة جديدة</Button>}
      </div>

      <div className="flex gap-2">
        {[{ key: 'grid', label: 'عرض الشبكة' }, { key: 'leaderboard', label: 'لوحة المتصدرين' }].map(({ key, label }) => (
          <button key={key} onClick={() => setTab(key as 'grid' | 'leaderboard')}
            className={`px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
              tab === key ? 'bg-primary text-white' : 'bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-slate-400'
            }`}>{label}</button>
        ))}
      </div>

      {isLoading ? <Spinner className="h-48" /> : tab === 'grid' ? (
        !units?.length ? (
          <EmptyState icon={Users} title="لا توجد طلائع" description="أنشئ أول طليعة للفوج" />
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {units.map((unit) => (
              <Card key={unit.id} className="cursor-pointer hover:shadow-card-hover transition-shadow"
                onClick={() => navigate(`/units/${unit.id}`)}>
                <div className="flex items-start gap-4">
                  <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-primary to-accent flex items-center justify-center text-white text-2xl font-bold flex-shrink-0">
                    {unit.emblem_url ? <img src={assetUrl(unit.emblem_url)} className="w-full h-full object-cover rounded-2xl" /> : unit.name[0]}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-bold text-gray-900 dark:text-white truncate">{unit.name}</h3>
                    <SectionBadge section={unit.section} />
                    {unit.motto && <p className="text-xs text-gray-500 dark:text-slate-400 mt-1 italic truncate">"{unit.motto}"</p>}
                  </div>
                  {isAdmin && (
                    <div className="flex items-center gap-1 flex-shrink-0" onClick={(e) => e.stopPropagation()}>
                      <button onClick={() => openEdit(unit)} title="تعديل"
                        className="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                        <Pencil size={15} />
                      </button>
                      <button onClick={() => handleDelete(unit)} title="حذف"
                        className="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                        <Trash2 size={15} />
                      </button>
                    </div>
                  )}
                </div>
                <div className="flex items-center gap-4 mt-4 pt-4 border-t border-gray-100 dark:border-slate-700 text-center">
                  <div className="flex-1">
                    <p className="text-xl font-bold text-primary tabular-nums">{unit.score_total}</p>
                    <p className="text-xs text-gray-500">نقاط</p>
                  </div>
                  <div className="flex-1">
                    <p className="text-xl font-bold text-accent tabular-nums">{unit.members?.length ?? 0}</p>
                    <p className="text-xs text-gray-500">عضو</p>
                  </div>
                  <div className="flex-1">
                    <p className="text-xl font-bold text-secondary tabular-nums">{unit.level}</p>
                    <p className="text-xs text-gray-500">مستوى</p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )
      ) : (
        <Card>
          <h3 className="section-title flex items-center gap-2"><Trophy size={18} className="text-secondary" />ترتيب الطلائع</h3>
          <div className="space-y-2">
            {leaderboard?.map((unit, i) => (
              <div key={unit.id} className="flex items-center gap-4 p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors cursor-pointer"
                onClick={() => navigate(`/units/${unit.id}`)}>
                <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm flex-shrink-0 ${
                  i === 0 ? 'bg-yellow-400 text-yellow-900' : i === 1 ? 'bg-gray-300 text-gray-700' : i === 2 ? 'bg-amber-700 text-white' : 'bg-gray-100 dark:bg-slate-700 text-gray-600 dark:text-slate-300'
                }`}>{i + 1}</div>
                <div className="flex-1 min-w-0">
                  <p className="font-semibold text-gray-900 dark:text-white">{unit.name}</p>
                  <SectionBadge section={unit.section} />
                </div>
                <div className="text-left">
                  <p className="font-bold text-lg tabular-nums text-primary">{unit.score_total}</p>
                  <p className="text-xs text-gray-400">نقطة</p>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}

      <UnitFormModal open={showForm} onClose={() => setShowForm(false)} unit={editing} />
    </div>
  )
}
