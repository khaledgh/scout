import { useState, useRef, useEffect } from 'react'
import { Outlet, useLocation } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { Menu, Bell, ChevronDown, LogOut, User } from 'lucide-react'
import { useAuth } from '@/features/auth/AuthContext'
import { useNotifications, useMarkNotificationRead } from '@/hooks/useNotifications'
import { Avatar } from '@/components/ui'
import { format } from 'date-fns'
import { arLB as ar } from '@/lib/arLB'

const pageTitles: Record<string, string> = {
  '/':             'لوحة التحكم',
  '/members':      'الأعضاء',
  '/units':        'الطلائع',
  '/activities':   'الأنشطة',
  '/badges':       'الشارات',
  '/training':     'التدريب',
  '/leaderboard':  'المتصدرون',
  '/communication':'التواصل',
  '/reports':      'التقارير',
  '/equipment':    'المعدات',
}

export function AppLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [showNotifs, setShowNotifs] = useState(false)
  const [showUser, setShowUser]     = useState(false)
  const { user, logout }            = useAuth()
  const { data: notifications }     = useNotifications()
  const markRead                    = useMarkNotificationRead()
  const location                    = useLocation()
  const notifsRef                   = useRef<HTMLDivElement>(null)
  const userRef                     = useRef<HTMLDivElement>(null)

  const unread   = notifications?.filter((n) => !n.read_at) ?? []
  const pageTitle = pageTitles[location.pathname] ?? pageTitles[Object.keys(pageTitles).find((k) => k !== '/' && location.pathname.startsWith(k)) ?? ''] ?? ''

  // close dropdowns on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (notifsRef.current && !notifsRef.current.contains(e.target as Node)) setShowNotifs(false)
      if (userRef.current   && !userRef.current.contains(e.target as Node))   setShowUser(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  return (
    <div className="flex h-screen overflow-hidden bg-gray-50 dark:bg-slate-950">
      <Sidebar open={sidebarOpen} />

      <div className="flex-1 flex flex-col overflow-hidden min-w-0">
        {/* ── Topbar ────────────────────────────────────────────────────── */}
        <header className="flex items-center gap-4 px-6 h-16 flex-shrink-0 bg-white dark:bg-slate-900 border-b border-gray-100 dark:border-slate-800 shadow-sm z-30">
          {/* Hamburger */}
          <button
            onClick={() => setSidebarOpen((v) => !v)}
            className="p-2 rounded-xl text-gray-500 hover:bg-gray-100 dark:text-slate-400 dark:hover:bg-slate-800 transition-colors flex-shrink-0"
          >
            <Menu size={20} />
          </button>

          {/* Page title */}
          <h2 className="font-bold text-gray-800 dark:text-slate-100 text-lg hidden sm:block" dir="rtl">
            {pageTitle}
          </h2>

          <div className="flex-1" />

          {/* Notification bell */}
          <div ref={notifsRef} className="relative">
            <button
              onClick={() => { setShowNotifs((v) => !v); setShowUser(false) }}
              className="relative p-2 rounded-xl text-gray-500 hover:bg-gray-100 dark:text-slate-400 dark:hover:bg-slate-800 transition-colors"
            >
              <Bell size={20} />
              {unread.length > 0 && (
                <span className="absolute top-1 right-1 w-4 h-4 rounded-full bg-red-500 text-white text-[9px] font-bold flex items-center justify-center">
                  {unread.length > 9 ? '9+' : unread.length}
                </span>
              )}
            </button>

            {showNotifs && (
              <div className="absolute left-0 top-12 w-80 bg-white dark:bg-slate-900 rounded-2xl shadow-xl border border-gray-100 dark:border-slate-700 z-50 overflow-hidden" dir="rtl">
                <div className="px-4 py-3 border-b border-gray-100 dark:border-slate-800 flex items-center justify-between">
                  <span className="font-semibold text-gray-900 dark:text-white text-sm">الإشعارات</span>
                  {unread.length > 0 && <span className="badge bg-primary/10 text-primary">{unread.length} جديد</span>}
                </div>
                <div className="max-h-72 overflow-y-auto divide-y divide-gray-50 dark:divide-slate-800">
                  {!notifications?.length ? (
                    <p className="text-center text-sm text-gray-400 py-8">لا توجد إشعارات</p>
                  ) : notifications.slice(0, 10).map((n) => (
                    <button
                      key={n.id}
                      onClick={() => markRead.mutate(n.id)}
                      className={`w-full text-right px-4 py-3 hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors ${!n.read_at ? 'bg-primary/5' : ''}`}
                    >
                      <p className="text-sm font-medium text-gray-900 dark:text-white">{n.title}</p>
                      <p className="text-xs text-gray-500 dark:text-slate-400 mt-0.5 truncate">{n.body}</p>
                      <p className="text-[10px] text-gray-400 mt-1">{format(new Date(n.created_at), 'PPp', { locale: ar })}</p>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* User avatar + dropdown */}
          <div ref={userRef} className="relative">
            <button
              onClick={() => { setShowUser((v) => !v); setShowNotifs(false) }}
              className="flex items-center gap-2 p-1.5 rounded-xl hover:bg-gray-100 dark:hover:bg-slate-800 transition-colors"
            >
              <Avatar name={user?.full_name ?? ''} url={user?.avatar_url} size="sm" />
              <div className="text-left hidden md:block">
                <p className="text-xs font-semibold text-gray-800 dark:text-slate-100 leading-none">{user?.full_name}</p>
              </div>
              <ChevronDown size={14} className="text-gray-400" />
            </button>

            {showUser && (
              <div className="absolute left-0 top-12 w-52 bg-white dark:bg-slate-900 rounded-2xl shadow-xl border border-gray-100 dark:border-slate-700 z-50 overflow-hidden" dir="rtl">
                <div className="px-4 py-3 border-b border-gray-100 dark:border-slate-800">
                  <p className="font-semibold text-gray-900 dark:text-white text-sm">{user?.full_name}</p>
                  <p className="text-xs text-gray-400 mt-0.5">{user?.phone}</p>
                </div>
                <div className="p-1.5">
                  <button className="flex items-center gap-3 w-full px-3 py-2 rounded-xl text-sm text-gray-700 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
                    <User size={16} /> الملف الشخصي
                  </button>
                  <button
                    onClick={logout}
                    className="flex items-center gap-3 w-full px-3 py-2 rounded-xl text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                  >
                    <LogOut size={16} /> تسجيل الخروج
                  </button>
                </div>
              </div>
            )}
          </div>
        </header>

        {/* ── Page content ──────────────────────────────────────────────── */}
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
