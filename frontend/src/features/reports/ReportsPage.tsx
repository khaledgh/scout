import { useState } from 'react'
import { BarChart3, Download } from 'lucide-react'
import { Card, Button, Spinner } from '@/components/ui'
import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts'
import { format } from 'date-fns'

export function ReportsPage() {
  const [month, setMonth] = useState(format(new Date(), 'yyyy-MM'))

  const { data: monthly, isLoading } = useQuery({
    queryKey: ['reports', 'monthly', month],
    queryFn: async () => {
      const { data } = await api.get('/reports/monthly', { params: { month } })
      return data.data as {
        month: string
        total_activities: number
        attendance_rate: number
        new_members: number
        xp_distributed: number
      }
    },
  })

  const handleExport = async () => {
    const response = await api.get('/reports/export', {
      params: { type: 'monthly', month },
      responseType: 'blob',
    })
    const url = URL.createObjectURL(response.data)
    const a = document.createElement('a')
    a.href = url
    a.download = `report-${month}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const chartData = monthly ? [
    { name: 'الأنشطة', value: monthly.total_activities },
    { name: 'نسبة الحضور', value: Math.round(monthly.attendance_rate) },
    { name: 'أعضاء جدد', value: monthly.new_members },
  ] : []

  return (
    <div dir="rtl" className="space-y-5">
      <div className="flex items-center justify-between">
        <h1 className="page-header">التقارير</h1>
        <Button variant="outline" onClick={handleExport}>
          <Download size={16} />تصدير CSV
        </Button>
      </div>

      <div className="flex items-center gap-3">
        <label className="label mb-0">الشهر:</label>
        <input type="month" value={month} onChange={(e) => setMonth(e.target.value)} className="input w-48" />
      </div>

      {isLoading ? <Spinner className="h-48" /> : monthly && (
        <>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            {[
              { label: 'إجمالي الأنشطة', value: monthly.total_activities },
              { label: 'نسبة الحضور', value: `${monthly.attendance_rate.toFixed(1)}%` },
              { label: 'أعضاء جدد', value: monthly.new_members },
              { label: 'XP موزعة', value: monthly.xp_distributed },
            ].map(({ label, value }) => (
              <Card key={label}>
                <p className="text-xs text-gray-500 dark:text-slate-400">{label}</p>
                <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1 tabular-nums">{value}</p>
              </Card>
            ))}
          </div>

          <Card>
            <h3 className="section-title flex items-center gap-2"><BarChart3 size={18} />إحصاءات الشهر</h3>
            <ResponsiveContainer width="100%" height={200}>
              <BarChart data={chartData} margin={{ top: 5, right: 5, left: 0, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                <XAxis dataKey="name" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip />
                <Bar dataKey="value" fill="#6D28D9" radius={[6, 6, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </>
      )}
    </div>
  )
}
