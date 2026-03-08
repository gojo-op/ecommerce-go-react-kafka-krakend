import React from 'react'
import { useForm } from 'react-hook-form'
import { useMutation } from 'react-query'
import { useNavigate, Link } from 'react-router-dom'
import { toast } from 'sonner'
import { apiService } from '@/services/api'
import { useAuth } from '@/store/authStore'
import { useAuthStore } from '@/store/authStore'
import LoadingSpinner from '@/components/LoadingSpinner'

interface RegisterFormData {
  email: string
  username: string
  password: string
  first_name: string
  last_name: string
}

const Register: React.FC = () => {
  const navigate = useNavigate()
  const { login } = useAuth()
  const { register, handleSubmit, formState: { errors } } = useForm<RegisterFormData>()

  const registerMutation = useMutation(
    async (data: RegisterFormData) => {
      const res = await apiService.register(data)
      const payload = (res.data as any)?.data || (res.data as any)
      const { access_token, refresh_token, expires_at } = payload
      const expiresMs = Number(expires_at) * 1000
      useAuthStore.setState({ accessToken: access_token })
      const profileRes = await apiService.getProfile()
      const user = (profileRes.data as any)?.data || (profileRes.data as any)
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
      login(normalizedUser as any, access_token, refresh_token, String(expiresMs))
    },
    {
      onSuccess: () => {
        toast.success('Registration successful!')
        navigate('/')
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || 'Registration failed')
      },
    }
  )

  const onSubmit = (data: RegisterFormData) => {
    registerMutation.mutate(data)
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900 mb-2">
            Create your account
          </h2>
          <p className="text-gray-600">
            Already have an account?{' '}
            <Link to="/login" className="text-primary-600 hover:text-primary-500 font-medium">
              Sign in
            </Link>
          </p>
        </div>

        <div className="card">
          <div className="card-body">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">Email</label>
                <input
                  {...register('email', { required: 'Email is required' })}
                  type="email"
                  id="email"
                  className="input"
                  placeholder="you@example.com"
                />
                {errors.email && <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>}
              </div>

              <div>
                <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-2">Username</label>
                <input
                  {...register('username', { required: 'Username is required' })}
                  type="text"
                  id="username"
                  className="input"
                  placeholder="yourusername"
                />
                {errors.username && <p className="mt-1 text-sm text-red-600">{errors.username.message}</p>}
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">Password</label>
                <input
                  {...register('password', { required: 'Password is required', minLength: { value: 8, message: 'Minimum 8 characters' } })}
                  type="password"
                  id="password"
                  className="input"
                  placeholder="********"
                />
                {errors.password && <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>}
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label htmlFor="first_name" className="block text-sm font-medium text-gray-700 mb-2">First name</label>
                  <input
                    {...register('first_name', { required: 'First name is required' })}
                    type="text"
                    id="first_name"
                    className="input"
                    placeholder="John"
                  />
                  {errors.first_name && <p className="mt-1 text-sm text-red-600">{errors.first_name.message}</p>}
                </div>
                <div>
                  <label htmlFor="last_name" className="block text-sm font-medium text-gray-700 mb-2">Last name</label>
                  <input
                    {...register('last_name', { required: 'Last name is required' })}
                    type="text"
                    id="last_name"
                    className="input"
                    placeholder="Doe"
                  />
                  {errors.last_name && <p className="mt-1 text-sm text-red-600">{errors.last_name.message}</p>}
                </div>
              </div>

              <div>
                <button type="submit" disabled={registerMutation.isLoading} className="btn btn-primary w-full">
                  {registerMutation.isLoading ? (
                    <div className="flex items-center justify-center space-x-2">
                      <LoadingSpinner size="sm" />
                      <span>Creating account...</span>
                    </div>
                  ) : (
                    'Create account'
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Register