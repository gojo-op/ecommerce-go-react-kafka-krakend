import React from 'react'
import { toast } from 'sonner'
import apiService from '@/services/api'

const Chat: React.FC = () => {
  const [messages, setMessages] = React.useState<string[]>([])
  const [input, setInput] = React.useState('')
  const wsRef = React.useRef<WebSocket | null>(null)

  React.useEffect(() => {
    const ws = apiService.connectChat()
    wsRef.current = ws
    ws.onopen = () => {}
    ws.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data)
        setMessages((prev) => [...prev, JSON.stringify(data)])
      } catch {
        setMessages((prev) => [...prev, String(e.data)])
      }
    }
    ws.onerror = () => { toast.error('Chat error') }
    ws.onclose = () => {}
    return () => { ws.close() }
  }, [])

  const sendMessage = () => {
    if (!input.trim()) return
    wsRef.current?.send(JSON.stringify({ message: input }))
    setInput('')
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Chat</h1>
      <div className="card">
        <div className="card-body">
          <div className="space-y-4">
            <div className="border rounded p-4 h-64 overflow-auto bg-gray-50">
              {messages.map((m, idx) => (
                <div key={idx} className="text-sm text-gray-800">{m}</div>
              ))}
            </div>
            <div className="flex space-x-2">
              <input className="input flex-1" value={input} onChange={(e)=>setInput(e.target.value)} placeholder="Type a message" />
              <button className="btn btn-primary" onClick={sendMessage}>Send</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Chat