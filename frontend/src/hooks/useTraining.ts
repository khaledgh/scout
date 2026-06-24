import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { TrainingLesson, Quiz, ApiResponse } from '@/types'

export function useTrainingLessons() {
  return useQuery({
    queryKey: ['training', 'lessons'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<TrainingLesson[]>>('/training/lessons')
      return data.data
    },
  })
}

export function useTrainingLesson(id: number | undefined) {
  return useQuery({
    queryKey: ['training', 'lessons', id],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<TrainingLesson>>(`/training/lessons/${id}`)
      return data.data
    },
    enabled: !!id,
  })
}

export function useLessonQuiz(lessonId: number | undefined) {
  return useQuery({
    queryKey: ['training', 'lessons', lessonId, 'quiz'],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Quiz>>(`/training/lessons/${lessonId}/quiz`)
      return data.data
    },
    enabled: !!lessonId,
  })
}

export function useSubmitQuizAttempt(quizId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (answers: number[]) => api.post(`/training/quizzes/${quizId}/attempt`, { answers }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['training', 'progress'] }),
  })
}

export function useCreateLesson() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post('/training/lessons', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['training', 'lessons'] }),
  })
}

export function useUpdateLesson(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.put(`/training/lessons/${id}`, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['training', 'lessons'] })
      qc.invalidateQueries({ queryKey: ['training', 'lessons', id] })
    },
  })
}

export function useDeleteLesson() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/training/lessons/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['training', 'lessons'] }),
  })
}

export function useCreateQuiz(lessonId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: unknown) => api.post(`/training/lessons/${lessonId}/quiz`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['training', 'lessons', lessonId, 'quiz'] }),
  })
}

export function useUploadLessonMedia(lessonId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => {
      const form = new FormData()
      form.append('file', file)
      return api.post(`/training/lessons/${lessonId}/media`, form)
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['training', 'lessons', lessonId] }),
  })
}

export function useDeleteLessonMedia(lessonId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (mediaId: number) => api.delete(`/training/lessons/${lessonId}/media/${mediaId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['training', 'lessons', lessonId] }),
  })
}

export function useUploadLessonCover(lessonId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => {
      const form = new FormData()
      form.append('file', file)
      return api.post(`/training/lessons/${lessonId}/cover`, form)
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['training', 'lessons', lessonId] }),
  })
}

export function useMyTrainingProgress() {
  return useQuery({
    queryKey: ['training', 'progress'],
    queryFn: async () => {
      const { data } = await api.get('/training/me/progress')
      return data.data as Array<{ lesson_id: number; lesson_title: string; passed: boolean; best_score: number }>
    },
  })
}
