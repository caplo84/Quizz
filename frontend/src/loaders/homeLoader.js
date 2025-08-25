import { getQuiz } from "../services/questions";

export async function homeLoader() {
  try {
    const data = await getQuiz();
    console.log("Loader data:", data);
    return data;
  } catch (error) {
    console.error("Error loading quiz data:", error);
    return []; // Return empty array as fallback
  }
}
