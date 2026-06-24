import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Member, MemberMedical, ApiResponse, PaginationMeta } from '@/types'

interface MemberFilters {
  page?: number
  page_size?: number
  unit_id?: number
  section?: string
  status?: string
  search?: string
}

export function useMembers(filters: MemberFilters = {}) {
  return useQuery({
    queryKey: ['members', filters],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Member[]> & { meta: PaginationMeta }>('/members', { params: filters })
      return data
    },
  })
}

export function useMember(id: number | undefined) {
  return useQuery({
    queryKey: ['members', id],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Member>>(`/members/${id}`)
      return data.data
    },
    enabled: !!id,
  })
}

export function useCreateMember() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post('/members', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['members'] }),
  })
}

export function useUpdateMember(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/members/${id}`, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['members'] })
      qc.invalidateQueries({ queryKey: ['members', id] })
    },
  })
}

export function useDeleteMember() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/members/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['members'] }),
  })
}

export function useMemberMedical(id: number | undefined) {
  return useQuery({
    queryKey: ['members', id, 'medical'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<MemberMedical>>(`/members/${id}/medical`)
      return data.data
    },
    enabled: !!id,
  })
}

export function useUpsertMedical(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/members/${id}/medical`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['members', id, 'medical'] }),
  })
}

export function useMemberTimeline(id: number | undefined) {
  return useQuery({
    queryKey: ['members', id, 'timeline'],
    queryFn: async () => {
      const { data } = await api.get(`/members/${id}/timeline`)
      return data.data
    },
    enabled: !!id,
  })
}

export function useCreateEvaluation(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post(`/members/${id}/evaluate`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['members', id, 'timeline'] }),
  })
}

export function useUploadMemberPhoto(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => {
      const form = new FormData()
      form.append('file', file)
      return api.post(`/members/${id}/photo`, form, { headers: { 'Content-Type': 'multipart/form-data' } })
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['members'] })
      qc.invalidateQueries({ queryKey: ['members', id] })
    },
  })
}

export function useMemberQR(id: number | undefined) {
  return useQuery({
    queryKey: ['members', id, 'qr'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<{ token: string; member_id: number }>>(`/members/${id}/qr`)
      return data.data
    },
    enabled: !!id,
    staleTime: 20 * 60 * 1000,
  })
}
