import { getQuiz } from "../services/questions";

export async function headerLoader() {
  const data = await getQuiz();
  return data;
}
