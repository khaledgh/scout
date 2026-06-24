import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Badge, MemberBadge, Skill, ApiResponse } from '@/types'

export function useBadges() {
  return useQuery({
    queryKey: ['badges'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Badge[]>>('/badges')
      return data.data
    },
  })
}

export function useMemberBadges(memberId: number | undefined) {
  return useQuery({
    queryKey: ['members', memberId, 'badges'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<MemberBadge[]>>(`/members/${memberId}/badges`)
      return data.data
    },
    enabled: !!memberId,
  })
}

export function useAwardBadge(memberId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (badgeId: number) => api.post(`/members/${memberId}/badges`, { badge_id: badgeId }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['members', memberId, 'badges'] })
      qc.invalidateQueries({ queryKey: ['members', memberId] })
    },
  })
}

export function useRevokeBadge(memberId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (badgeId: number) => api.delete(`/members/${memberId}/badges/${badgeId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['members', memberId, 'badges'] }),
  })
}

export function useCreateBadge() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post('/badges', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['badges'] }),
  })
}

export function useUpdateBadge(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/badges/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['badges'] }),
  })
}

export function useDeleteBadge() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/badges/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['badges'] }),
  })
}

export function useSkills() {
  return useQuery({
    queryKey: ['skills'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Skill[]>>('/skills')
      return data.data
    },
  })
}

export function useMemberLeaderboard(section?: string) {
  return useQuery({
    queryKey: ['leaderboard', 'members', section],
    queryFn: async () => {
      const { data } = await api.get('/leaderboard/members', { params: { section } })
      return data.data
    },
  })
}

export function useUnitLeaderboardGamification() {
  return useQuery({
    queryKey: ['leaderboard', 'units'],
    queryFn: async () => {
      const { data } = await api.get('/leaderboard/units')
      return data.data
    },
  })
}

export function useXPHistory(memberId: number | undefined) {
  return useQuery({
    queryKey: ['members', memberId, 'xp'],
    queryFn: async () => {
      const { data } = await api.get(`/members/${memberId}/xp`)
      return data.data
    },
    enabled: !!memberId,
  })
}
