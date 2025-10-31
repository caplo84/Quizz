import { RouterProvider, createBrowserRouter } from "react-router-dom";
import AppLayout from "./ui/AppLayout";
import Home from "./features/home/Home";
import HomeModern from "./features/home/HomeModern";
import { homeLoader } from "./loaders/homeLoader";
import Quiz from "./features/quiz/Quiz";
import Error from "./ui/Error";
import FinishedScreen from "./ui/FinishedScreen";
import ConnectionStatus from "./ui/ConnectionStatus";
import TestImageDisplay from './components/TestImageDisplay';
import UserLogin from "./features/user/UserLogin";
import Leaderboard from "./features/user/Leaderboard";
import Dashboard from "./features/user/Dashboard";
import AdminLayout from "./features/admin/AdminLayout";
import AdminLogin from "./features/admin/AdminLogin";
import AdminDashboard from "./features/admin/AdminDashboard";
import QuizManagement from "./features/admin/QuizManagement";
import QuizForm from "./features/admin/QuizForm";
import TopicManagement from "./features/admin/TopicManagement";
import SyncManagement from "./features/admin/SyncManagement";
import BulkOperations from "./features/admin/BulkOperations";
import AISettings from "./features/admin/AISettings";

function App() {
  const router = createBrowserRouter([
    {
      path: "/admin/login",
      element: <AdminLogin />,
    },
    {
      path: "/admin",
      element: <AdminLayout />,
      children: [
        {
          index: true,
          element: <AdminDashboard />,
        },
        {
          path: "quizzes",
          element: <QuizManagement />,
        },
        {
          path: "quizzes/new",
          element: <QuizForm />,
        },
        {
          path: "quizzes/:id/edit",
          element: <QuizForm />,
        },
        {
          path: "topics",
          element: <TopicManagement />,
        },
        {
          path: "topics/new",
          element: <TopicManagement />,
        },
        {
          path: "sync",
          element: <SyncManagement />,
        },
        {
          path: "bulk",
          element: <BulkOperations />,
        },
        {
          path: "analytics",
          element: <AdminDashboard />,
        },
        {
          path: "settings",
          element: <AISettings />,
        },
      ],
    },
    {
      element: <AppLayout />,

      children: [
        {
          path: "/",
          element: <HomeModern />,
          loader: homeLoader,
        },
        {
          path: "/legacy",
          element: <Home />,
          loader: homeLoader,
        },
        {
          path: "/login",
          element: <UserLogin />,
        },
        {
          path: "/leaderboard",
          element: <Leaderboard />,
        },
        {
          path: "/dashboard",
          element: <Dashboard />,
        },
        {
          path: "/:type",
          element: <Quiz />,
          errorElement: <Error />,
        },
        {
          path: "/finished",
          element: <FinishedScreen />,
        },
        {
          path: "/test-images",
          element: <TestImageDisplay />,
        },
      ],
    },
  ]);

  return (
    <>
      <RouterProvider router={router} />
      <ConnectionStatus />
    </>
  );
}

export default App;
