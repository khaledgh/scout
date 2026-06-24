import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useUnit, useAddUnitMembers, useRemoveUnitMember, useDeleteUnit } from '@/hooks/useUnits'
import { useMembers } from '@/hooks/useMembers'
import { useAuth } from '@/features/auth/AuthContext'
import { isLeaderRole, isAdmin as isAdminRole } from '@/lib/permissions'
import { Card, Spinner, Button, Avatar, Badge, SectionBadge, Modal } from '@/components/ui'
import { Users, Trash2, Plus, ArrowRight, Shield, Pencil } from 'lucide-react'
import { toast } from 'sonner'
import { UnitFormModal } from './UnitFormModal'

export function UnitDetailPage() {
  const { id } = useParams<{ id: string }>()
  const unitId = Number(id)
  const navigate = useNavigate()
  const { user } = useAuth()

  const [showAddModal, setShowAddModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [selectedIds, setSelectedIds] = useState<number[]>([])

  const { data: unit, isLoading } = useUnit(unitId)
  const { data: allMembersData } = useMembers({ page_size: 200 })
  const addMembers = useAddUnitMembers(unitId)
  const removeMember = useRemoveUnitMember(unitId)
  const deleteUnit = useDeleteUnit()
  const isLeader = isLeaderRole(user)
  const isAdmin = isAdminRole(user)

  const currentMemberIds = new Set(unit?.members?.map((m) => m.member_id) ?? [])
  const availableMembers = (allMembersData?.data ?? []).filter((m) => !currentMemberIds.has(m.id))

  const handleRemove = async (memberId: number, name: string) => {
    if (!confirm(`إزالة ${name} من الطليعة؟`)) return
    await removeMember.mutateAsync(memberId)
    toast.success('تم إزالة العضو')
  }

  const handleDeleteUnit = async () => {
    if (!unit || !confirm(`حذف طليعة ${unit.name}؟ لا يمكن التراجع عن هذا الإجراء.`)) return
    await deleteUnit.mutateAsync(unit.id)
    toast.success('تم حذف الطليعة')
    navigate('/units')
  }

  const handleAdd = async () => {
    if (!selectedIds.length) return
    await addMembers.mutateAsync(selectedIds)
    toast.success(`تم إضافة ${selectedIds.length} عضو`)
    setSelectedIds([])
    setShowAddModal(false)
  }

  if (isLoading) return <Spinner className="h-64" />
  if (!unit) return (
    <div dir="rtl" className="text-center py-20">
      <p className="text-gray-500 dark:text-slate-400 mb-4">الطليعة غير موجودة</p>
      <Button variant="outline" onClick={() => navigate('/units')}>
        <ArrowRight size={16} /> العودة للطلائع
      </Button>
    </div>
  )

  return (
    <div dir="rtl" className="space-y-5 max-w-4xl mx-auto">
      <button
        onClick={() => navigate('/units')}
        className="flex items-center gap-2 text-sm text-gray-500 dark:text-slate-400 hover:text-primary transition-colors"
      >
        <ArrowRight size={16} /> الطلائع
      </button>

      {/* Header */}
      <Card>
        <div className="flex items-start gap-5">
          <div className="w-16 h-16 rounded-2xl bg-primary/10 dark:bg-primary/20 flex items-center justify-center flex-shrink-0">
            <Shield size={28} className="text-primary" />
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-3 flex-wrap">
              <h1 className="text-2xl font-extrabold text-gray-900 dark:text-white">{unit.name}</h1>
              <SectionBadge section={unit.section} />
              {!unit.is_active && <Badge variant="gray">غير نشطة</Badge>}
              {isAdmin && (
                <div className="flex items-center gap-1 mr-auto">
                  <button
                    onClick={() => setShowEditModal(true)}
                    className="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors"
                    title="تعديل الطليعة"
                  >
                    <Pencil size={15} />
                  </button>
                  <button
                    onClick={handleDeleteUnit}
                    className="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                    title="حذف الطليعة"
                  >
                    <Trash2 size={15} />
                  </button>
                </div>
              )}
            </div>
            {unit.motto && (
              <p className="text-sm text-gray-500 dark:text-slate-400 mt-1 italic">"{unit.motto}"</p>
            )}
            <div className="flex items-center gap-6 mt-3">
              <div className="text-center">
                <p className="text-2xl font-extrabold text-primary tabular-nums">{unit.score_total}</p>
                <p className="text-xs text-gray-400 dark:text-slate-500">النقاط</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-extrabold text-accent tabular-nums">{unit.level}</p>
                <p className="text-xs text-gray-400 dark:text-slate-500">المستوى</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-extrabold text-emerald-600 tabular-nums">{unit.members?.length ?? 0}</p>
                <p className="text-xs text-gray-400 dark:text-slate-500">الأعضاء</p>
              </div>
            </div>
          </div>
        </div>
      </Card>

      {/* Leaders */}
      {(unit.leaders?.length ?? 0) > 0 && (
        <Card>
          <h2 className="section-title">القادة</h2>
          <div className="space-y-1.5">
            {unit.leaders!.map((l) => (
              <div key={l.id} className="flex items-center gap-3 p-2.5 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                <Avatar name={l.user?.full_name ?? ''} size="sm" />
                <div className="flex-1 min-w-0">
                  <p className="font-semibold text-sm text-gray-900 dark:text-white truncate">{l.user?.full_name}</p>
                </div>
                <Badge variant="blue">{l.role_in_unit === 'leader' ? 'قائد' : 'مساعد'}</Badge>
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* Members */}
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h2 className="section-title mb-0 flex items-center gap-2">
            <Users size={18} className="text-primary" />
            الأعضاء ({unit.members?.length ?? 0})
          </h2>
          {isLeader && (
            <Button size="sm" onClick={() => { setSelectedIds([]); setShowAddModal(true) }}>
              <Plus size={14} /> إضافة أعضاء
            </Button>
          )}
        </div>

        {!unit.members?.length ? (
          <div className="text-center py-10">
            <Users size={40} className="mx-auto text-gray-200 dark:text-slate-600 mb-3" />
            <p className="text-sm text-gray-400 dark:text-slate-500">لا يوجد أعضاء في هذه الطليعة</p>
          </div>
        ) : (
          <div className="space-y-1.5">
            {unit.members.map((um) => (
              <div key={um.id} className="flex items-center gap-3 p-2.5 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                <Avatar name={um.member?.full_name ?? ''} url={um.member?.photo_url} size="sm" />
                <div className="flex-1 min-w-0">
                  <p className="font-semibold text-sm text-gray-900 dark:text-white truncate">{um.member?.full_name}</p>
                  {um.member?.rank_stage && (
                    <p className="text-xs text-gray-400 dark:text-slate-500">{um.member.rank_stage}</p>
                  )}
                </div>
                {um.is_primary && <Badge variant="gold">رئيسي</Badge>}
                {isLeader && (
                  <button
                    onClick={() => handleRemove(um.member_id, um.member?.full_name ?? '')}
                    className="p-1.5 rounded-lg text-gray-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                    title="إزالة من الطليعة"
                  >
                    <Trash2 size={14} />
                  </button>
                )}
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Add Members Modal */}
      <Modal open={showAddModal} onClose={() => setShowAddModal(false)} title="إضافة أعضاء للطليعة" size="md">
        <div dir="rtl" className="space-y-4">
          {!availableMembers.length ? (
            <p className="text-sm text-gray-400 text-center py-4">لا يوجد أعضاء متاحون للإضافة</p>
          ) : (
            <div className="space-y-1 max-h-80 overflow-y-auto">
              {availableMembers.map((m) => (
                <label
                  key={m.id}
                  className="flex items-center gap-3 p-2.5 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 cursor-pointer transition-colors"
                >
                  <input
                    type="checkbox"
                    checked={selectedIds.includes(m.id)}
                    onChange={() =>
                      setSelectedIds((prev) =>
                        prev.includes(m.id) ? prev.filter((x) => x !== m.id) : [...prev, m.id]
                      )
                    }
                    className="w-4 h-4 rounded accent-primary flex-shrink-0"
                  />
                  <Avatar name={m.full_name} url={m.photo_url} size="sm" />
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm text-gray-900 dark:text-white truncate">{m.full_name}</p>
                  </div>
                  <SectionBadge section={m.section} />
                </label>
              ))}
            </div>
          )}
          <div className="flex gap-3 justify-end pt-2 border-t border-gray-100 dark:border-slate-700">
            <Button variant="ghost" onClick={() => setShowAddModal(false)}>إلغاء</Button>
            <Button onClick={handleAdd} loading={addMembers.isPending} disabled={!selectedIds.length}>
              إضافة ({selectedIds.length})
            </Button>
          </div>
        </div>
      </Modal>

      {/* Edit Unit Modal */}
      <UnitFormModal open={showEditModal} onClose={() => setShowEditModal(false)} unit={unit} />
    </div>
  )
}
