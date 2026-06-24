import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Activity, ActivityAttendance, ApiResponse } from '@/types'

interface ActivityFilters {
  type?: string
  status?: string
  from?: string
  to?: string
  unit_id?: number
  page?: number
}

export function useActivities(filters: ActivityFilters = {}) {
  return useQuery({
    queryKey: ['activities', filters],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Activity[]>>('/activities', { params: filters })
      return data
    },
  })
}

export function useActivity(id: number | undefined) {
  return useQuery({
    queryKey: ['activities', id],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Activity>>(`/activities/${id}`)
      return data.data
    },
    enabled: !!id,
  })
}

export function useCreateActivity() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post('/activities', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['activities'] }),
  })
}

export function useUpdateActivity(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/activities/${id}`, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['activities'] })
      qc.invalidateQueries({ queryKey: ['activities', id] })
    },
  })
}

export function useDeleteActivity() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/activities/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['activities'] }),
  })
}

export function useActivityAttendance(activityId: number | undefined) {
  return useQuery({
    queryKey: ['activities', activityId, 'attendance'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<ActivityAttendance[]>>(`/activities/${activityId}/attendance`)
      return data.data
    },
    enabled: !!activityId,
    refetchInterval: 5000,
  })
}

export function useRecordAttendance(activityId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (records: Array<{ member_id: number; status: string }>) =>
      api.post(`/activities/${activityId}/attendance`, { records }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['activities', activityId, 'attendance'] }),
  })
}

export function useCheckIn(activityId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: { method: string; qr_token?: string; lat?: number; lng?: number }) =>
      api.post(`/activities/${activityId}/checkin`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['activities', activityId, 'attendance'] }),
  })
}

export function useSubmitFeedback(activityId: number) {
  return useMutation({
    mutationFn: (body: unknown) => api.post(`/activities/${activityId}/feedback`, body),
  })
}

export function useUploadActivityMedia(activityId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => {
      const form = new FormData()
      form.append('file', file)
      return api.post(`/activities/${activityId}/media`, form)
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['activities', activityId] }),
  })
}

export function useFeedbackSummary(activityId: number | undefined) {
  return useQuery({
    queryKey: ['activities', activityId, 'feedback'],
    queryFn: async () => {
      const { data } = await api.get(`/activities/${activityId}/feedback/summary`)
      return data.data
    },
    enabled: !!activityId,
  })
}
