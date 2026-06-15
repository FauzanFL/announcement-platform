import { useEffect, useState } from 'react'
import Layout from '@/components/Layout'
import { ChevronRightIcon, MegaphoneIcon } from '@/components/icons'
import api from '@/lib/api'
import useNotificationStore from '@/store/notificationStore'
import useSSE from '@/hooks/useSSE'
import type { Announcement } from '@/types'

const formatDate = (iso: string) =>
  new Date(iso).toLocaleDateString('en-US', {
    day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })

function AnnouncementCard({ ann }: { ann: Announcement }) {
  const [expanded, setExpanded] = useState(false)
  return (
    <div className="card overflow-hidden animate-slide-in hover:border-slate-700 transition-colors duration-150">
      <button
        className="w-full text-left px-5 py-4 flex items-start gap-4 group"
        onClick={() => setExpanded((v) => !v)}
      >
        <div className="mt-1.5 w-2 h-2 rounded-full bg-brand-400 shrink-0 ring-4 ring-brand-400/10" />
        <div className="flex-1 min-w-0">
          <p className="text-sm font-semibold text-slate-100 group-hover:text-brand-300 transition-colors">
            {ann.title}
          </p>
          <p className="text-xs text-slate-500 mt-0.5 font-mono">{formatDate(ann.created_at)}</p>
        </div>
        <ChevronRightIcon
          className={`w-4 h-4 text-slate-500 shrink-0 mt-0.5 transition-transform duration-200
                      ${expanded ? 'rotate-90' : ''}`}
        />
      </button>
      {expanded && (
        <div className="px-5 pb-5 animate-fade-in">
          <div className="ml-6 pt-3 border-t border-slate-800">
            <p className="text-sm text-slate-300 leading-relaxed whitespace-pre-wrap">{ann.content}</p>
          </div>
        </div>
      )}
    </div>
  )
}

function Skeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 4 }).map((_, i) => (
        <div key={i} className="card px-5 py-4 animate-pulse">
          <div className="flex items-start gap-4">
            <div className="w-2 h-2 rounded-full bg-slate-700 mt-1.5 shrink-0" />
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-slate-800 rounded w-3/4" />
              <div className="h-3 bg-slate-800 rounded w-1/3" />
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

export default function AnnouncementsPage() {
  const [announcements, setAnnouncements] = useState<Announcement[]>([])
  const [loading, setLoading]             = useState(true)
  const { latestEvent, clearLatestEvent } = useNotificationStore()

  useSSE()

  const fetchAnnouncements = async () => {
    try {
      const res = await api.get<Announcement[]>('/announcements')
      setAnnouncements(res.data ?? [])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void fetchAnnouncements() }, [])

  useEffect(() => {
    if (latestEvent) { void fetchAnnouncements(); clearLatestEvent() }
  }, [latestEvent, clearLatestEvent])

  return (
    <Layout>
      <div className="max-w-3xl mx-auto px-6 py-8">
        <div className="flex items-center gap-3 mb-8">
          <div className="flex items-center justify-center w-10 h-10 rounded-xl bg-brand-500/10 border border-brand-500/20">
            <MegaphoneIcon className="w-5 h-5 text-brand-400" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-slate-100">Announcement</h1>
            <p className="text-sm text-slate-500">{announcements.length} announcement available</p>
          </div>
        </div>

        {loading ? (
          <Skeleton />
        ) : announcements.length === 0 ? (
          <div className="text-center py-16">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-slate-800 mb-4">
              <MegaphoneIcon className="w-8 h-8 text-slate-600" />
            </div>
            <p className="text-slate-400 font-medium">No announcement</p>
            <p className="text-slate-600 text-sm mt-1">New announcement will appear here automatically</p>
          </div>
        ) : (
          <div className="space-y-3">
            {announcements.map((a) => <AnnouncementCard key={a.id} ann={a} />)}
          </div>
        )}
      </div>
    </Layout>
  )
}