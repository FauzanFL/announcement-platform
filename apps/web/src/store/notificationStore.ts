import { create } from 'zustand'
import type { AnnouncementEvent } from '@/types'

interface NotificationState {
  unreadCount: number
  latestEvent: AnnouncementEvent | null
  setUnreadCount: (count: number) => void
  setLatestEvent: (event: AnnouncementEvent) => void
  clearLatestEvent: () => void
}

const useNotificationStore = create<NotificationState>((set) => ({
  unreadCount: 0,
  latestEvent: null,
  setUnreadCount: (unreadCount) => set({ unreadCount }),
  setLatestEvent: (latestEvent) => set({ latestEvent }),
  clearLatestEvent: () => set({ latestEvent: null }),
}))

export default useNotificationStore