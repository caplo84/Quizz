import { RouterProvider, createBrowserRouter } from "react-router-dom";
import AppLayout from "./ui/AppLayout";
import Home from "./features/home/Home";
import { homeLoader } from "./loaders/homeLoader";
import Quiz from "./features/quiz/Quiz";
import Error from "./ui/Error";
import FinishedScreen from "./ui/FinishedScreen";
import ConnectionStatus from "./ui/ConnectionStatus";
import TestImageDisplay from './components/TestImageDisplay';

function App() {
  const router = createBrowserRouter([
    {
      element: <AppLayout />,

      children: [
        {
          path: "/",
          element: <Home />,
          loader: homeLoader,
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
