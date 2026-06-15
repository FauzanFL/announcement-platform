import { useState, type FormEvent, type ChangeEvent } from 'react'
import { Link, useNavigate } from 'react-router'
import api from '@/lib/api'
import { MegaphoneIcon, SpinnerIcon } from '@/components/icons'
import type { Role } from '@/types'

interface RegisterForm { username: string; password: string; role: Role }

export default function Register() {
  const [form, setForm]       = useState<RegisterForm>({ username: '', password: '', role: 'user' })
  const [error, setError]     = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleChange =
    (key: keyof RegisterForm) =>
    (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
      setForm((f) => ({ ...f, [key]: e.target.value }))

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await api.post('/register', form)
      void navigate('/login')
    } catch {
      setError('Pendaftaran gagal. Username mungkin sudah digunakan.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-slate-950 flex items-center justify-center px-4">
      <div className="w-full max-w-md animate-fade-in">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-2xl bg-brand-500 mb-4">
            <MegaphoneIcon className="w-6 h-6 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-slate-100">Buat Akun</h1>
          <p className="text-slate-500 text-sm mt-1">Daftarkan diri untuk mulai menerima pengumuman</p>
        </div>

        <div className="card p-6">
          {error && (
            <div className="mb-4 px-4 py-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
              {error}
            </div>
          )}
          <form onSubmit={(e) => { void handleSubmit(e) }} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Username</label>
              <input
                className="input-field"
                placeholder="Pilih username"
                value={form.username}
                onChange={handleChange('username')}
                required
                autoFocus
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Password</label>
              <input
                type="password"
                className="input-field"
                placeholder="Min. 6 karakter"
                value={form.password}
                onChange={handleChange('password')}
                required
                minLength={6}
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Role</label>
              <select className="input-field" value={form.role} onChange={handleChange('role')}>
                <option value="user">User — menerima pengumuman</option>
                <option value="admin">Admin — mengelola pengumuman</option>
              </select>
            </div>
            <button type="submit" disabled={loading} className="btn-primary w-full justify-center py-2.5">
              {loading ? <><SpinnerIcon className="w-4 h-4" /> Mendaftarkan...</> : 'Daftar'}
            </button>
          </form>
        </div>

        <p className="text-center text-sm text-slate-500 mt-4">
          Sudah punya akun?{' '}
          <Link to="/login" className="text-brand-400 hover:text-brand-300 font-medium transition-colors">
            Masuk di sini
          </Link>
        </p>
      </div>
    </div>
  )
}