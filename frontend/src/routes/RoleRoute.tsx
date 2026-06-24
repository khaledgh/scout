import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '@/features/auth/AuthContext'
import type { Role } from '@/types'

interface Props {
  roles: Role[]
}

export function RoleRoute({ roles }: Props) {
  const { user } = useAuth()

  if (!user || !roles.includes(user.role)) {
    return <Navigate to="/403" replace />
  }

  return <Outlet />
}
