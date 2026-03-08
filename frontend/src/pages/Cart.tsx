import React from 'react'
import { useAuthStore } from '@/store/authStore'
import { useCartStore } from '@/store/cartStore'
import { toast } from 'sonner'

const Cart: React.FC = () => {
  const { user } = useAuthStore()
  const { items, count, loading, error, loadCart, removeItem, clear } = useCartStore()

  React.useEffect(() => {
    if (user?.id) {
      loadCart(user.id)
    }
  }, [user?.id, loadCart])

  const handleRemove = async (sku: string) => {
    if (!user?.id) return
    try {
      await removeItem(user.id, sku)
      toast.success('Item removed')
    } catch {}
  }

  const handleClear = async () => {
    if (!user?.id) return
    await clear(user.id)
    toast.success('Cart cleared')
  }

  const totalCents = items.reduce((sum, it) => sum + (it.price || 0) * (it.quantity || 0), 0)

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Shopping Cart</h1>
      {loading && <div className="text-gray-600">Loading cart...</div>}
      {error && <div className="text-red-600">{error}</div>}
      {!loading && !error && (
        <div>
          <div className="flex items-center justify-between mb-4">
            <div className="text-gray-700">Items: <span className="font-semibold">{count}</span></div>
            <button className="btn btn-secondary" onClick={handleClear} disabled={items.length === 0}>Clear Cart</button>
          </div>
          {items.length === 0 ? (
            <p className="text-gray-600">Your cart is empty.</p>
          ) : (
            <div className="space-y-4">
              {items.map((it) => (
                <div key={it.sku} className="card">
                  <div className="card-body flex items-center justify-between">
                    <div>
                      <div className="font-semibold text-gray-900">{it.name}</div>
                      <div className="text-gray-600 text-sm">SKU: {it.sku}</div>
                      <div className="text-gray-600 text-sm">Qty: {it.quantity}</div>
                    </div>
                    <div className="text-right">
                      <div className="text-primary-600 font-bold mb-2">${((it.price || 0) / 100).toFixed(2)}</div>
                      <button className="btn btn-danger" onClick={() => handleRemove(it.sku)}>Remove</button>
                    </div>
                  </div>
                </div>
              ))}
              <div className="flex items-center justify-between mt-4">
                <div className="text-xl font-bold text-gray-900">Total: ${ (totalCents / 100).toFixed(2) }</div>
                <a href="/checkout" className="btn btn-primary">Proceed to Checkout</a>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default Cart