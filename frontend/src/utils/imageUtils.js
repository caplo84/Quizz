/**
 * Transforms re  // Handle di  // Handle different image URL formats
  
  // Case 1: Already standardized filename (topic_filename format)
  if (imageUrl.includes('_') && !imageUrl.includes('/')) {
    return `/images/quiz-images/${imageUrl}`;
  }
  
  // Case 2: Relative path format (image/04.jpeg, css/img/Q-141.png)
  if (imageUrl.includes('/')) {
    const filename = imageUrl.split('/').pop();
    return `http://localhost:8080/static/quiz-images/${filename}`;
  }
  
  // Case 3: Just filename, use as is
  return `http://localhost:8080/static/quiz-images/${imageUrl}`;ormats
  
  // Case 1: Already standardized filename (topic_filename format)
  if (imageUrl.includes('_') && !imageUrl.includes('/')) {
    return `/images/quiz-images/${imageUrl}`;
  }
  
  // Case 2: Relative path format (image/04.jpeg, css/img/Q-141.png)
  if (imageUrl.includes('/')) {
    const filename = imageUrl.split('/').pop();
    return `/images/quiz-images/${filename}`;
  }
  
  // Case 3: Just filename, use as is
  return `/images/quiz-images/${imageUrl}`;to local quiz image URLs
 * @param {string} imageUrl - The relative image URL from GitHub (e.g., "image/04.jpeg")
 * @param {string} topic - The quiz topic (e.g., "android", "javascript") 
 * @returns {string} - The local image URL
 */
export function getQuizImageUrl(imageUrl, topic = '') {
  if (!imageUrl || imageUrl.trim() === '') {
    return null;
  }

  // If it's already a full URL, return as is
  if (imageUrl.startsWith('http://') || imageUrl.startsWith('https://')) {
    return imageUrl;
  }

  // If it's already a local path starting with /, return as is
  if (imageUrl.startsWith('/')) {
    return imageUrl;
  }

  // Handle different image URL formats
  console.log('🖼️ Processing image URL:', imageUrl, 'for topic:', topic);
  
  // Case 1: Already standardized filename (topic_filename format)
  if (imageUrl.includes('_') && !imageUrl.includes('/')) {
    console.log('✅ Standardized format detected:', imageUrl);
    return `/images/quiz-images/${imageUrl}`;
  }
  
  // Case 2: Relative path format (image/04.jpeg, css/img/Q-141.png)
  if (imageUrl.includes('/')) {
    const filename = imageUrl.split('/').pop();
    console.log('� Using original filename:', filename);
    return `http://localhost:8080/static/quiz-images/${filename}`;
  }
  
  // Case 3: Just filename, use as is
  console.log('📝 Using filename as is:', imageUrl);
  return `http://localhost:8080/static/quiz-images/${imageUrl}`;
}

/**
 * Gets the full GitHub raw URL for downloading images
 * @param {string} imageUrl - The relative image URL (e.g., "image/04.jpeg")
 * @param {string} topic - The quiz topic (e.g., "android")
 * @returns {string} - The full GitHub raw URL
 */
export function getGitHubImageUrl(imageUrl, topic = 'android') {
  if (!imageUrl || imageUrl.trim() === '') {
    return null;
  }

  const baseUrl = 'https://raw.githubusercontent.com/Ebazhanov/linkedin-skill-assessments-quizzes/main';
  return `${baseUrl}/${topic}/${imageUrl}`;
}