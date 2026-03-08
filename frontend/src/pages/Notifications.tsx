import React from 'react'
import { useAuthStore } from '@/store/authStore'
import apiService from '@/services/api'
import { toast } from 'sonner'

const Notifications: React.FC = () => {
  const { user } = useAuthStore()
  const [items, setItems] = React.useState<any[]>([])

  const load = async () => {
    if (!user) return
    try {
      const res = await apiService.getNotifications(user.id)
      const data = (res.data as any)?.data || []
      setItems(data)
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Failed to load notifications')
    }
  }

  React.useEffect(() => { load() }, [user?.id])

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Notifications</h1>
      <div className="space-y-4">
        {items.length === 0 && <div className="text-gray-600">No notifications</div>}
        {items.map((n, idx) => (
          <div key={idx} className="card">
            <div className="card-body">
              <div className="font-medium">{String(n.type)}</div>
              <pre className="text-sm text-gray-700 whitespace-pre-wrap">{JSON.stringify(n.data)}</pre>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default Notifications