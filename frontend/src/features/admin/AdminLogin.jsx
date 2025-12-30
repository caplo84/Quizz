import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { setUser } from './authSlice';
import { setDarkMode } from '../home/homeSlice';
import { Moon, Sun } from 'lucide-react';

const AdminLogin = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const dispatch = useDispatch();
  const { darkMode } = useSelector((state) => state.home);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const reason = searchParams.get('reason');
  const reasonMessage =
    reason === 'unauthorized'
      ? 'Admin access is required for this page.'
      : reason === 'auth_required'
        ? 'Please sign in to continue.'
        : '';

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      // TODO: Replace with real backend auth endpoint
      if (formData.email === 'admin@quiz.com' && formData.password === 'admin123') {
        dispatch(
          setUser({
            id: 1,
            email: formData.email,
            name: 'Admin User',
            isAdmin: true,
          }),
        );
        navigate('/admin');
      } else {
        setError('Invalid email or password');
      }
    } catch (err) {
      setError(err.message || 'Login failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className={`min-h-screen flex items-center justify-center p-4 ${
        darkMode
          ? 'bg-gradient-to-br from-slate-950 via-slate-900 to-indigo-950'
          : 'bg-gradient-to-br from-blue-500 to-purple-600'
      }`}
    >
      <button
        type="button"
        onClick={() => dispatch(setDarkMode())}
        className={`absolute top-6 right-6 p-2 rounded-lg transition-colors ${
          darkMode ? 'bg-white/10 text-white hover:bg-white/20' : 'bg-white/20 text-white hover:bg-white/30'
        }`}
        aria-label="Toggle dark mode"
        title="Toggle dark mode"
      >
        {darkMode ? <Sun size={18} /> : <Moon size={18} />}
      </button>

      <div
        className={`rounded-xl shadow-2xl p-8 w-full max-w-md border ${
          darkMode
            ? 'bg-slate-900/90 border-slate-800 text-white'
            : 'bg-white border-transparent'
        }`}
      >
        {/* Logo/Header */}
        <div className="text-center mb-8">
          <h1 className={`text-3xl font-bold mb-2 ${darkMode ? 'text-white' : 'text-gray-800'}`}>Quiz Admin</h1>
          <p className={darkMode ? 'text-slate-300' : 'text-gray-600'}>Sign in to access the admin panel</p>
        </div>

        {reasonMessage && (
          <div className={`mb-4 px-4 py-3 rounded-lg text-sm border ${darkMode ? 'bg-amber-900/30 border-amber-700 text-amber-200' : 'bg-amber-50 border-amber-200 text-amber-700'}`}>
            {reasonMessage}
          </div>
        )}

        {/* Login Form */}
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Email Field */}
          <div>
            <label htmlFor="email" className={`block text-sm font-medium mb-2 ${darkMode ? 'text-slate-200' : 'text-gray-700'}`}>
              Email Address
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
              className={`w-full px-4 py-3 rounded-lg transition-colors focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${darkMode ? 'border border-slate-700 bg-slate-950 text-white placeholder:text-slate-400' : 'border border-gray-300'}`}
              placeholder="admin@quiz.com"
            />
          </div>

          {/* Password Field */}
          <div>
            <label htmlFor="password" className={`block text-sm font-medium mb-2 ${darkMode ? 'text-slate-200' : 'text-gray-700'}`}>
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              className={`w-full px-4 py-3 rounded-lg transition-colors focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${darkMode ? 'border border-slate-700 bg-slate-950 text-white placeholder:text-slate-400' : 'border border-gray-300'}`}
              placeholder="••••••••"
            />
          </div>

          {/* Error Message */}
          {error && (
            <div className={`px-4 py-3 rounded-lg text-sm border ${darkMode ? 'bg-red-900/30 border-red-700 text-red-200' : 'bg-red-50 border-red-200 text-red-600'}`}>
              {error}
            </div>
          )}

          {/* Submit Button */}
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>

        {/* Footer */}
        <div className="mt-6 text-center">
          <a href="/" className={`text-sm ${darkMode ? 'text-violet-300 hover:text-violet-200' : 'text-blue-600 hover:text-blue-700'}`}>
            ← Back to Main Site
          </a>
        </div>

        {/* Dev Note */}
        <div className={`mt-4 p-3 rounded text-xs border ${darkMode ? 'bg-slate-800 border-slate-700 text-slate-300' : 'bg-yellow-50 border-yellow-200 text-yellow-800'}`}>
          <strong>Dev Mode:</strong> Use admin@quiz.com / admin123
        </div>
      </div>
    </div>
  );
};

export default AdminLogin;