import { useEffect, useState, useCallback } from 'react'
import Layout from '@/components/Layout'
import { BellIcon, CheckIcon, SpinnerIcon } from '@/components/icons'
import api from '@/lib/api'
import useNotificationStore from '@/store/notificationStore'
import useToastStore from '@/store/toastStore'
import useSSE from '@/hooks/useSSE'
import type { Notification, UnreadCountResponse } from '@/types'

const formatDate = (iso: string) =>
  new Date(iso).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })

function NotifCard({ notif, onMarkRead }: { notif: Notification; onMarkRead: (id: string) => Promise<void> }) {
  const [loading, setLoading] = useState(false)

  const handle = async () => {
    if (notif.is_read) return
    setLoading(true)
    await onMarkRead(notif.id)
    setLoading(false)
  }

  return (
    <div className={`card overflow-hidden animate-slide-in transition-all duration-200
                     ${!notif.is_read ? 'border-amber-400/20 bg-amber-400/5' : ''}`}>
      <div className="px-5 py-4 flex gap-4">
        <div className="mt-1.5 shrink-0">
          {notif.is_read
            ? <div className="w-2 h-2 rounded-full bg-slate-700" />
            : <div className="w-2 h-2 rounded-full bg-amber-400 ring-4 ring-amber-400/20 animate-pulse" />
          }
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <p className={`text-sm font-semibold ${notif.is_read ? 'text-slate-400' : 'text-slate-100'}`}>
              {notif.announcement.title}
            </p>
            {!notif.is_read && (
              <span className="shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-xs
                               font-medium bg-amber-400/15 text-amber-400 border border-amber-400/20">
                New
              </span>
            )}
          </div>
          <p className={`text-sm mt-1 leading-relaxed line-clamp-2
                         ${notif.is_read ? 'text-slate-600' : 'text-slate-400'}`}>
            {notif.announcement.content}
          </p>
          <div className="flex items-center justify-between mt-3">
            <span className="text-xs text-slate-600 font-mono">{formatDate(notif.created_at)}</span>
            {!notif.is_read && (
              <button
                onClick={() => { void handle() }}
                disabled={loading}
                className="inline-flex items-center gap-1.5 text-xs text-brand-400 hover:text-brand-300
                           font-medium transition-colors disabled:opacity-50"
              >
                {loading ? <SpinnerIcon className="w-3 h-3" /> : <CheckIcon className="w-3 h-3" />}
                Mark as read
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

function Skeleton() {
  return (
    <div className="space-y-2">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="card px-5 py-4 animate-pulse">
          <div className="flex gap-4">
            <div className="w-2 h-2 rounded-full bg-slate-700 mt-1.5 shrink-0" />
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-slate-800 rounded w-2/3" />
              <div className="h-3 bg-slate-800 rounded w-full" />
              <div className="h-3 bg-slate-800 rounded w-1/4" />
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading]             = useState(true)
  const [markingAll, setMarkingAll]       = useState(false)
  const { unreadCount, setUnreadCount, latestEvent, clearLatestEvent } = useNotificationStore()
  const toast = useToastStore()

  useSSE()

  const fetchAll = useCallback(async () => {
    try {
      const [notifRes, countRes] = await Promise.all([
        api.get<Notification[]>('/notifications'),
        api.get<UnreadCountResponse>('/notifications/unread-count'),
      ])
      setNotifications(notifRes.data ?? [])
      setUnreadCount(countRes.data.unread_count)
    } finally {
      setLoading(false)
    }
  }, [setUnreadCount])

  useEffect(() => { void fetchAll() }, [fetchAll])

  useEffect(() => {
    if (latestEvent) { void fetchAll(); clearLatestEvent() }
  }, [latestEvent, fetchAll, clearLatestEvent])

  const markRead = async (id: string) => {
    try {
      await api.put(`/notifications/${id}/read`)
      setNotifications((prev) => prev.map((n) => n.id === id ? { ...n, is_read: true } : n))
      setUnreadCount(Math.max(0, unreadCount - 1))
    } catch {
      toast.error('Gagal menandai notifikasi')
    }
  }

  const markAllRead = async () => {
    setMarkingAll(true)
    try {
      await api.put('/notifications/read-all')
      setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })))
      setUnreadCount(0)
      toast.success('All notifications marked as read')
    } catch {
      toast.error('Failed to mark all notifications as read')
    } finally {
      setMarkingAll(false)
    }
  }

  const unread = notifications.filter((n) => !n.is_read)
  const read   = notifications.filter((n) => n.is_read)

  return (
    <Layout>
      <div className="max-w-3xl mx-auto px-6 py-8">
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-3">
            <div className="flex items-center justify-center w-10 h-10 rounded-xl bg-amber-400/10 border border-amber-400/20">
              <BellIcon className="w-5 h-5 text-amber-400" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-slate-100">Notification</h1>
              <p className="text-sm text-slate-500">
                {unreadCount > 0 ? `${unreadCount} unread` : 'All read'}
              </p>
            </div>
          </div>
          {unreadCount > 0 && (
            <button
              onClick={() => { void markAllRead() }}
              disabled={markingAll}
              className="btn-secondary text-xs"
            >
              {markingAll ? <SpinnerIcon className="w-3 h-3" /> : <CheckIcon className="w-3 h-3" />}
              Mark all as read
            </button>
          )}
        </div>

        {loading ? (
          <Skeleton />
        ) : notifications.length === 0 ? (
          <div className="text-center py-16">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-slate-800 mb-4">
              <BellIcon className="w-8 h-8 text-slate-600" />
            </div>
            <p className="text-slate-400 font-medium">No notification</p>
            <p className="text-slate-600 text-sm mt-1">Notification will appear when there is a new announcement</p>
          </div>
        ) : (
          <div className="space-y-6">
            {unread.length > 0 && (
              <section>
                <p className="text-xs font-semibold text-amber-400 uppercase tracking-wider mb-3 px-1">
                  Unread
                </p>
                <div className="space-y-2">
                  {unread.map((n) => <NotifCard key={n.id} notif={n} onMarkRead={markRead} />)}
                </div>
              </section>
            )}
            {read.length > 0 && (
              <section>
                <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 px-1">
                  Read
                </p>
                <div className="space-y-2">
                  {read.map((n) => <NotifCard key={n.id} notif={n} onMarkRead={markRead} />)}
                </div>
              </section>
            )}
          </div>
        )}
      </div>
    </Layout>
  )
}