import { getQuiz } from "../services/api";

export async function homeLoader() {
  try {
    // Add a small delay to ensure everything is loaded
    await new Promise(resolve => setTimeout(resolve, 100));
    
    const data = await getQuiz();
    
    // Ensure we return an array
    if (!Array.isArray(data)) {
      return data ? [data] : [];
    }
    
    if (data.length === 0) {
      console.warn('⚠️ No quiz data received from API');
    }
    
    return data;
  } catch (error) {
    console.error("Error loading quiz data:", error);
    
    // Return empty array instead of test data to avoid confusion
    return [];
  }
}
