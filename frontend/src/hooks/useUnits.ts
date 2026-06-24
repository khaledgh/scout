import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Unit, ApiResponse } from '@/types'

export function useUnits() {
  return useQuery({
    queryKey: ['units'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Unit[]>>('/units')
      return data.data
    },
  })
}

export function useUnit(id: number | undefined) {
  return useQuery({
    queryKey: ['units', id],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Unit>>(`/units/${id}`)
      return data.data
    },
    enabled: !!id,
  })
}

export function useUnitLeaderboard() {
  return useQuery({
    queryKey: ['units', 'leaderboard'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Unit[]>>('/units/leaderboard')
      return data.data
    },
  })
}

export function useCreateUnit() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post('/units', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['units'] }),
  })
}

export function useUpdateUnit(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/units/${id}`, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['units'] })
      qc.invalidateQueries({ queryKey: ['units', id] })
    },
  })
}

export function useDeleteUnit() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/units/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['units'] }),
  })
}

export function useAddUnitMembers(unitId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (memberIds: number[]) => api.post(`/units/${unitId}/members`, { member_ids: memberIds }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['units', unitId] }),
  })
}

export function useAssignMemberToUnit() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ unitId, memberId }: { unitId: number; memberId: number }) =>
      api.post(`/units/${unitId}/members`, { member_ids: [memberId] }),
    onSuccess: (_, { memberId, unitId }) => {
      qc.invalidateQueries({ queryKey: ['members', memberId] })
      qc.invalidateQueries({ queryKey: ['units', unitId] })
      qc.invalidateQueries({ queryKey: ['units'] })
    },
  })
}

export function useRemoveUnitMember(unitId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (memberId: number) => api.delete(`/units/${unitId}/members/${memberId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['units', unitId] }),
  })
}
