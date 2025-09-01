import { getQuiz } from "../services/questions";

export async function homeLoader() {
  console.log('🏠 Home loader starting...');
  
  try {
    // Add a small delay to ensure everything is loaded
    await new Promise(resolve => setTimeout(resolve, 100));
    
    const data = await getQuiz();
    console.log("📊 Home loader data received:", data);
    console.log("📊 Data type:", typeof data);
    console.log("📊 Is array:", Array.isArray(data));
    console.log("📊 Data length:", data ? data.length : 'no data');
    
    // Ensure we return an array
    if (!Array.isArray(data)) {
      console.warn('⚠️ Data is not an array, wrapping in array');
      return data ? [data] : [];
    }
    
    if (data.length === 0) {
      console.warn('⚠️ No quiz data received from API');
    }
    
    return data;
  } catch (error) {
    console.error("💥 Error loading quiz data:", error);
    console.error("💥 Error stack:", error.stack);
    
    // Return empty array instead of test data to avoid confusion
    console.log("🔄 Returning empty array due to error");
    return [];
  }
}
