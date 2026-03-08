import axios, { AxiosInstance, AxiosError, AxiosResponse } from 'axios'
import { toast } from 'sonner'
import { useAuthStore } from '@/store/authStore'

export class ApiService {
  private api: AxiosInstance

  constructor() {
    this.api = axios.create({
      baseURL: import.meta.env.VITE_API_URL || '/api/v1',
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor
    const attachAuthToken = (config: any) => {
      const token = useAuthStore.getState().accessToken
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    }

    this.api.interceptors.request.use(
      attachAuthToken,
      (error) => {
        return Promise.reject(error)
      }
    )


    // Response interceptor
    const handleResponseError = async (error: AxiosError) => {
      const originalRequest: any = error.config

      if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
        originalRequest._retry = true
        try {
          const refreshToken = useAuthStore.getState().refreshToken
          if (refreshToken) {
            const response = await this.api.post('/auth/refresh', { refresh_token: refreshToken })
            const { access_token } = (response.data as any)?.data || (response.data as any)
            useAuthStore.setState({ accessToken: access_token })
            originalRequest.headers = originalRequest.headers || {}
            originalRequest.headers.Authorization = `Bearer ${access_token}`
            return axios(originalRequest)
          }
        } catch (refreshError) {
          useAuthStore.getState().logout()
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      }

      if (error.response) {
        const message = (error.response.data as any)?.error || 'An error occurred'
        switch (error.response.status) {
          case 400:
          case 422:
            toast.error(message)
            break
          case 401:
            toast.error('Please login to continue')
            break
          case 403:
            toast.error('You do not have permission to perform this action')
            break
          case 404:
            toast.error('Resource not found')
            break
          case 500:
            toast.error('Internal server error. Please try again later.')
            break
          default:
            toast.error(message)
        }
      } else if (error.request) {
        toast.error('Network error. Please check your connection.')
      } else {
        toast.error('An unexpected error occurred')
      }

      return Promise.reject(error)
    }

    this.api.interceptors.response.use(
      (response: AxiosResponse) => {
        return response
      },
      handleResponseError
    )

  }

  // Auth endpoints
  async register(data: any) {
    return this.api.post('/auth/register', data)
  }

  async login(data: any) {
    return this.api.post('/auth/login', data)
  }

  async refreshToken(refreshToken: string) {
    return this.api.post('/auth/refresh', { refresh_token: refreshToken })
  }

  async getProfile() {
    return this.api.get('/auth/profile')
  }

  // Address endpoints
  async getAddresses() {
    return this.api.get('/auth/addresses')
  }
  async createAddress(data: any) {
    return this.api.post('/auth/addresses', data)
  }
  async updateAddress(id: string, data: any) {
    return this.api.put(`/auth/addresses/${id}`, data)
  }
  async deleteAddress(id: string) {
    return this.api.delete(`/auth/addresses/${id}`)
  }

  // Product endpoints
  async getProducts(params?: any) {
    return this.api.get('/products', { params })
  }

  async getProduct(id: string) {
    return this.api.get(`/products/${id}`)
  }

  async getProductBySKU(sku: string) {
    return this.api.get(`/products/sku/${sku}`)
  }

  async createProduct(data: any) {
    return this.api.post('/products', data)
  }

  async updateProduct(id: string, data: any) {
    return this.api.put(`/products/${id}`, data)
  }

  async deleteProduct(id: string) {
    return this.api.delete(`/products/${id}`)
  }

  async updateProductStock(id: string, quantity: number) {
    return this.api.patch(`/products/${id}/stock`, { quantity })
  }

  // Cart endpoints
  async getCart(userId: string) {
    return this.api.get(`/cart/${userId}`)
  }
  async addToCart(userId: string, item: { sku: string; name: string; price: number; quantity: number }) {
    return this.api.post(`/cart/${userId}/items`, item)
  }
  async removeFromCart(userId: string, sku: string) {
    return this.api.delete(`/cart/${userId}/items/${sku}`)
  }
  async clearCart(userId: string) {
    return this.api.delete(`/cart/${userId}`)
  }

  // Order endpoints
  async createOrder(data: any) {
    return this.api.post('/orders', data)
  }

  async checkout(data: any) {
    return this.api.post('/orders/checkout', data)
  }

  async getOrders(params?: any) {
    return this.api.get('/orders', { params })
  }

  async getOrder(id: string) {
    return this.api.get(`/orders/${id}`)
  }

  async updateOrderStatus(id: string, status: string) {
    return this.api.patch(`/orders/${id}/status`, { status })
  }

  async cancelOrder(id: string) {
    return this.api.patch(`/orders/${id}/status`, { status: 'cancelled' })
  }

  // Payment endpoints
  async createPaymentIntent(data: any) {
    return this.api.post('/payments/intent', data)
  }

  // Chat helpers
  connectChat() {
    const origin = window.location.origin
    const wsUrl = origin.replace('http', 'ws') + '/ws'
    return new WebSocket(wsUrl)
  }

  // Notifications
  async getNotifications(userId: string) {
    return this.api.get('/notifications', { params: { user_id: userId } })
  }
}

export const apiService = new ApiService()
export default apiService