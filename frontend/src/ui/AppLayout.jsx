import { Outlet, useLocation } from "react-router-dom";
import Header from "./Header";
import PageDesign from "./PageDesign";
import { useSelector } from "react-redux";

function AppLayout() {
  const { darkMode } = useSelector((state) => state.home);
  const location = useLocation();
  const isHomePage = location.pathname === '/';
  const staticRoutes = new Set(['/legacy', '/login', '/leaderboard', '/dashboard', '/finished', '/test-images']);
  const isQuizPage = !isHomePage && !staticRoutes.has(location.pathname);
  const isFullBleedPage = isHomePage || isQuizPage;

  return (
    <div
      className={`relative min-h-screen font-rubik transition-all duration-300 ${
        isFullBleedPage
          ? '' // Full-bleed pages handle their own layout (home + quiz)
          : `laptop:px-24 desktop:py-16 mobile:py-8 mobile:px-8 grid grid-rows-[auto_1fr] overflow-hidden px-56 py-32 ${
              darkMode ? "bg-dark-navy" : "bg-light-grey"
            }`
      }`}
    >
      {!isFullBleedPage && <Header />}
      <main className={isFullBleedPage ? '' : 'z-[2]'}>
        <Outlet />
      </main>
      {!isFullBleedPage && <PageDesign />}
    </div>
  );
}

export default AppLayout;
