import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'sonner'
import { queryClient } from '@/lib/queryClient'
import { AuthProvider } from '@/features/auth/AuthContext'
import { ProtectedRoute } from '@/routes/ProtectedRoute'
import { AppLayout } from '@/components/layout/AppLayout'
import { LoginPage } from '@/features/auth/LoginPage'
import { DashboardPage } from '@/features/dashboard/DashboardPage'
import { MembersPage } from '@/features/members/MembersPage'
import { MemberProfilePage } from '@/features/members/MemberProfilePage'
import { UnitsPage } from '@/features/units/UnitsPage'
import { UnitDetailPage } from '@/features/units/UnitDetailPage'
import { ActivitiesPage } from '@/features/activities/ActivitiesPage'
import { ActivityDetailPage } from '@/features/activities/ActivityDetailPage'
import { BadgesPage } from '@/features/badges/BadgesPage'
import { TrainingPage } from '@/features/training/TrainingPage'
import { TrainingLessonDetailPage } from '@/features/training/TrainingLessonDetailPage'
import { CommunicationPage } from '@/features/communication/CommunicationPage'
import { ReportsPage } from '@/features/reports/ReportsPage'
import { EquipmentPage } from '@/features/equipment/EquipmentPage'
import { LeaderboardPage } from '@/features/leaderboard/LeaderboardPage'
import '@/lib/i18n'

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />

            <Route element={<ProtectedRoute />}>
              <Route element={<AppLayout />}>
                <Route index element={<DashboardPage />} />
                <Route path="members" element={<MembersPage />} />
                <Route path="members/:id" element={<MemberProfilePage />} />
                <Route path="units" element={<UnitsPage />} />
                <Route path="units/:id" element={<UnitDetailPage />} />
                <Route path="activities" element={<ActivitiesPage />} />
                <Route path="activities/:id" element={<ActivityDetailPage />} />
                <Route path="badges" element={<BadgesPage />} />
                <Route path="training" element={<TrainingPage />} />
                <Route path="training/:id" element={<TrainingLessonDetailPage />} />
                <Route path="communication" element={<CommunicationPage />} />
                <Route path="reports" element={<ReportsPage />} />
                <Route path="equipment" element={<EquipmentPage />} />
                <Route path="leaderboard" element={<LeaderboardPage />} />
              </Route>
            </Route>

            <Route path="/403" element={
              <div className="flex h-screen items-center justify-center text-gray-500">
                403 — Forbidden
              </div>
            } />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>

          <Toaster position="top-center" richColors />
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
