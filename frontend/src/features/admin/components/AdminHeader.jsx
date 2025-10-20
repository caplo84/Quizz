import React, { useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useLocation, useNavigate } from 'react-router-dom';
import { logout } from '../authSlice';
import { 
  Bell as BellIcon, 
  UserCircle as UserCircleIcon,
  LogOut as ArrowRightOnRectangleIcon 
} from 'lucide-react';

const AdminHeader = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useSelector((state) => state.auth || {});
  const [showUserMenu, setShowUserMenu] = useState(false);

  const titleByPath = {
    '/admin': 'Admin Dashboard',
    '/admin/quizzes': 'Quiz Management',
    '/admin/topics': 'Topic Management',
    '/admin/sync': 'GitHub Sync',
    '/admin/bulk': 'Bulk Operations',
    '/admin/analytics': 'Analytics',
  };

  const currentTitle = titleByPath[location.pathname] || 'Admin Panel';

  const handleLogout = () => {
    dispatch(logout());
    navigate('/admin/login');
  };

  return (
    <header className="bg-white shadow-sm border-b border-gray-200">
      <div className="flex items-center justify-between px-6 py-4">
        {/* Page Title / Breadcrumb */}
        <div>
          <h2 className="text-xl font-semibold text-gray-800">
            {currentTitle}
          </h2>
        </div>

        {/* Right Side Actions */}
        <div className="flex items-center gap-4">
          {/* Notifications */}
          <button className="relative p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors">
            <BellIcon className="w-6 h-6" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>

          {/* User Menu */}
          <div className="relative">
            <button 
              onClick={() => setShowUserMenu(!showUserMenu)}
              className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <UserCircleIcon className="w-8 h-8" />
              <div className="text-left">
                <p className="text-sm font-medium">{user?.name || 'Admin User'}</p>
                <p className="text-xs text-gray-500">{user?.email || 'admin@example.com'}</p>
              </div>
            </button>

            {showUserMenu && (
              <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-10">
                <a
                  href="#"
                  className="flex items-center gap-2 px-4 py-2 text-sm hover:bg-gray-100"
                >
                  <UserCircleIcon className="w-4 h-4" />
                  Admin Profile
                </a>
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-2 w-full px-4 py-2 text-sm text-red-600 hover:bg-gray-100"
                >
                  <ArrowRightOnRectangleIcon className="w-4 h-4" />
                  Logout
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default AdminHeader;