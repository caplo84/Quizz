import { useLoaderData } from "react-router-dom";
import HomeItem from "./HomeItem";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { setQuizzes } from "./homeSlice";

function Home() {
  const { darkMode } = useSelector((state) => state.home);
  const data = useLoaderData();
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(setQuizzes(data));
  }, [dispatch, data]);

  // Add debugging to see what data looks like
  console.log("Data received in Home:", data);
  console.log("Is data an array?", Array.isArray(data));

  return (
    <div className="desktop:grid-cols-1 desktop:gap-24 mobile:gap-16 grid grid-cols-2">
      <div>
        <h1
          className={`mobile:text-[4rem] mobile:pb-8 pb-20 text-[6.4rem] font-light leading-[100%] transition-all duration-300 ${
            darkMode ? "text-white" : "text-dark-navy"
          }`}
        >
          Welcome to the <br />
          <span className="font-medium">Frontend Quiz!</span>
        </h1>
        <p
          className={`mobile:text-[1.8rem] text-[2rem] italic leading-[150%] transition-all duration-300 ${
            darkMode ? "text-light-bluish" : "text-grey-navy"
          }`}
        >
          Pick a subject to get started.
        </p>
      </div>
      <div>
        <ul className="mobile:space-y-5 space-y-8">
          {data && Array.isArray(data) && data.length > 0 ? (
            data.map((item) => (
              <HomeItem key={item.title} img={item.icon} text={item.title} slug={item.slug} />
            ))
          ) : (
            <li className="text-center py-12">
              <div className={`${darkMode ? "text-white" : "text-dark-navy"} space-y-4`}>
                <div className="mb-6">
                  <svg className="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                  </svg>
                </div>
                <h3 className="text-xl font-semibold mb-2">No Quizzes Available</h3>
                <p className="text-gray-500 max-w-md mx-auto">
                  Quiz topics are being prepared. Please check back soon for exciting challenges!
                </p>
              </div>
            </li>
          )}
        </ul>
      </div>
    </div>
  );
}

export default Home;
