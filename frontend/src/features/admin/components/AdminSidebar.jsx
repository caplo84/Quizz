import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { 
  Home as HomeIcon, 
  FileText as DocumentTextIcon, 
  Folder as FolderIcon, 
  RefreshCw as ArrowPathIcon,
  Layers as LayersIcon,
  BarChart3 as ChartBarIcon,
  ArrowLeft
} from 'lucide-react';

const AdminSidebar = () => {
  const location = useLocation();

  const navItems = [
    { path: '/admin', icon: HomeIcon, label: 'Dashboard', exact: true },
    { path: '/admin/quizzes', icon: DocumentTextIcon, label: 'Quizzes' },
    { path: '/admin/topics', icon: FolderIcon, label: 'Topics' },
    { path: '/admin/sync', icon: ArrowPathIcon, label: 'GitHub Sync' },
    { path: '/admin/bulk', icon: LayersIcon, label: 'Bulk Operations' },
    { path: '/admin/analytics', icon: ChartBarIcon, label: 'Analytics' },
  ];

  const isActive = (item) => {
    if (item.exact) {
      return location.pathname === item.path;
    }
    return location.pathname.startsWith(item.path);
  };

  return (
    <div className="w-64 bg-gray-800 text-white flex flex-col">
      {/* Logo/Brand */}
      <div className="p-4 border-b border-gray-700">
        <h1 className="text-2xl font-bold">Quiz Admin</h1>
        <p className="text-sm text-gray-400 mt-1">Management Panel</p>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto p-4">
        <ul className="space-y-2">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item);
            
            return (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                    active
                      ? 'bg-blue-600 text-white'
                      : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                  }`}
                >
                  <Icon className="w-5 h-5" />
                  <span className="font-medium">{item.label}</span>
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Footer */}
      <div className="p-4 border-t border-gray-700">
        <Link
          to="/"
          className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
        >
          <ArrowLeft className="w-5 h-5" />
          <span>Back to Main Site</span>
        </Link>
      </div>
    </div>
  );
};

export default AdminSidebar;