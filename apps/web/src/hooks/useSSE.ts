import { useEffect, useRef } from 'react'
import useAuthStore from '@/store/authStore'
import useNotificationStore from '@/store/notificationStore'
import useToastStore from '@/store/toastStore'
import type { AnnouncementEvent, UnreadCountResponse } from '@/types'

const BASE_URL = import.meta.env.VITE_API_URL

export default function useSSE(): void {
  const token  = useAuthStore((s) => s.token)
  const user   = useAuthStore((s) => s.user)
  const setUnreadCount = useNotificationStore((s) => s.setUnreadCount)
  const setLatestEvent = useNotificationStore((s) => s.setLatestEvent)
  const toastInfo      = useToastStore((s) => s.info)
  const esRef = useRef<EventSource | null>(null)

  useEffect(() => {
    if (!token || user?.role !== 'user') return

    const url = `${BASE_URL}/stream?token=${encodeURIComponent(token)}`
    const es = new EventSource(url)
    esRef.current = es

    es.addEventListener('announcement', (e: MessageEvent<string>) => {
      const data = JSON.parse(e.data) as AnnouncementEvent
      setLatestEvent(data)

      if (data.type === 'created' && data.announcement) {
        toastInfo(data.announcement.title, 'Pengumuman Baru')
        if (Notification.permission === 'granted') {
          new Notification('Pengumuman Baru', { body: data.announcement.title })
        }
      }
    })

    es.addEventListener('unread_count', (e: MessageEvent<string>) => {
      const data = JSON.parse(e.data) as UnreadCountResponse
      setUnreadCount(data.unread_count)
    })

    es.onerror = () => {
      console.warn('SSE reconnecting...')
    }

    if (Notification.permission === 'default') {
      void Notification.requestPermission()
    }

    return () => {
      es.close()
      esRef.current = null
    }
  }, [token, user?.role])
}