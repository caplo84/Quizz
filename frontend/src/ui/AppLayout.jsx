import { Outlet, useLocation } from "react-router-dom";
import Header from "./Header";
import PageDesign from "./PageDesign";
import { useSelector } from "react-redux";

function AppLayout() {
  const { darkMode } = useSelector((state) => state.home);
  const location = useLocation();
  const isHomePage = location.pathname === '/';

  return (
    <div
      className={`relative min-h-screen font-rubik transition-all duration-300 ${
        isHomePage 
          ? '' // No background styling for home page (let Home component handle it)
          : `laptop:px-24 desktop:py-16 mobile:py-8 mobile:px-8 grid grid-rows-[auto_1fr] overflow-hidden px-56 py-32 ${
              darkMode ? "bg-dark-navy" : "bg-light-grey"
            }`
      }`}
    >
      {!isHomePage && <Header />}
      <main className={isHomePage ? '' : 'z-[2]'}>
        <Outlet />
      </main>
      {!isHomePage && <PageDesign />}
    </div>
  );
}

export default AppLayout;
