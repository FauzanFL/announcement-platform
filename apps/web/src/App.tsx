import { BrowserRouter, Routes, Route, Navigate } from 'react-router'
import Login             from '@/pages/Login'
import Register          from '@/pages/Register'
import AnnouncementsPage from '@/pages/AnnouncementsPage'
import NotificationsPage from '@/pages/NotificationsPage'
import AdminDashboard    from '@/pages/AdminDashboard'
import useAuthStore      from '@/store/authStore'
import type { Role }     from '@/types'
import type { ReactNode } from 'react'

function RequireAuth({ children }: { children: ReactNode }) {
  const token = useAuthStore((s) => s.token)
  return token ? <>{children}</> : <Navigate to="/login" replace />
}

function RequireRole({ children, role }: { children: ReactNode; role: Role }) {
  const user = useAuthStore((s) => s.user)
  if (!user) return <Navigate to="/login" replace />
  if (user.role !== role) return <Navigate to={user.role === 'admin' ? '/admin' : '/'} replace />
  return <>{children}</>
}

function RedirectIfAuthed({ children }: { children: ReactNode }) {
  const { token, user } = useAuthStore()
  if (token) return <Navigate to={user?.role === 'admin' ? '/admin' : '/'} replace />
  return <>{children}</>
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login"    element={<RedirectIfAuthed><Login /></RedirectIfAuthed>} />
        <Route path="/register" element={<RedirectIfAuthed><Register /></RedirectIfAuthed>} />

        <Route path="/" element={
          <RequireRole role="user"><AnnouncementsPage /></RequireRole>
        } />
        <Route path="/notifications" element={
          <RequireRole role="user"><NotificationsPage /></RequireRole>
        } />
        <Route path="/admin" element={
          <RequireRole role="admin"><AdminDashboard /></RequireRole>
        } />

        <Route path="*" element={
          <RequireAuth><Navigate to="/" replace /></RequireAuth>
        } />
      </Routes>
    </BrowserRouter>
  )
}