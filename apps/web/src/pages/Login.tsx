import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router'
import api from '@/lib/api'
import { MegaphoneIcon, SpinnerIcon } from '@/components/icons'
import type { AuthResponse } from '@/types'
import useAuthStore from '@/store/authStore'

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError]       = useState('')
  const [loading, setLoading]   = useState(false)
  const login    = useAuthStore((s) => s.login)
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await api.post<AuthResponse>('/login', { username, password })
      login(res.data.token, res.data.user)
      void navigate(res.data.user.role === 'admin' ? '/admin' : '/')
    } catch {
      setError('Login failed. Check your username and password.')
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
          <h1 className="text-2xl font-bold text-slate-100">Login to Information Board</h1>
          <p className="text-slate-500 text-sm mt-1">Centralized announcement board</p>
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
                placeholder="Enter username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
                autoFocus
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Password</label>
              <input
                type="password"
                className="input-field"
                placeholder="Enter password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <button type="submit" disabled={loading} className="btn-primary w-full justify-center py-2.5">
              {loading ? <><SpinnerIcon className="w-4 h-4" /> Processing...</> : 'Login'}
            </button>
          </form>
        </div>

        <p className="text-center text-sm text-slate-500 mt-4">
          Don't have an account?{' '}
          <Link to="/register" className="text-brand-400 hover:text-brand-300 font-medium transition-colors">
            Register here
          </Link>
        </p>
      </div>
    </div>
  )
}