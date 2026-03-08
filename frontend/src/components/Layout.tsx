import React from 'react'
import { Link, Outlet } from 'react-router-dom'
import { useAuth } from '@/store/authStore'
import { useCartStore } from '@/store/cartStore'
import { ShoppingCart, User, LogOut, Package, Search, Menu, X } from 'lucide-react'
import { toast } from 'sonner'

const Layout: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth()
  const { count, loadCart } = useCartStore()
  const [isMenuOpen, setIsMenuOpen] = React.useState(false)
  const [searchQuery, setSearchQuery] = React.useState('')

  React.useEffect(() => {
    if (isAuthenticated && user?.id) {
      loadCart(user.id)
    }
  }, [isAuthenticated, user?.id, loadCart])

  const handleLogout = () => {
    logout()
    toast.success('Logged out successfully')
  }

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      // Navigate to search results
      window.location.href = `/products?search=${encodeURIComponent(searchQuery)}`
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center space-x-2">
              <Package className="h-8 w-8 text-primary-600" />
              <span className="text-xl font-bold text-gray-900">DemoShop</span>
            </Link>

            {/* Search Bar */}
            <form onSubmit={handleSearch} className="hidden md:flex flex-1 max-w-lg mx-8">
              <div className="relative w-full">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Search products..."
                  className="input pl-10 pr-4 py-2"
                />
              </div>
            </form>

            {/* Navigation */}
            <nav className="hidden md:flex items-center space-x-6">
              <Link to="/products" className="text-gray-700 hover:text-primary-600 font-medium">
                Products
              </Link>
              
              {isAuthenticated ? (
                <>
                  <Link to="/orders" className="text-gray-700 hover:text-primary-600 font-medium">
                    Orders
                  </Link>
                  <Link to="/payments" className="text-gray-700 hover:text-primary-600 font-medium">
                    Payments
                  </Link>
                  <Link to="/chat" className="text-gray-700 hover:text-primary-600 font-medium">
                    Chat
                  </Link>
                  <Link to="/notifications" className="text-gray-700 hover:text-primary-600 font-medium">
                    Notifications
                  </Link>
                  <div className="flex items-center space-x-4">
                    <Link to="/cart" className="relative text-gray-700 hover:text-primary-600">
                      <ShoppingCart className="h-6 w-6" />
                      <span className="absolute -top-1 -right-1 bg-primary-600 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                        {count}
                      </span>
                    </Link>
                    <div className="relative group">
                      <button className="flex items-center space-x-2 text-gray-700 hover:text-primary-600">
                        <User className="h-6 w-6" />
                        <span className="font-medium">{user?.firstName}</span>
                      </button>
                      <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-10 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200">
                        <Link to="/profile" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                          Profile
                        </Link>
                        <button
                          onClick={handleLogout}
                          className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        >
                          <div className="flex items-center space-x-2">
                            <LogOut className="h-4 w-4" />
                            <span>Logout</span>
                          </div>
                        </button>
                      </div>
                    </div>
                  </div>
                </>
              ) : (
                <div className="flex items-center space-x-4">
                  <Link to="/login" className="text-gray-700 hover:text-primary-600 font-medium">
                    Login
                  </Link>
                  <Link to="/register" className="btn btn-primary">
                    Sign Up
                  </Link>
                </div>
              )}
            </nav>

            {/* Mobile menu button */}
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="md:hidden p-2 rounded-md text-gray-700 hover:text-primary-600"
            >
              {isMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
            </button>
          </div>

          {/* Mobile menu */}
          {isMenuOpen && (
            <div className="md:hidden border-t border-gray-200 py-4">
              <form onSubmit={handleSearch} className="mb-4">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search products..."
                    className="input pl-10 pr-4 py-2 w-full"
                  />
                </div>
              </form>
              
              <nav className="space-y-2">
                <Link to="/products" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                  Products
                </Link>
                
                {isAuthenticated ? (
                  <>
                    <Link to="/orders" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Orders
                    </Link>
                    <Link to="/payments" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Payments
                    </Link>
                    <Link to="/chat" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Chat
                    </Link>
                    <Link to="/notifications" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Notifications
                    </Link>
                    <Link to="/cart" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Cart
                    </Link>
                    <Link to="/profile" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Profile
                    </Link>
                    <button
                      onClick={handleLogout}
                      className="block w-full text-left px-4 py-2 text-gray-700 hover:text-primary-600 font-medium"
                    >
                      Logout
                    </button>
                  </>
                ) : (
                  <>
                    <Link to="/login" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Login
                    </Link>
                    <Link to="/register" className="block px-4 py-2 text-gray-700 hover:text-primary-600 font-medium">
                      Sign Up
                    </Link>
                  </>
                )}
              </nav>
            </div>
          )}
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center text-gray-600">
            <p>&copy; 2024 DemoShop. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default Layout