import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'
import type { DashboardKPIs } from '@/types'

export function useDashboard() {
  return useQuery({
    queryKey: ['dashboard'],
    queryFn: async () => {
      const { data } = await api.get<{ data: DashboardKPIs }>('/dashboard')
      return data.data
    },
    refetchInterval: 60_000,
  })
}
