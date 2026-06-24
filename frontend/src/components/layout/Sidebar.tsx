import { useEffect, useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import {
  LayoutDashboard, Users, Shield, CalendarDays, Award,
  BookOpen, Trophy, MessageSquare, BarChart3, Package, Compass, X,
} from 'lucide-react'
import { useAuth } from '@/features/auth/AuthContext'
import { Avatar } from '@/components/ui'

const navItems = [
  { key: 'dashboard',      path: '/',               icon: LayoutDashboard },
  { key: 'members',        path: '/members',         icon: Users },
  { key: 'units',          path: '/units',           icon: Shield },
  { key: 'activities',     path: '/activities',      icon: CalendarDays },
  { key: 'badges',         path: '/badges',          icon: Award },
  { key: 'training',       path: '/training',        icon: BookOpen },
  { key: 'leaderboard',    path: '/leaderboard',     icon: Trophy },
  { key: 'communication',  path: '/communication',   icon: MessageSquare },
  { key: 'reports',        path: '/reports',         icon: BarChart3 },
  { key: 'equipment',      path: '/equipment',       icon: Package },
]

interface Props {
  open: boolean
  onClose: () => void
}

function SidebarContent({
  t, user, collapsed, onClose, showClose,
}: {
  t: (k: string) => string
  user: ReturnType<typeof useAuth>['user']
  collapsed: boolean
  onClose?: () => void
  showClose?: boolean
}) {
  return (
    <div dir="rtl" className="flex flex-col h-full">
      {/* Logo */}
      <div className={`flex items-center h-16 border-b border-gray-100 dark:border-slate-800 flex-shrink-0 ${collapsed ? 'justify-center px-2' : 'gap-3 px-5'}`}>
        <div className="flex-shrink-0 w-9 h-9 rounded-xl bg-gradient-to-br from-primary to-accent flex items-center justify-center shadow-sm">
          <Compass size={20} className="text-white" />
        </div>
        {!collapsed && (
          <div className="flex-1 min-w-0">
            <p className="font-extrabold text-gray-900 dark:text-white text-base leading-none">كشفي</p>
            <p className="text-[10px] text-gray-400 dark:text-slate-500 leading-none mt-0.5">Scout Management</p>
          </div>
        )}
        {showClose && onClose && (
          <button onClick={onClose} className="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 dark:hover:bg-slate-800 transition-colors flex-shrink-0">
            <X size={18} />
          </button>
        )}
      </div>

      {/* Nav */}
      <nav className="flex-1 overflow-y-auto py-3 px-2 space-y-0.5">
        {navItems.map(({ key, path, icon: Icon }) => (
          <NavLink
            key={key}
            to={path}
            end={path === '/'}
            title={collapsed ? t(`nav.${key}`) : undefined}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-150 ${
                isActive
                  ? 'bg-primary text-white shadow-sm'
                  : 'text-gray-500 dark:text-slate-400 hover:bg-gray-50 dark:hover:bg-slate-800 hover:text-gray-900 dark:hover:text-slate-100'
              }`
            }
          >
            <Icon size={18} className="flex-shrink-0" />
            {!collapsed && <span>{t(`nav.${key}`)}</span>}
          </NavLink>
        ))}
      </nav>

      {/* User section */}
      <div className={`border-t border-gray-100 dark:border-slate-800 p-3 flex-shrink-0 ${collapsed ? 'flex justify-center' : ''}`}>
        {user && (
          collapsed ? (
            <Avatar name={user.full_name} url={user.avatar_url} size="sm" />
          ) : (
            <div className="flex items-center gap-3 px-2 py-2 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
              <Avatar name={user.full_name} url={user.avatar_url} size="sm" />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-semibold text-gray-800 dark:text-white truncate">{user.full_name}</p>
                <p className="text-[11px] text-gray-400 dark:text-slate-500 truncate">
                  {user.role === 'super_admin' ? 'مسؤول عام' : user.role === 'leader' ? 'قائد' : 'مساعد'}
                </p>
              </div>
            </div>
          )
        )}
      </div>
    </div>
  )
}

export function Sidebar({ open, onClose }: Props) {
  const { t }    = useTranslation()
  const { user } = useAuth()
  const location = useLocation()

  // Track whether we're on desktop using matchMedia (more reliable than CSS classes)
  const [isDesktop, setIsDesktop] = useState(() => window.innerWidth >= 768)

  useEffect(() => {
    const mq = window.matchMedia('(min-width: 768px)')
    const handler = (e: MediaQueryListEvent) => {
      setIsDesktop(e.matches)
      if (!e.matches) onClose() // always close when switching to mobile
    }
    mq.addEventListener('change', handler)
    return () => mq.removeEventListener('change', handler)
  }, [onClose])

  // Close overlay on navigation (mobile only)
  useEffect(() => {
    if (!isDesktop) onClose()
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname])

  // Prevent body scroll when mobile drawer is open
  useEffect(() => {
    if (!isDesktop && open) {
      document.body.classList.add('overflow-hidden')
    } else {
      document.body.classList.remove('overflow-hidden')
    }
    return () => document.body.classList.remove('overflow-hidden')
  }, [isDesktop, open])

  const commonAside = 'flex flex-col bg-white dark:bg-slate-900 border-e border-gray-100 dark:border-slate-800 shadow-sm flex-shrink-0'

  return (
    <>
      {/* ── Mobile: fixed overlay from the right (RTL-correct) ────── */}
      {!isDesktop && (
        <>
          {/* Backdrop */}
          <div
            className={`fixed inset-0 bg-black/50 z-40 transition-opacity duration-300 ${
              open ? 'opacity-100' : 'opacity-0 pointer-events-none'
            }`}
            onClick={onClose}
          />
          {/* Drawer — slides in from the right */}
          <aside
            className={`${commonAside} fixed inset-y-0 right-0 z-50 w-72 transition-transform duration-300 ${
              open ? 'translate-x-0 shadow-2xl' : 'translate-x-full'
            }`}
          >
            <SidebarContent t={t} user={user} collapsed={false} onClose={onClose} showClose />
          </aside>
        </>
      )}

      {/* ── Desktop: inline sidebar ────────────────────────────────── */}
      {isDesktop && (
        <aside
          className={`${commonAside} h-screen transition-all duration-300 ${open ? 'w-64' : 'w-[68px]'}`}
        >
          <SidebarContent t={t} user={user} collapsed={!open} />
        </aside>
      )}
    </>
  )
}
