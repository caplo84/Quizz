import { Outlet, Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import AdminSidebar from './components/AdminSidebar';
import AdminHeader from './components/AdminHeader';

const AdminLayout = () => {
  const { user } = useSelector((state) => state.auth);

  // Redirect to login if not authenticated or not admin
  if (!user) {
    return <Navigate to="/admin/login?reason=auth_required" replace />;
  }

  if (!user.isAdmin) {
    return <Navigate to="/admin/login?reason=unauthorized" replace />;
  }

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <AdminSidebar />

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <AdminHeader />

        {/* Content Area */}
        <main className="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default AdminLayout;
