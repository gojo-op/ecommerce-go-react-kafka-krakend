import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import apiService from '@/services/api'

export interface CartItem {
  sku: string
  name: string
  price: number
  quantity: number
}

interface CartState {
  items: CartItem[]
  count: number
  loading: boolean
  error: string | null
  loadCart: (userId: string) => Promise<void>
  addItem: (userId: string, item: CartItem) => Promise<void>
  removeItem: (userId: string, sku: string) => Promise<void>
  clear: (userId: string) => Promise<void>
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      count: 0,
      loading: false,
      error: null,
      async loadCart(userId: string) {
        set({ loading: true, error: null })
        try {
          const res = await apiService.getCart(userId)
          const cart = (res.data as any)?.data || (res.data as any)
          const items: CartItem[] = Array.isArray(cart?.items) ? cart.items.map((it: any) => ({
            sku: it.sku ?? it.SKU,
            name: it.name ?? it.Name,
            price: Number(it.price ?? it.Price ?? 0),
            quantity: Number(it.quantity ?? it.Quantity ?? 0),
          })) : []
          const count = items.reduce((sum, it) => sum + (it.quantity || 0), 0)
          set({ items, count, loading: false })
        } catch (e: any) {
          set({ error: e?.response?.data?.error || 'Failed to load cart', loading: false })
        }
      },
      async addItem(userId: string, item: CartItem) {
        set({ loading: true, error: null })
        try {
          await apiService.addToCart(userId, item)
          await get().loadCart(userId)
        } catch (e: any) {
          set({ error: e?.response?.data?.error || 'Failed to add item', loading: false })
          throw e
        }
      },
      async removeItem(userId: string, sku: string) {
        set({ loading: true, error: null })
        try {
          await apiService.removeFromCart(userId, sku)
          await get().loadCart(userId)
        } catch (e: any) {
          set({ error: e?.response?.data?.error || 'Failed to remove item', loading: false })
        }
      },
      async clear(userId: string) {
        set({ loading: true, error: null })
        try {
          await apiService.clearCart(userId)
          set({ items: [], count: 0, loading: false })
        } catch (e: any) {
          set({ error: e?.response?.data?.error || 'Failed to clear cart', loading: false })
        }
      },
    }),
    {
      name: 'cart-storage',
      partialize: (state) => ({ items: state.items, count: state.count }),
    }
  )
)