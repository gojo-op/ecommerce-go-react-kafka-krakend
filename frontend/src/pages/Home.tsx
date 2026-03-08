import React from 'react'
import { Link } from 'react-router-dom'
import { Package, Shield, Truck, Star } from 'lucide-react'

const Home: React.FC = () => {
  const features = [
    {
      icon: Package,
      title: 'Quality Products',
      description: 'Discover amazing products from trusted sellers',
    },
    {
      icon: Shield,
      title: 'Secure Shopping',
      description: 'Your payments and data are always protected',
    },
    {
      icon: Truck,
      title: 'Fast Delivery',
      description: 'Get your orders delivered quickly to your doorstep',
    },
    {
      icon: Star,
      title: 'Excellent Service',
      description: '24/7 customer support to help with your needs',
    },
  ]

  return (
    <div className="space-y-16">
      {/* Hero Section */}
      <section className="text-center py-20 bg-gradient-to-br from-primary-50 to-primary-100 rounded-2xl">
        <div className="max-w-4xl mx-auto px-4">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
            Welcome to{' '}
            <span className="text-primary-600">DemoShop</span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Discover amazing products, enjoy secure shopping, and experience fast delivery. 
            Your perfect shopping destination awaits.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link to="/products" className="btn btn-primary text-lg px-8 py-3">
              Shop Now
            </Link>
            <Link to="/register" className="btn btn-outline text-lg px-8 py-3">
              Create Account
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section>
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Why Choose DemoShop?
          </h2>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            We provide the best shopping experience with quality products, secure payments, and excellent customer service.
          </p>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {features.map((feature, index) => (
            <div key={index} className="card text-center">
              <div className="card-body">
                <div className="mx-auto mb-4 p-3 bg-primary-100 rounded-full w-16 h-16 flex items-center justify-center">
                  <feature.icon className="h-8 w-8 text-primary-600" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  {feature.title}
                </h3>
                <p className="text-gray-600">
                  {feature.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-primary-600 rounded-2xl py-16">
        <div className="text-center text-white">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Ready to Start Shopping?
          </h2>
          <p className="text-xl mb-8 opacity-90 max-w-2xl mx-auto">
            Join thousands of satisfied customers and discover amazing products today.
          </p>
          <Link to="/products" className="btn bg-white text-primary-600 hover:bg-gray-50 text-lg px-8 py-3">
            Browse Products
          </Link>
        </div>
      </section>
    </div>
  )
}

export default Home