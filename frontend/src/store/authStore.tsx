import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  phone?: string
  avatar?: string
  isActive: boolean
  isEmailVerified: boolean
  createdAt: string
  updatedAt: string
}

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  expiresAt: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (user: User, accessToken: string, refreshToken: string, expiresAt: string) => void
  logout: () => void
  setUser: (user: User) => void
  setLoading: (loading: boolean) => void
  clearAuth: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      expiresAt: null,
      isAuthenticated: false,
      isLoading: false,
      
      login: (user, accessToken, refreshToken, expiresAt) => {
        const normalized = {
          id: (user as any)?.ID ?? (user as any)?.id ?? '',
          email: (user as any)?.email ?? '',
          firstName: (user as any)?.firstName ?? (user as any)?.first_name ?? '',
          lastName: (user as any)?.lastName ?? (user as any)?.last_name ?? '',
          phone: (user as any)?.phone ?? '',
          avatar: (user as any)?.avatar ?? (user as any)?.avatar_url ?? '',
          isActive: (user as any)?.isActive ?? (user as any)?.is_active ?? true,
          isEmailVerified: (user as any)?.isEmailVerified ?? (user as any)?.is_email_verified ?? false,
          createdAt: (user as any)?.createdAt ?? (user as any)?.created_at ?? '',
          updatedAt: (user as any)?.updatedAt ?? (user as any)?.updated_at ?? '',
        } as any
        if (!normalized.id && accessToken) {
          const claims = parseJwt(accessToken)
          if (claims) normalized.id = claims.user_id ?? claims.sub ?? normalized.id
        }
        set({
          user: normalized,
          accessToken,
          refreshToken,
          expiresAt,
          isAuthenticated: true,
          isLoading: false,
        })
      },
      
      logout: () => {
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          expiresAt: null,
          isAuthenticated: false,
          isLoading: false,
        })
      },
      
      setUser: (user) => {
        set({ user })
      },
      
      setLoading: (loading) => {
        set({ isLoading: loading })
      },
      
      clearAuth: () => {
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          expiresAt: null,
          isAuthenticated: false,
          isLoading: false,
        })
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        expiresAt: state.expiresAt,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)

import React, { createContext, useContext, useEffect } from 'react'
import apiService from '@/services/api'
import { parseJwt } from '@/utils/jwt'

interface AuthContextType {
  isAuthenticated: boolean
  user: User | null
  login: (user: User, accessToken: string, refreshToken: string, expiresAt: string) => void
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, accessToken, isAuthenticated, isLoading, login, logout } = useAuthStore()

  useEffect(() => {
    if (accessToken && useAuthStore.getState().expiresAt) {
      const expiresAt = new Date(useAuthStore.getState().expiresAt!)
      if (expiresAt < new Date()) {
        logout()
      }
    }
  }, [accessToken, logout])

  useEffect(() => {
    const fixUser = async () => {
      if (!isAuthenticated) return
      if (!user?.id && accessToken) {
        try {
          const res = await apiService.getProfile()
          const u = (res.data as any)?.data || (res.data as any)
          const normalized = {
            id: (u?.ID ?? u?.id ?? ''),
            email: (u?.email ?? ''),
            firstName: (u?.first_name ?? u?.firstName ?? ''),
            lastName: (u?.last_name ?? u?.lastName ?? ''),
            phone: (u?.phone ?? ''),
            avatar: (u?.avatar_url ?? u?.avatar ?? ''),
            isActive: (u?.is_active ?? u?.isActive ?? true),
            isEmailVerified: (u?.is_email_verified ?? u?.isEmailVerified ?? false),
            createdAt: (u?.created_at ?? u?.createdAt ?? ''),
            updatedAt: (u?.updated_at ?? u?.updatedAt ?? ''),
          }
          useAuthStore.setState({ user: normalized as any })
        } catch {}
        if (!useAuthStore.getState().user?.id) {
          const claims = parseJwt(accessToken)
          if (claims) useAuthStore.setState({ user: { ...(useAuthStore.getState().user as any), id: claims.user_id ?? claims.sub } })
        }
      }
    }
    fixUser()
  }, [isAuthenticated, user?.id, accessToken])

  return (
    <AuthContext.Provider value={{
      isAuthenticated,
      user,
      login,
      logout,
      isLoading,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}