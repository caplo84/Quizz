import { useState } from 'react';
import { useSelector } from 'react-redux';
import { getQuizImageUrl, getGitHubImageUrl } from '../utils/imageUtils';

function ImageDisplay({ imageUrl, altText, className = '', topic = 'android' }) {
  const { darkMode } = useSelector((state) => state.home);
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);
  const [isZoomed, setIsZoomed] = useState(false);

  if (!imageUrl || imageUrl.trim() === '') {
    return null;
  }

  // Try local image first, fallback to GitHub if not found
  const localImageUrl = getQuizImageUrl(imageUrl, topic);
  const githubImageUrl = getGitHubImageUrl(imageUrl, topic);
  const [currentImageUrl, setCurrentImageUrl] = useState(localImageUrl);



  const handleImageLoad = () => {
    setIsLoading(false);
    setHasError(false);
  };

  const handleImageError = () => {
    // If local image fails and we haven't tried GitHub yet, try GitHub URL
    if (currentImageUrl === localImageUrl && githubImageUrl) {
      setCurrentImageUrl(githubImageUrl);
      setIsLoading(true);
      setHasError(false);
      return;
    }
    
    // If both local and GitHub fail, show error
    setIsLoading(false);
    setHasError(true);
  };

  const toggleZoom = () => {
    setIsZoomed(!isZoomed);
  };

  // Fallback for broken images
  if (hasError) {
    return (
      <div className={`flex items-center justify-center p-8 rounded-lg border-2 border-dashed ${className} ${
        darkMode 
          ? 'bg-gray-800 border-gray-600 text-gray-400' 
          : 'bg-gray-100 border-gray-300 text-gray-500'
      }`}>
        <div className="text-center">
          <svg className="w-12 h-12 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <p className="text-sm font-medium">Image unavailable</p>
          {altText && <p className="text-xs mt-1">{altText}</p>}
        </div>
      </div>
    );
  }

  return (
    <>
      <div className={`relative rounded-lg overflow-hidden border ${className} ${
        darkMode 
          ? 'bg-gray-900 border-gray-700' 
          : 'bg-gray-50 border-gray-200'
      }`}>
        {/* Loading state */}
        {isLoading && (
          <div className={`flex items-center justify-center p-8 ${
            darkMode ? 'bg-gray-800 text-gray-400' : 'bg-gray-100 text-gray-500'
          }`}>
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-current mx-auto mb-2"></div>
              <p className="text-sm">Loading image...</p>
            </div>
          </div>
        )}
        
        {/* Image */}
        <img
          src={currentImageUrl}
          alt={altText || 'Quiz image'}
          className={`max-w-full h-auto cursor-pointer transition-all duration-200 hover:opacity-90 ${
            isLoading ? 'hidden' : 'block'
          }`}
          onLoad={handleImageLoad}
          onError={handleImageError}
          onClick={toggleZoom}
          loading="lazy"
        />
        
        {/* Zoom icon */}
        {!isLoading && !hasError && (
          <button
            onClick={toggleZoom}
            className={`absolute top-2 right-2 p-2 rounded-full transition-all duration-200 ${
              darkMode
                ? 'bg-gray-800 bg-opacity-80 hover:bg-opacity-100 text-gray-300 hover:text-white'
                : 'bg-white bg-opacity-80 hover:bg-opacity-100 text-gray-600 hover:text-gray-800'
            }`}
            title="Click to zoom"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
            </svg>
          </button>
        )}
      </div>

      {/* Zoom modal */}
      {isZoomed && (
        <div 
          className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-80 p-4"
          onClick={toggleZoom}
        >
          <div className="relative max-w-full max-h-full">
            <img
              src={currentImageUrl}
              alt={altText || 'Quiz image'}
              className="max-w-full max-h-full object-contain rounded-lg"
              onClick={(e) => e.stopPropagation()}
            />
            <button
              onClick={toggleZoom}
              className="absolute top-4 right-4 p-3 bg-black bg-opacity-60 text-white rounded-full hover:bg-opacity-80 transition-all duration-200"
              title="Close"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      )}
    </>
  );
}

export default ImageDisplay;