import { useEffect, useState } from 'react'
import { CheckIcon, XMarkIcon, BellIcon } from './icons'
import useToastStore from '@/store/toastStore'
import type { ToastItem, ToastType } from '@/types'

const VARIANTS: Record<ToastType, { bar: string; icon: React.ReactNode; label: string }> = {
  success: {
    bar:   'bg-emerald-500',
    icon:  <CheckIcon className="w-5 h-5 text-emerald-400" />,
    label: 'text-emerald-400',
  },
  error: {
    bar:   'bg-red-500',
    icon:  <XMarkIcon className="w-5 h-5 text-red-400" />,
    label: 'text-red-400',
  },
  info: {
    bar:   'bg-brand-500',
    icon:  <BellIcon className="w-5 h-5 text-brand-400" />,
    label: 'text-brand-400',
  },
}

interface ToastItemProps {
  toast: ToastItem
  onDismiss: (id: number) => void
}

function ToastCard({ toast, onDismiss }: ToastItemProps) {
  const [visible, setVisible] = useState(false)
  const variant = VARIANTS[toast.type]

  useEffect(() => {
    requestAnimationFrame(() => setVisible(true))
    const timer = setTimeout(() => {
      setVisible(false)
      setTimeout(() => onDismiss(toast.id), 300)
    }, toast.duration ?? 4000)
    return () => clearTimeout(timer)
  }, [toast.id, toast.duration, onDismiss])

  return (
    <div
      className={`pointer-events-auto relative overflow-hidden rounded-xl border border-slate-700
                  bg-slate-900 shadow-2xl shadow-black/40 transition-all duration-300
                  ${visible ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-2'}`}
    >
      <div className={`absolute left-0 top-0 bottom-0 w-1 ${variant.bar}`} />
      <div className="flex items-start gap-3 px-4 py-3 pl-5">
        <div className="mt-0.5 shrink-0">{variant.icon}</div>
        <div className="flex-1 min-w-0">
          {toast.title && (
            <p className={`text-xs font-semibold uppercase tracking-wider mb-0.5 ${variant.label}`}>
              {toast.title}
            </p>
          )}
          <p className="text-sm text-slate-200 leading-snug">{toast.message}</p>
        </div>
        <button
          onClick={() => onDismiss(toast.id)}
          className="shrink-0 text-slate-500 hover:text-slate-300 transition-colors mt-0.5"
        >
          <XMarkIcon className="w-4 h-4" />
        </button>
      </div>
    </div>
  )
}

export default function Toast() {
  const toasts  = useToastStore((s) => s.toasts)
  const dismiss = useToastStore((s) => s.dismiss)

  return (
    <div className="fixed top-4 right-4 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
      {toasts.map((t) => (
        <ToastCard key={t.id} toast={t} onDismiss={dismiss} />
      ))}
    </div>
  )
}