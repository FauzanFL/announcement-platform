import { NavLink, useNavigate } from "react-router";
import useAuthStore from "@/store/authStore";
import useNotificationStore from "@/store/notificationStore";
import { MegaphoneIcon, BellIcon, LogoutIcon, XMarkIcon } from "./icons";
import { useState } from "react";
import api from "@/lib/api";

interface NavItemProps {
  to: string;
  icon: React.ReactNode;
  label: string;
  badge?: number | null;
  onClick?: () => void;
}

function NavItem({ to, icon, label, badge, onClick }: NavItemProps) {
  return (
    <NavLink
      to={to}
      end
      onClick={onClick}
      className={({ isActive }) =>
        `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-150
         ${
           isActive
             ? "bg-brand-500/15 text-brand-400 border border-brand-500/20"
             : "text-slate-400 hover:text-slate-200 hover:bg-slate-800"
         }`
      }
    >
      <span className="shrink-0">{icon}</span>
      <span className="flex-1">{label}</span>
      {badge != null && badge > 0 && (
        <span
          className="flex items-center justify-center min-w-5 h-5 px-1.5 rounded-full
                   bg-amber-400/15 text-amber-400 text-xs font-semibold border border-amber-400/20"
        >
          {badge > 99 ? "99+" : badge}
        </span>
      )}
    </NavLink>
  );
}

function HamburgerIcon({ className = "w-5 h-5" }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
      />
    </svg>
  );
}

export default function Sidebar() {
  const { user, logout } = useAuthStore();
  const unreadCount = useNotificationStore((s) => s.unreadCount);
  const navigate = useNavigate();
  const [drawerOpen, setDrawerOpen] = useState(false);

  const handleLogout = async () => {
    try {
      await api.post("/logout");
    } finally {
      logout();
      void navigate("/login");
      setDrawerOpen(false);
    }
  };

  const isAdmin = user?.role === "admin";

  const navContent = (
    <>
      <nav className="flex-1 px-3 py-4 space-y-1">
        {isAdmin ? (
          <NavItem
            to="/admin"
            icon={<MegaphoneIcon />}
            label="Kelola Pengumuman"
            onClick={() => setDrawerOpen(false)}
          />
        ) : (
          <>
            <NavItem
              to="/"
              icon={<MegaphoneIcon />}
              label="Pengumuman"
              onClick={() => setDrawerOpen(false)}
            />
            <NavItem
              to="/notifications"
              icon={<BellIcon />}
              label="Notifikasi"
              badge={unreadCount}
              onClick={() => setDrawerOpen(false)}
            />
          </>
        )}
      </nav>

      <div className="border-t border-slate-800 p-3">
        <div className="flex items-center gap-3 px-3 py-2 rounded-lg">
          <div
            className="w-8 h-8 rounded-full bg-brand-500/20 border border-brand-500/30
                          flex items-center justify-center text-brand-400 font-semibold text-sm shrink-0"
          >
            {user?.username?.[0]?.toUpperCase()}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-slate-200 truncate">
              {user?.username}
            </p>
            <p className="text-xs text-slate-500 capitalize">{user?.role}</p>
          </div>
          <button
            onClick={() => {
              void handleLogout();
            }}
            className="text-slate-500 hover:text-red-400 transition-colors shrink-0"
            title="Logout"
          >
            <LogoutIcon className="w-4 h-4" />
          </button>
        </div>
      </div>
    </>
  );

  return (
    <>
      <aside className="hidden lg:flex flex-col w-64 min-h-screen bg-slate-900 border-r border-slate-800">
        {/* Logo */}
        <div className="flex items-center gap-3 px-6 py-5 border-b border-slate-800">
          <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-brand-500">
            <MegaphoneIcon className="w-4 h-4 text-white" />
          </div>
          <div>
            <p className="text-sm font-semibold text-slate-100">
              Information Board
            </p>
            <p className="text-xs text-slate-500 capitalize">{user?.role}</p>
          </div>
        </div>

        {navContent}
      </aside>

      {/* ─── Tablet/mobile topbar ───────────────────────────────────────────── */}
      <div
        className="lg:hidden fixed top-0 left-0 right-0 z-40 flex items-center justify-between
                      px-4 h-14 bg-slate-900 border-b border-slate-800"
      >
        <div className="flex items-center gap-2.5">
          <div className="flex items-center justify-center w-7 h-7 rounded-lg bg-brand-500">
            <MegaphoneIcon className="w-3.5 h-3.5 text-white" />
          </div>
          <p className="text-sm font-semibold text-slate-100">Papan Info</p>
        </div>

        <div className="flex items-center gap-3">
          {/* Unread badge di topbar (mobile) */}
          {!isAdmin && unreadCount > 0 && (
            <NavLink
              to="/notifications"
              className="relative text-slate-400 hover:text-slate-200 transition-colors"
            >
              <BellIcon className="w-5 h-5" />
              <span
                className="absolute -top-1 -right-1 flex items-center justify-center
                               min-w-4 h-4 px-1 rounded-full bg-amber-400 text-slate-900
                               text-[10px] font-bold"
              >
                {unreadCount > 99 ? "99+" : unreadCount}
              </span>
            </NavLink>
          )}
          <button
            onClick={() => setDrawerOpen(true)}
            className="text-slate-400 hover:text-slate-200 transition-colors"
          >
            <HamburgerIcon />
          </button>
        </div>
      </div>

      {/* ─── Drawer overlay ─────────────────────────────────────────────────── */}
      {drawerOpen && (
        <div
          className="lg:hidden fixed inset-0 z-50 flex"
          onClick={() => setDrawerOpen(false)}
        >
          {/* Backdrop */}
          <div className="absolute inset-0 bg-black/60 backdrop-blur-sm animate-fade-in" />

          {/* Drawer panel */}
          <div
            className="relative ml-auto w-72 h-full bg-slate-900 border-l border-slate-800
                        flex flex-col animate-slide-in shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between px-6 py-4 border-b border-slate-800">
              <div className="flex items-center gap-2.5">
                <div className="flex items-center justify-center w-7 h-7 rounded-lg bg-brand-500">
                  <MegaphoneIcon className="w-3.5 h-3.5 text-white" />
                </div>
                <p className="text-sm font-semibold text-slate-100">
                  Papan Info
                </p>
              </div>
              <button
                onClick={() => setDrawerOpen(false)}
                className="text-slate-500 hover:text-slate-300 transition-colors"
              >
                <XMarkIcon className="w-5 h-5" />
              </button>
            </div>
            {navContent}
          </div>
        </div>
      )}
    </>
  );
}
