import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import { ApiService } from './api'

describe('ApiService', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('uses authApi base URL for auth endpoints', async () => {
    const mockAuthInstance = { post: vi.fn().mockResolvedValue({ data: { data: { access_token: 'a', refresh_token: 'r', expires_at: 1, user: {} } } }), interceptors: { request: { use: vi.fn() }, response: { use: vi.fn() } } }
    const mockCoreInstance = { get: vi.fn(), post: vi.fn(), interceptors: { request: { use: vi.fn() }, response: { use: vi.fn() } } }
    vi.spyOn(axios, 'create').mockReturnValueOnce(mockCoreInstance as any).mockReturnValueOnce(mockAuthInstance as any)
    const svc = new (ApiService as any)()
    await svc.login({ email: 'e@example.com', password: 'p' })
    expect(mockAuthInstance.post).toHaveBeenCalled()
    expect((mockAuthInstance.post as any).mock.calls[0][0]).toBe('/auth/login')
  })

  it('uses authApi for refreshToken', async () => {
    const mockAuthInstance = { post: vi.fn().mockResolvedValue({ data: { data: { access_token: 'new-token' } } }), interceptors: { request: { use: vi.fn() }, response: { use: vi.fn() } } }
    const mockCoreInstance = { get: vi.fn(), post: vi.fn(), interceptors: { request: { use: vi.fn() }, response: { use: vi.fn() } } }
    vi.spyOn(axios, 'create').mockReturnValueOnce(mockCoreInstance as any).mockReturnValueOnce(mockAuthInstance as any)
    const svc = new (ApiService as any)()
    await svc.refreshToken('refresh')
    expect(mockAuthInstance.post).toHaveBeenCalledWith('/auth/refresh', { refresh_token: 'refresh' })
  })
})