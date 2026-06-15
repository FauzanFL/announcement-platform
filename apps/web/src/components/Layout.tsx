import type { ReactNode } from 'react'
import Sidebar from './Sidebar'
import Toast from './Toast'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="flex min-h-screen bg-slate-950">
      <Sidebar />
      <main className="flex-1 overflow-y-auto">
        {children}
      </main>
      <Toast />
    </div>
  )
}