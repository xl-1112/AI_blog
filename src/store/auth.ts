import { create } from 'zustand'
import type { Role, User } from '../types'

const TOKEN_KEY = 'liang-blog-jwt'
const USER_KEY = 'liang-blog-user'

type AuthState = {
  token: string
  user: User | null
  setSession: (token: string, user: User) => void
  setUser: (user: User) => void
  logout: () => void
  can: (roles: Role[]) => boolean
}

function readUser() {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: localStorage.getItem(TOKEN_KEY) ?? '',
  user: readUser(),
  setSession: (token, user) => {
    localStorage.setItem(TOKEN_KEY, token)
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    set({ token, user })
  },
  setUser: (user) => {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    set({ user })
  },
  logout: () => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    set({ token: '', user: null })
  },
  can: (roles) => {
    const role = get().user?.role
    return Boolean(role && roles.includes(role))
  },
}))

