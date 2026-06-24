import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, Search, Filter, Pencil, Trash2 } from 'lucide-react'
import { useMembers, useDeleteMember } from '@/hooks/useMembers'
import { Button, Avatar, Badge, SectionBadge, Spinner, EmptyState } from '@/components/ui'
import { useAuth } from '@/features/auth/AuthContext'
import { canManageMembers } from '@/lib/permissions'
import { toast } from 'sonner'
import type { Member } from '@/types'
import { MemberFormModal } from './MemberFormModal'

const statusLabel: Record<string, string> = { active: 'نشط', inactive: 'غير نشط' }

export function MembersPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const [search, setSearch] = useState('')
  const [section, setSection] = useState('')
  const [status, setStatus] = useState('')
  const [page, setPage] = useState(1)
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Member | undefined>(undefined)

  const canManage = canManageMembers(user)

  const { data, isLoading } = useMembers({ search: search || undefined, section: section || undefined, status: status || undefined, page })
  const deleteMutation = useDeleteMember()

  const handleDelete = async (m: Member) => {
    if (!confirm(`هل أنت متأكد من حذف ${m.full_name}؟`)) return
    await deleteMutation.mutateAsync(m.id)
    toast.success('تم حذف العضو')
  }

  const openCreate = () => { setEditing(undefined); setShowForm(true) }
  const openEdit = (m: Member) => { setEditing(m); setShowForm(true) }

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">الأعضاء</h1>
        {canManage && <Button onClick={openCreate}><Plus size={16} />إضافة عضو</Button>}
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-3">
        <div className="relative flex-1 min-w-48">
          <Search size={16} className="absolute top-2.5 right-3 text-gray-400" />
          <input
            type="text" placeholder="بحث بالاسم أو الهاتف..." value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="input pr-9 w-full"
          />
        </div>
        <select value={section} onChange={(e) => setSection(e.target.value)} className="input w-40">
          <option value="">كل الشعب</option>
          <option value="ashbal">أشبال</option>
          <option value="kashaf">كشاف</option>
          <option value="jawala">جوالة</option>
          <option value="mukashe">مكاشفة</option>
        </select>
        <select value={status} onChange={(e) => setStatus(e.target.value)} className="input w-36">
          <option value="">كل الحالات</option>
          <option value="active">نشط</option>
          <option value="inactive">غير نشط</option>
        </select>
      </div>

      {/* Table */}
      {isLoading ? <Spinner className="h-48" /> : !data?.data?.length ? (
        <EmptyState icon={Filter} title="لا توجد نتائج" description="جرب تغيير معايير البحث" />
      ) : (
        <>
          <div className="card overflow-hidden" style={{ padding: 0 }}>
            <table className="w-full text-sm">
              <thead className="bg-gray-50 dark:bg-slate-900/50">
                <tr>
                  {['العضو', 'الشعبة', 'الهاتف', 'المستوى / XP', 'الحالة', ''].map((h) => (
                    <th key={h} className="px-4 py-3 text-right font-medium text-gray-500 dark:text-slate-400 text-xs">
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-slate-700">
                {data.data.map((m) => (
                  <tr key={m.id} className="hover:bg-gray-50 dark:hover:bg-slate-700/30 cursor-pointer transition-colors"
                    onClick={() => navigate(`/members/${m.id}`)}>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        <Avatar name={m.full_name} url={m.photo_url} size="sm" />
                        <span className="font-medium text-gray-900 dark:text-white">{m.full_name}</span>
                      </div>
                    </td>
                    <td className="px-4 py-3"><SectionBadge section={m.section} /></td>
                    <td className="px-4 py-3 text-gray-500 dark:text-slate-400 tabular-nums ltr">{m.parent_phone}</td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        <span className="font-semibold text-primary">{m.level}</span>
                        <span className="text-gray-400 text-xs">· {m.xp_total} XP</span>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <Badge variant={m.status === 'active' ? 'green' : 'gray'}>{statusLabel[m.status]}</Badge>
                    </td>
                    <td className="px-4 py-3 text-left" onClick={(e) => e.stopPropagation()}>
                      {canManage && (
                        <div className="flex items-center gap-1 justify-end">
                          <button onClick={() => openEdit(m)} title="تعديل"
                            className="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10 transition-colors">
                            <Pencil size={15} />
                          </button>
                          <button onClick={() => handleDelete(m)} title="حذف"
                            className="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                            <Trash2 size={15} />
                          </button>
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {/* Pagination */}
          {data.meta && data.meta.total_pages > 1 && (
            <div className="flex items-center justify-between text-sm text-gray-500">
              <span>{data.meta.total} عضو</span>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(p => p - 1)}>السابق</Button>
                <span className="px-3 py-1">{page} / {data.meta.total_pages}</span>
                <Button variant="outline" size="sm" disabled={page >= data.meta.total_pages} onClick={() => setPage(p => p + 1)}>التالي</Button>
              </div>
            </div>
          )}
        </>
      )}

      <MemberFormModal open={showForm} onClose={() => setShowForm(false)} member={editing} />
    </div>
  )
}
