import { create } from 'zustand'
import type { ToastItem, ToastType } from '@/types'

let _id = 0

interface ToastState {
  toasts: ToastItem[]
  push: (message: string, type: ToastType, title?: string, duration?: number) => void
  dismiss: (id: number) => void
  success: (message: string, title?: string) => void
  error: (message: string, title?: string) => void
  info: (message: string, title?: string) => void
}

const useToastStore = create<ToastState>((set) => {
  const add = (message: string, type: ToastType, title?: string, duration?: number) => {
    const id = ++_id
    set((s) => ({ toasts: [...s.toasts, { id, message, type, title, duration }] }))
  }

  return {
    toasts: [],
    push: add,
    dismiss: (id) => set((s) => ({ toasts: s.toasts.filter((t) => t.id !== id) })),
    success: (message, title) => add(message, 'success', title),
    error:   (message, title) => add(message, 'error',   title),
    info:    (message, title) => add(message, 'info',    title),
  }
})

export default useToastStore