import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Equipment, ApiResponse } from '@/types'

export function useEquipment() {
  return useQuery({
    queryKey: ['equipment'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Equipment[]>>('/equipment')
      return data.data
    },
  })
}

export function useCreateEquipment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post('/equipment', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['equipment'] }),
  })
}

export function useUpdateEquipment(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/equipment/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['equipment'] }),
  })
}

export function useDeleteEquipment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/equipment/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['equipment'] }),
  })
}
