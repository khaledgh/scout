import { QueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60_000,
      retry: (failureCount, error: unknown) => {
        const status = (error as { response?: { status?: number } })?.response?.status
        if (status === 401 || status === 403 || status === 404) return false
        return failureCount < 2
      },
    },
    mutations: {
      onError: (error: unknown) => {
        const msg =
          (error as { response?: { data?: { error?: { message?: string } } } })?.response?.data
            ?.error?.message ?? 'An unexpected error occurred'
        toast.error(msg)
      },
    },
  },
})
