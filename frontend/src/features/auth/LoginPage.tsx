import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { loginSchema, type LoginInput } from '@/schemas/auth.schema'
import { useAuth } from './AuthContext'

export function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginInput>({ resolver: zodResolver(loginSchema) })

  const onSubmit = async (data: LoginInput) => {
    try {
      await login(data.phone, data.password)
      navigate('/', { replace: true })
    } catch {
      toast.error('رقم الهاتف أو كلمة المرور غير صحيحة')
    }
  }

  return (
    <div className="min-h-screen flex" dir="rtl">
      {/* Left panel — branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-primary items-center justify-center p-12">
        <div className="text-center text-white">
          <div className="w-24 h-24 rounded-full bg-white/20 flex items-center justify-center mx-auto mb-6">
            <span className="text-5xl font-bold">ك</span>
          </div>
          <h1 className="text-4xl font-bold mb-3">كشفي</h1>
          <p className="text-primary-100 text-lg">نظام إدارة الفوج الكشفي</p>
        </div>
      </div>

      {/* Right panel — form */}
      <div className="flex-1 flex items-center justify-center px-6 py-12 bg-gray-50 dark:bg-slate-950">
        <div className="w-full max-w-md">
          <div className="card p-8">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
              تسجيل الدخول
            </h2>
            <p className="text-gray-500 dark:text-slate-400 mb-8 text-sm">
              أدخل بيانات حسابك للمتابعة
            </p>

            <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-5">
              <div>
                <label className="label">رقم الهاتف</label>
                <input
                  {...register('phone')}
                  type="tel"
                  placeholder="03XXXXXXX"
                  className="input"
                  dir="ltr"
                />
                {errors.phone && (
                  <p className="mt-1 text-xs text-red-600">{errors.phone.message}</p>
                )}
              </div>

              <div>
                <label className="label">كلمة المرور</label>
                <input
                  {...register('password')}
                  type="password"
                  placeholder="••••••••"
                  className="input"
                  dir="ltr"
                />
                {errors.password && (
                  <p className="mt-1 text-xs text-red-600">{errors.password.message}</p>
                )}
              </div>

              <button type="submit" disabled={isSubmitting} className="btn-primary w-full py-2.5 mt-2">
                {isSubmitting ? 'جاري الدخول...' : 'دخول'}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
