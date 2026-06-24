import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Announcement, Notification, Channel, Message, ApiResponse } from '@/types'

export function useAnnouncements() {
  return useQuery({
    queryKey: ['announcements'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Announcement[]>>('/announcements')
      return data.data
    },
  })
}

export function useCreateAnnouncement() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post('/announcements', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['announcements'] }),
  })
}

export function useUpdateAnnouncement(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/announcements/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['announcements'] }),
  })
}

export function useDeleteAnnouncement() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/announcements/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['announcements'] }),
  })
}

export function useNotifications() {
  return useQuery({
    queryKey: ['notifications'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Notification[]>>('/notifications')
      return data.data
    },
    refetchInterval: 30_000,
  })
}

export function useMarkNotificationRead() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.put(`/notifications/${id}/read`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['notifications'] }),
  })
}

export function useChannels() {
  return useQuery({
    queryKey: ['channels'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Channel[]>>('/channels')
      return data.data
    },
  })
}

export function useChannelMessages(channelId: number | undefined) {
  return useQuery({
    queryKey: ['channels', channelId, 'messages'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Message[]>>(`/channels/${channelId}/messages`)
      return data.data
    },
    enabled: !!channelId,
  })
}
