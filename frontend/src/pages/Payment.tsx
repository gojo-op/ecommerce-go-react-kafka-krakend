import React from 'react'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import apiService from '@/services/api'

type FormData = {
  orderId: string
  amount: number
  currency: string
  provider: 'stripe' | 'razorpay'
}

const Payment: React.FC = () => {
  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    defaultValues: { currency: 'USD', provider: 'stripe' }
  })

  const onSubmit = async (data: FormData) => {
    try {
      await apiService.createPaymentIntent({
        order_id: data.orderId,
        amount: Number(data.amount),
        currency: data.currency,
        provider: data.provider,
      })
      toast.success('Payment intent created')
      reset()
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Failed to create payment intent')
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Payments</h1>
      <div className="card">
        <div className="card-body">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 max-w-lg">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Order ID</label>
              <input className="input" {...register('orderId', { required: 'Order ID is required' })} placeholder="Order ID" />
              {errors.orderId && <p className="mt-1 text-sm text-red-600">{errors.orderId.message}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Amount</label>
              <input type="number" className="input" {...register('amount', { required: 'Amount is required', min: 1 })} placeholder="Amount" />
              {errors.amount && <p className="mt-1 text-sm text-red-600">{errors.amount.message}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Currency</label>
              <input className="input" {...register('currency', { required: 'Currency is required' })} placeholder="USD" />
              {errors.currency && <p className="mt-1 text-sm text-red-600">{errors.currency.message}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Provider</label>
              <select className="input" {...register('provider', { required: 'Provider is required' })}>
                <option value="stripe">Stripe</option>
                <option value="razorpay">Razorpay</option>
              </select>
            </div>
            <button type="submit" className="btn btn-primary">Create Intent</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default Payment