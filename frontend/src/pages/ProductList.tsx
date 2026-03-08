import React from 'react'
import { toast } from 'sonner'
import apiService from '@/services/api'
import { useCartStore } from '@/store/cartStore'
import { useAuthStore } from '@/store/authStore'

const ProductList: React.FC = () => {
  const { user, isAuthenticated } = useAuthStore()

  const [products, setProducts] = React.useState<any[]>([])
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    const load = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await apiService.getProducts()
        const data = (res.data as any)?.data || []
        const normalized = (Array.isArray(data) ? data : []).map((raw: any) => ({
          id: raw.ID ?? raw.id ?? raw.Id ?? raw._id ?? Math.random().toString(36).slice(2),
          name: raw.name ?? raw.Name ?? '',
          sku: raw.sku ?? raw.SKU ?? '',
          description: raw.description ?? raw.Description ?? '',
          price: Number(raw.price ?? raw.Price ?? 0),
          currency: raw.currency ?? raw.Currency ?? 'USD',
          stock: Number(raw.stock ?? raw.Stock ?? 0),
          category: raw.category ?? raw.Category ?? '',
          image_url: raw.image_url ?? raw.ImageURL ?? '',
        }))
        setProducts(normalized)
      } catch (e: any) {
        setError(e?.response?.data?.error || 'Failed to load products')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const { addItem } = useCartStore()
  const handleAddToCart = async (item: { sku: string; name: string; price: number }) => {
    try {
      if (!user) throw new Error('Not authenticated')
      await addItem(user.id, { ...item, quantity: 1 })
      toast.success('Added to cart')
    } catch (e: any) {
      toast.error(e?.response?.data?.error || e?.message || 'Failed to add to cart')
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Products</h1>
      {loading && (
        <div className="text-gray-600">Loading products...</div>
      )}
      {error && (
        <div className="text-red-600">{error}</div>
      )}
      {!error && !loading && products && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {products.map((p: any) => (
            <div key={p.id} className="card overflow-hidden">
              {p.image_url ? (
                <img src={p.image_url} alt={p.name} className="w-full h-40 object-cover" />
              ) : null}
              <div className="card-body">
                <div className="text-lg font-semibold text-gray-900">{p.name}</div>
                <div className="text-sm text-gray-600 mb-2">{p.category}</div>
                <div className="text-primary-600 font-bold mb-4">${((p.price || 0) / 100).toFixed(2)} {p.currency}</div>
                <div className="flex space-x-2">
                  <button
                    className="btn btn-primary"
                    disabled={!isAuthenticated}
                    onClick={() => handleAddToCart({ sku: p.sku, name: p.name, price: p.price })}
                  >
                    Add to Cart
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default ProductList