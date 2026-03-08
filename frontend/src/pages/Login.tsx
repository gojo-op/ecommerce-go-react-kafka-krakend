import React from 'react'
import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation } from 'react-query'
import { toast } from 'sonner'
import { apiService } from '@/services/api'
import { useAuth } from '@/store/authStore'
import { validateEmail } from '@/utils/helpers'
import LoadingSpinner from '@/components/LoadingSpinner'

interface LoginFormData {
  email: string
  password: string
}

const Login: React.FC = () => {
  const navigate = useNavigate()
  const { login } = useAuth()
  const { register, handleSubmit, formState: { errors } } = useForm<LoginFormData>()

  const loginMutation = useMutation(
    (data: LoginFormData) => apiService.login(data),
    {
      onSuccess: (response) => {
        const payload = (response.data as any)?.data || (response.data as any)
        const { user, access_token, refresh_token, expires_at } = payload
        const normalizedUser = {
          id: (user?.ID ?? user?.id ?? ''),
          email: (user?.email ?? ''),
          firstName: (user?.first_name ?? user?.firstName ?? ''),
          lastName: (user?.last_name ?? user?.lastName ?? ''),
          phone: (user?.phone ?? ''),
          avatar: (user?.avatar_url ?? user?.avatar ?? ''),
          isActive: (user?.is_active ?? user?.isActive ?? true),
          isEmailVerified: (user?.is_email_verified ?? user?.isEmailVerified ?? false),
          createdAt: (user?.created_at ?? user?.createdAt ?? ''),
          updatedAt: (user?.updated_at ?? user?.updatedAt ?? ''),
        }
        const expiresMs = Number(expires_at) * 1000
        login(normalizedUser as any, access_token, refresh_token, String(expiresMs))
        toast.success('Login successful!')
        navigate('/')
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || 'Login failed')
      },
    }
  )

  const onSubmit = (data: LoginFormData) => {
    loginMutation.mutate(data)
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900 mb-2">
            Sign in to your account
          </h2>
          <p className="text-gray-600">
            Or{' '}
            <Link to="/register" className="text-primary-600 hover:text-primary-500 font-medium">
              create a new account
            </Link>
          </p>
        </div>

        <div className="card">
          <div className="card-body">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                  Email address
                </label>
                <input
                  {...register('email', {
                    required: 'Email is required',
                    validate: (value) => validateEmail(value) || 'Please enter a valid email'
                  })}
                  type="email"
                  id="email"
                  className="input"
                  placeholder="Enter your email"
                />
                {errors.email && (
                  <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                  Password
                </label>
                <input
                  {...register('password', {
                    required: 'Password is required',
                    minLength: {
                      value: 8,
                      message: 'Password must be at least 8 characters'
                    }
                  })}
                  type="password"
                  id="password"
                  className="input"
                  placeholder="Enter your password"
                />
                {errors.password && (
                  <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
                )}
              </div>

              <div>
                <button
                  type="submit"
                  disabled={loginMutation.isLoading}
                  className="btn btn-primary w-full"
                >
                  {loginMutation.isLoading ? (
                    <div className="flex items-center justify-center space-x-2">
                      <LoadingSpinner size="sm" />
                      <span>Signing in...</span>
                    </div>
                  ) : (
                    'Sign in'
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>

        <div className="text-center text-sm text-gray-600">
          <p>
            By signing in, you agree to our{' '}
            <a href="#" className="text-primary-600 hover:text-primary-500">
              Terms of Service
            </a>{' '}
            and{' '}
            <a href="#" className="text-primary-600 hover:text-primary-500">
              Privacy Policy
            </a>
          </p>
        </div>
      </div>
    </div>
  )
}

export default Login