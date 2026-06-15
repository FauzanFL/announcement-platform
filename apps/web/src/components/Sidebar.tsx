import { NavLink, useNavigate } from 'react-router'
import useAuthStore from '@/store/authStore'
import useNotificationStore from '@/store/notificationStore'
import { MegaphoneIcon, BellIcon, LogoutIcon } from './icons'

interface NavItemProps {
  to: string
  icon: React.ReactNode
  label: string
  badge?: number | null
}

function NavItem({ to, icon, label, badge }: NavItemProps) {
  return (
    <NavLink
      to={to}
      end
      className={({ isActive }) =>
        `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-150
         ${isActive
          ? 'bg-brand-500/15 text-brand-400 border border-brand-500/20'
          : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
        }`
      }
    >
      <span className="shrink-0">{icon}</span>
      <span className="flex-1">{label}</span>
      {badge != null && badge > 0 && (
        <span className="flex items-center justify-center min-w-[20px] h-5 px-1.5 rounded-full
                         bg-amber-400/15 text-amber-400 text-xs font-semibold border border-amber-400/20">
          {badge > 99 ? '99+' : badge}
        </span>
      )}
    </NavLink>
  )
}

export default function Sidebar() {
  const { user, logout } = useAuthStore()
  const unreadCount = useNotificationStore((s) => s.unreadCount)
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const isAdmin = user?.role === 'admin'

  return (
    <aside className="flex flex-col w-64 min-h-screen bg-slate-900 border-r border-slate-800">
      {/* Logo */}
      <div className="flex items-center gap-3 px-6 py-5 border-b border-slate-800">
        <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-brand-500">
          <MegaphoneIcon className="w-4 h-4 text-white" />
        </div>
        <div>
          <p className="text-sm font-semibold text-slate-100">Papan Info</p>
          <p className="text-xs text-slate-500 capitalize">{user?.role}</p>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-4 space-y-1">
        {isAdmin ? (
          <NavItem to="/admin" icon={<MegaphoneIcon />} label="Kelola Pengumuman" />
        ) : (
          <>
            <NavItem to="/" icon={<MegaphoneIcon />} label="Pengumuman" />
            <NavItem
              to="/notifications"
              icon={<BellIcon />}
              label="Notifikasi"
              badge={unreadCount}
            />
          </>
        )}
      </nav>

      {/* User footer */}
      <div className="border-t border-slate-800 p-3">
        <div className="flex items-center gap-3 px-3 py-2 rounded-lg">
          <div className="w-8 h-8 rounded-full bg-brand-500/20 border border-brand-500/30
                          flex items-center justify-center text-brand-400 font-semibold text-sm">
            {user?.username?.[0]?.toUpperCase()}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-slate-200 truncate">{user?.username}</p>
            <p className="text-xs text-slate-500 capitalize">{user?.role}</p>
          </div>
          <button
            onClick={handleLogout}
            className="text-slate-500 hover:text-red-400 transition-colors"
            title="Logout"
          >
            <LogoutIcon className="w-4 h-4" />
          </button>
        </div>
      </div>
    </aside>
  )
}