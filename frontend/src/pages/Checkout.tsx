import React from 'react'
import { useAuthStore } from '@/store/authStore'
import { useCartStore } from '@/store/cartStore'
import apiService from '@/services/api'
import { toast } from 'sonner'
import { useNavigate } from 'react-router-dom'

const Checkout: React.FC = () => {
  const { user } = useAuthStore()
  const { items, count, loadCart } = useCartStore()
  const [addresses, setAddresses] = React.useState<any[]>([])
  const [selectedId, setSelectedId] = React.useState<string>('')
  const [saving, setSaving] = React.useState(false)
  const [paymentMethod, setPaymentMethod] = React.useState<'online'|'cod'>('online')
  const navigate = useNavigate()

  const totalCents = items.reduce((sum, it) => sum + (it.price || 0) * (it.quantity || 0), 0)

  const fetchAddresses = async () => {
    try {
      const res = await apiService.getAddresses()
      const payload = (res.data as any)?.data || (res.data as any)
      setAddresses(payload || [])
      const def = (payload || []).find((a: any) => a.is_default)
      if (def) setSelectedId(def.id || def.ID)
    } catch {}
  }

  React.useEffect(() => {
    if (user?.id) {
      loadCart(user.id)
      fetchAddresses()
    }
  }, [user?.id])

  const [form, setForm] = React.useState({
    name: '', phone: '', address1: '', address2: '', city: '', state: '', country: '', postal_code: '', is_default: false,
  })

  const handleAddAddress = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      const res = await apiService.createAddress({
        first_name: form.name.split(' ')[0] || form.name,
        last_name: form.name.split(' ').slice(1).join(' '),
        phone: form.phone,
        address1: form.address1,
        address2: form.address2,
        city: form.city,
        state: form.state,
        country: form.country,
        postal_code: form.postal_code,
        is_default: form.is_default,
        type: 'shipping',
      })
      const payload = (res.data as any)?.data || (res.data as any)
      toast.success('Address saved')
      setSelectedId(payload?.id || payload?.ID)
      fetchAddresses()
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Failed to save address')
    } finally { setSaving(false) }
  }

  const handleCheckout = async () => {
    if (!user?.id) { toast.error('Please login'); return }
    if (count === 0) { toast.error('Cart is empty'); return }
    const addr = addresses.find(a => (a.id || a.ID) === selectedId)
    if (!addr) { toast.error('Select or add an address'); return }
    try {
      const res = await apiService.checkout({
        user_id: user.id,
        items: items.map(it => ({ sku: it.sku, name: it.name, unit_price: it.price, quantity: it.quantity })),
        currency: 'USD',
        payment_method: paymentMethod,
        shipping: {
          name: `${addr.first_name || ''} ${addr.last_name || ''}`.trim(),
          phone: addr.phone || '',
          address1: addr.address1,
          address2: addr.address2 || '',
          city: addr.city,
          state: addr.state,
          country: addr.country,
          postal: addr.postal_code,
        }
      })
      const order = (res.data as any)?.data || (res.data as any)
      toast.success('Order created')
      if (paymentMethod === 'cod') {
        navigate('/orders')
      } else {
        navigate('/payments', { state: { order_id: order?.id || order?.ID, amount: order?.total, currency: order?.currency } })
      }
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Checkout failed')
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Checkout</h1>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <div className="card-body">
            <h2 className="text-xl font-semibold mb-4">Shipping Address</h2>
            {addresses.length > 0 && (
              <div className="space-y-3 mb-6">
                {addresses.map((a) => (
                  <label key={a.id || a.ID} className="flex items-start space-x-3">
                    <input type="radio" name="address" checked={selectedId === (a.id || a.ID)} onChange={() => setSelectedId(a.id || a.ID)} />
                    <div>
                      <div className="font-medium">{a.first_name} {a.last_name} {a.is_default ? '(Default)' : ''}</div>
                      <div className="text-sm text-gray-600">{a.address1}{a.address2 ? ', ' + a.address2 : ''}</div>
                      <div className="text-sm text-gray-600">{a.city}, {a.state}, {a.country} {a.postal_code}</div>
                      <div className="text-sm text-gray-600">{a.phone}</div>
                    </div>
                  </label>
                ))}
              </div>
            )}
            <form onSubmit={handleAddAddress} className="space-y-3">
              <input className="input" placeholder="Full Name" value={form.name} onChange={e=>setForm({...form, name: e.target.value})} />
              <input className="input" placeholder="Phone" value={form.phone} onChange={e=>setForm({...form, phone: e.target.value})} />
              <input className="input" placeholder="Address Line 1" value={form.address1} onChange={e=>setForm({...form, address1: e.target.value})} />
              <input className="input" placeholder="Address Line 2" value={form.address2} onChange={e=>setForm({...form, address2: e.target.value})} />
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <input className="input" placeholder="City" value={form.city} onChange={e=>setForm({...form, city: e.target.value})} />
                <input className="input" placeholder="State" value={form.state} onChange={e=>setForm({...form, state: e.target.value})} />
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <input className="input" placeholder="Country" value={form.country} onChange={e=>setForm({...form, country: e.target.value})} />
                <input className="input" placeholder="Postal Code" value={form.postal_code} onChange={e=>setForm({...form, postal_code: e.target.value})} />
              </div>
              <label className="flex items-center space-x-2">
                <input type="checkbox" checked={form.is_default} onChange={e=>setForm({...form, is_default: e.target.checked})} />
                <span>Set as default</span>
              </label>
              <button className="btn btn-secondary" type="submit" disabled={saving}>{saving ? 'Saving...' : 'Save Address'}</button>
            </form>
          </div>
        </div>
        <div className="card">
          <div className="card-body">
            <h2 className="text-xl font-semibold mb-4">Order Summary</h2>
            <div className="space-y-3">
              {items.map(it => (
                <div key={it.sku} className="flex items-center justify-between">
                  <div>
                    <div className="font-medium">{it.name}</div>
                    <div className="text-sm text-gray-600">Qty: {it.quantity}</div>
                  </div>
                  <div className="text-right font-semibold">${(it.price/100).toFixed(2)}</div>
                </div>
              ))}
            </div>
            <div className="flex items-center justify-between mt-4">
              <div className="text-xl font-bold text-gray-900">Total: ${ (totalCents / 100).toFixed(2) }</div>
              <button className="btn btn-primary" onClick={handleCheckout}>Pay Now</button>
            </div>
            <div className="mt-3">
              <label className="font-medium mr-3">Payment Method:</label>
              <label className="mr-4">
                <input type="radio" name="pm" checked={paymentMethod==='online'} onChange={()=>setPaymentMethod('online')} /> <span>Pay Online</span>
              </label>
              <label>
                <input type="radio" name="pm" checked={paymentMethod==='cod'} onChange={()=>setPaymentMethod('cod')} /> <span>Cash on Delivery (COD)</span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Checkout