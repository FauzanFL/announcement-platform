import { useEffect, useState, useCallback, type FormEvent } from "react";
import Layout from "@/components/Layout";
import {
  MegaphoneIcon,
  PlusIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
  SpinnerIcon,
} from "@/components/icons";
import api from "@/lib/api";
import useToastStore from "@/store/toastStore";
import type { Announcement } from "@/types";

const formatDate = (iso: string) =>
  new Date(iso).toLocaleDateString("en-US", {
    day: "numeric",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

function StatCard({
  label,
  value,
  color,
}: {
  label: string;
  value: number;
  color: string;
}) {
  return (
    <div className="card px-5 py-4">
      <p className="text-xs text-slate-500 uppercase tracking-wider font-medium">
        {label}
      </p>
      <p className={`text-2xl sm:text-3xl font-bold mt-1 ${color}`}>{value}</p>
    </div>
  );
}

function AnnouncementModal({
  editing,
  onClose,
  onSaved,
}: {
  editing: Announcement | null;
  onClose: () => void;
  onSaved: () => void;
}) {
  const [title, setTitle] = useState(editing?.title ?? "");
  const [content, setContent] = useState(editing?.content ?? "");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const toast = useToastStore();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      if (editing) {
        await api.put(`/announcements/${editing.id}`, { title, content });
        toast.success("Announcement updated");
      } else {
        await api.post("/announcements", { title, content });
        toast.success("Announcement created", "Send to all users");
      }
      onSaved();
    } catch {
      setError("Failed to save announcement. Try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-end sm:items-center justify-center
                 bg-black/60 backdrop-blur-sm animate-fade-in"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div
        className="card w-full sm:max-w-lg animate-slide-in rounded-b-none sm:rounded-xl 
                      max-h[95dhv] flex flex-col"
      >
        <div className="flex items-center justify-between px-4 py-6 border-b border-slate-800 shrink-0">
          <div className="flex items-center gap-2.5">
            <div
              className="w-7 h-7 rounded-lg bg-brand-500/10 border border-brand-500/20
                            flex items-center justify-center"
            >
              <MegaphoneIcon className="w-3.5 h-3.5 text-brand-400" />
            </div>
            <h2 className="text-base font-semibold text-slate-100">
              {editing ? "Edit Announcement" : "New Announcement"}
            </h2>
          </div>
          <button
            onClick={onClose}
            className="text-slate-500 hover:text-slate-300 transition-colors"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        <form
          onSubmit={(e) => {
            void handleSubmit(e);
          }}
          className="p-4 sm:p-6 space-y-4 overflow-y-auto"
        >
          {error && (
            <div className="px-4 py-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
              {error}
            </div>
          )}
          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1.5">
              Title
            </label>
            <input
              className="input-field"
              placeholder="Anouncement title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              autoFocus
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1.5">
              Content
            </label>
            <textarea
              className="input-field resize-none"
              placeholder="Write your announcement here..."
              rows={5}
              value={content}
              onChange={(e) => setContent(e.target.value)}
              required
            />
          </div>
          <div className="flex items-center justify-end gap-3 pt-2 pb-safe">
            <button type="button" onClick={onClose} className="btn-secondary">
              Cancel
            </button>
            <button type="submit" disabled={loading} className="btn-primary">
              {loading ? (
                <>
                  <SpinnerIcon className="w-4 h-4" /> Saving...
                </>
              ) : editing ? (
                "Save Changes"
              ) : (
                "Publish"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

function Skeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 4 }).map((_, i) => (
        <div key={i} className="card px-4 sm:px-5 py-4 animate-pulse">
          <div className="flex items-center gap-4">
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-slate-800 rounded w-1/2" />
              <div className="h-3 bg-slate-800 rounded w-3/4" />
            </div>
            <div className="flex gap-2 shrink-0">
              <div className="h-8 w-16 bg-slate-800 rounded-lg" />
              <div className="h-8 w-16 bg-slate-800 rounded-lg" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default function AdminDashboard() {
  const [announcements, setAnnouncements] = useState<Announcement[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<Announcement | null>(null);
  const [deleting, setDeleting] = useState<string | null>(null);
  const toast = useToastStore();

  const fetchAnnouncements = useCallback(async () => {
    try {
      const res = await api.get<Announcement[]>("/announcements");
      setAnnouncements(res.data ?? []);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchAnnouncements();
  }, [fetchAnnouncements]);

  const openCreate = () => {
    setEditing(null);
    setShowModal(true);
  };
  const openEdit = (ann: Announcement) => {
    setEditing(ann);
    setShowModal(true);
  };
  const closeModal = () => {
    setShowModal(false);
    setEditing(null);
  };
  const handleSaved = () => {
    closeModal();
    void fetchAnnouncements();
  };

  const handleDelete = async (id: string) => {
    setDeleting(id);
    try {
      await api.delete(`/announcements/${id}`);
      toast.success("Announcement deleted");
      void fetchAnnouncements();
    } catch {
      toast.error("Failed to delete announcement");
    } finally {
      setDeleting(null);
    }
  };

  const now = new Date();
  const countThisMonth = announcements.filter((a) => {
    const d = new Date(a.created_at);
    return (
      d.getMonth() === now.getMonth() && d.getFullYear() === now.getFullYear()
    );
  }).length;
  const countToday = announcements.filter(
    (a) => new Date(a.created_at).toDateString() === now.toDateString(),
  ).length;

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-6 sm:py-8">
        <div className="flex items-center justify-between mb-6 sm:mb-8">
          <div className="flex items-center gap-3">
            <div className="flex items-center justify-center w-9 h-9 sm:h-10 rounded-xl bg-brand-500/10 border border-brand-500/20">
              <MegaphoneIcon className="w-4 h-4 sm:w-5 sm:h-5 text-brand-400" />
            </div>
            <div>
              <h1 className="text-lg sm:text-xl font-bold text-slate-100">
                Manage Announcement
              </h1>
              <p className="text-xs sm:text-sm text-slate-500">
                {announcements.length} active announcement
              </p>
            </div>
          </div>
          <button
            onClick={openCreate}
            className="btn-primary text-xs sm:text-sm px-3 sm:px-4"
          >
            <PlusIcon className="w-4 h-4" />
            <span className="hidden sm:inline">Create New</span>
            <span className="sm:hidden">Create</span>
          </button>
        </div>

        <div className="grid grid-cols-3 gap-2 sm:gap-3 mb-6 sm:mb-8">
          <StatCard
            label="Total"
            value={announcements.length}
            color="text-brand-400"
          />
          <StatCard
            label="This Month"
            value={countThisMonth}
            color="text-emerald-400"
          />
          <StatCard label="Today" value={countToday} color="text-amber-400" />
        </div>

        {loading ? (
          <Skeleton />
        ) : announcements.length === 0 ? (
          <div className="text-center py-16">
            <div className="inline-flex items-center justify-center w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-slate-800 mb-4">
              <MegaphoneIcon className="w-7 h-7 sm:w-8 sm:h-8 text-slate-600" />
            </div>
            <p className="text-slate-400 font-medium">No announcement</p>
            <p className="text-slate-600 text-sm mt-1 mb-4">
              Start your first announcement
            </p>
            <button onClick={openCreate} className="btn-primary mx-auto">
              <PlusIcon /> Create Your First Announcement
            </button>
          </div>
        ) : (
          <>
            <div className="hidden sm:block card overflow-hidden">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-slate-800">
                    <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                      Title
                    </th>
                    <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider hidden sm:table-cell">
                      Created At
                    </th>
                    <th className="text-right px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                      Action
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800">
                  {announcements.map((ann) => (
                    <tr
                      key={ann.id}
                      className="hover:bg-slate-800/40 transition-colors group"
                    >
                      <td className="px-5 py-4">
                        <p className="text-sm font-medium text-slate-200 group-hover:text-slate-100 line-clamp-1">
                          {ann.title}
                        </p>
                        <p className="text-xs text-slate-500 mt-0.5 line-clamp-1">
                          {ann.content}
                        </p>
                      </td>
                      <td className="px-5 py-4 hidden sm:table-cell">
                        <span className="text-xs text-slate-500 font-mono">
                          {formatDate(ann.created_at)}
                        </span>
                      </td>
                      <td className="px-5 py-4">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            onClick={() => openEdit(ann)}
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs
                                      font-medium text-slate-400 hover:text-slate-200 hover:bg-slate-700
                                      transition-colors"
                          >
                            <PencilIcon /> Edit
                          </button>
                          <button
                            onClick={() => {
                              void handleDelete(ann.id);
                            }}
                            disabled={deleting === ann.id}
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs
                                      font-medium text-red-400 hover:text-red-300 hover:bg-red-500/10
                                      transition-colors disabled:opacity-50"
                          >
                            {deleting === ann.id ? (
                              <SpinnerIcon className="w-3 h-3" />
                            ) : (
                              <TrashIcon />
                            )}
                            Delete
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="sm:hidden space-y-3">
              {announcements.map((ann) => (
                <div key={ann.id} className="card px-4 py-4">
                  <p className="text-sm font-semibold text-slate-100 line-clamp-1 mb-1">
                    {ann.title}
                  </p>
                  <p className="text-xs text-slate-500 line-clamp-2 mb-3">
                    {ann.content}
                  </p>
                  <div className="flex items-center justify-between">
                    <span className="text-xs text-slate-600 font-mono">
                      {formatDate(ann.created_at)}
                    </span>
                    <div className="flex gap-2">
                      <button
                        onClick={() => openEdit(ann)}
                        className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs
                                   font-medium text-slate-400 hover:text-slate-200 hover:bg-slate-700
                                   transition-colors border border-slate-700"
                      >
                        <PencilIcon /> Edit
                      </button>
                      <button
                        onClick={() => {
                          void handleDelete(ann.id);
                        }}
                        disabled={deleting === ann.id}
                        className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs
                                   font-medium text-red-400 hover:bg-red-500/10
                                   transition-colors border border-red-500/20 disabled:opacity-50"
                      >
                        {deleting === ann.id ? (
                          <SpinnerIcon className="w-3 h-3" />
                        ) : (
                          <TrashIcon />
                        )}
                        Hapus
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </>
        )}
      </div>

      {showModal && (
        <AnnouncementModal
          editing={editing}
          onClose={closeModal}
          onSaved={handleSaved}
        />
      )}
    </Layout>
  );
}
