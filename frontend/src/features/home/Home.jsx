import { useLoaderData } from "react-router-dom";
import HomeItem from "./HomeItem";
import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { setQuizzes, setDarkMode } from "./homeSlice";

function Home() {
  const { darkMode } = useSelector((state) => state.home);
  const data = useLoaderData();
  const dispatch = useDispatch();
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    dispatch(setQuizzes(data));
  }, [dispatch, data]);

  // Filter data based on search term
  const filteredData = data && Array.isArray(data) 
    ? data.filter(item => 
        item.title.toLowerCase().includes(searchTerm.toLowerCase())
      )
    : [];

  return (
    <div className={`min-h-screen relative overflow-hidden transition-all duration-500 ${
      darkMode 
        ? 'bg-gradient-to-br from-gray-900 via-purple-900 to-indigo-900' 
        : 'bg-gradient-to-br from-slate-50 to-gray-100'
    }`}>
      {/* Animated Background Elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className={`absolute top-1/4 left-1/4 w-64 h-64 rounded-full blur-3xl animate-pulse ${
          darkMode ? 'bg-purple-500/10' : 'bg-blue-200/30'
        }`}></div>
        <div className={`absolute top-3/4 right-1/4 w-48 h-48 rounded-full blur-2xl animate-bounce ${
          darkMode ? 'bg-indigo-500/10' : 'bg-purple-200/20'
        }`} style={{animationDelay: '1s', animationDuration: '3s'}}></div>
        <div className={`absolute top-1/2 left-1/2 w-32 h-32 rounded-full blur-xl animate-ping ${
          darkMode ? 'bg-blue-500/10' : 'bg-indigo-200/25'
        }`} style={{animationDelay: '2s'}}></div>
        <div className={`absolute top-10 right-10 w-20 h-20 rounded-full blur-lg animate-bounce ${
          darkMode ? 'bg-purple-400/20' : 'bg-pink-200/30'
        }`} style={{animationDelay: '0.5s', animationDuration: '4s'}}></div>
      </div>
      
      <div className="container mx-auto px-8 py-12 relative z-10">
        {/* Dark Mode Toggle - Top Right */}
        <div className="absolute top-6 right-6 z-20">
          <div className={`flex items-center gap-4 backdrop-blur-lg rounded-full px-4 py-2 border shadow-lg ${
            darkMode 
              ? 'bg-gray-800/50 border-gray-600' 
              : 'bg-white/80 border-gray-200'
          }`}>
            {darkMode ? (
              <img src="/icon-sun-light.svg" alt="sun icon" className="w-6 h-6" />
            ) : (
              <img src="/icon-sun-dark.svg" alt="sun icon" className="w-6 h-6" />
            )}
            <label
              htmlFor="toggle"
              className={`relative block h-7 w-12 cursor-pointer rounded-full transition-colors duration-300 before:absolute before:left-1 before:top-1 before:h-5 before:w-5 before:rounded-full before:bg-white before:transition-all before:duration-300 ${
                darkMode 
                  ? 'bg-purple-500 hover:bg-purple-600 before:translate-x-5' 
                  : 'bg-slate-300 hover:bg-slate-400 before:translate-x-0'
              }`}
            >
              <input
                type="checkbox"
                id="toggle"
                className="hidden"
                checked={darkMode}
                onChange={() => dispatch(setDarkMode())}
              />
            </label>
            {darkMode ? (
              <img src="/icon-moon-light.svg" alt="moon icon" className="w-6 h-6" />
            ) : (
              <img src="/icon-moon-dark.svg" alt="moon icon" className="w-6 h-6" />
            )}
          </div>
        </div>

        {/* Header Section */}
        <div className="text-center mb-16">
          <h1 className={`text-7xl font-bold mb-6 tracking-tight drop-shadow-2xl transform hover:scale-105 transition-transform duration-300 ${
            darkMode ? 'text-white' : 'text-slate-800'
          }`}>
            🚀 Programming Quiz ⚡
          </h1>
          <p className={`text-2xl mb-12 font-medium drop-shadow-lg animate-pulse ${
            darkMode ? 'text-gray-200' : 'text-slate-600'
          }`}>
            💻 Test your coding knowledge with fun quizzes! 🎯
          </p>
          
          {/* Search Bar */}
          <div className="relative inline-block mb-4">
            <div className="relative group">
              <input
                type="text"
                placeholder="🔍 Search for amazing quizzes..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className={`px-8 py-5 rounded-2xl backdrop-blur-sm border-none outline-none focus:ring-4 focus:scale-105 w-96 text-lg shadow-2xl transform transition-all duration-300 hover:shadow-3xl group-hover:bg-white ${
                  darkMode 
                    ? 'bg-gray-800/50 border border-gray-600 text-white placeholder-gray-300 focus:ring-purple-400/50' 
                    : 'bg-white/90 border border-gray-200 text-slate-700 placeholder-slate-400 focus:ring-blue-400/50'
                }`}
                style={{
                  boxShadow: '0 20px 40px rgba(0,0,0,0.1), 0 0 0 1px rgba(255,255,255,0.2)'
                }}
              />
              <div className="absolute right-6 top-1/2 transform -translate-y-1/2 group-focus-within:scale-110 transition-transform duration-300">
                <svg className={`w-6 h-6 group-focus-within:text-purple-500 transition-colors duration-300 ${
                  darkMode ? 'text-gray-300' : 'text-slate-400'
                }`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        {/* Topics Grid */}
        {filteredData && filteredData.length > 0 ? (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-4 gap-6 max-w-6xl mx-auto">
            {filteredData.map((item) => (
              <HomeItem 
                key={item.title} 
                img={item.icon} 
                text={item.title} 
                slug={item.slug}
                source={item.source}
                quizCount={item.quizCount}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-20">
            <div className="space-y-6">
              <div className="mb-8">
                <svg className={`w-20 h-20 mx-auto mb-6 ${
                  darkMode ? 'text-gray-400' : 'text-slate-400'
                }`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
              <h3 className={`text-2xl font-semibold mb-4 ${
                darkMode ? 'text-gray-200' : 'text-slate-700'
              }`}>
                {searchTerm ? 'No matching quizzes found' : 'No Quizzes Available'}
              </h3>
              <p className={`max-w-lg mx-auto text-lg ${
                darkMode ? 'text-gray-300' : 'text-slate-500'
              }`}>
                {searchTerm 
                  ? `No quizzes match "${searchTerm}". Try a different search term.`
                  : 'Quiz topics are being prepared. Please check back soon for exciting challenges!'
                }
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default Home;
