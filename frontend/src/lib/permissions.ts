import type { User } from '@/types'

/** Roles that may manage members/content at the UI level.
 *  The backend additionally enforces unit-scoping for leaders/assistants. */
export function isLeaderRole(user?: User | null): boolean {
  return !!user && (user.role === 'super_admin' || user.role === 'leader' || user.role === 'assistant')
}

export function isAdmin(user?: User | null): boolean {
  return user?.role === 'super_admin'
}

/** Whether the current user can create/edit/delete members from the UI.
 *  Final authority for unit-scoped leaders lives on the server (403 on mismatch). */
export function canManageMembers(user?: User | null): boolean {
  return isLeaderRole(user)
}
