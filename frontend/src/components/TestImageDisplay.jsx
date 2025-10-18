
import ImageDisplay from '../ui/ImageDisplay';
import CodeBlock from '../ui/CodeBlock';

// Test component to verify images and code blocks are working
function TestImageDisplay() {
  const testImages = [
    'Q-141.png',
    'machine-learning_Q10.jpg', 
    'matlab_Q39.jpg',
    'solidworks_q3.jpg',
    'Git-WorkFlow.png'
  ];

  const testCode = `function bubbleSort(arr) {
    let n = arr.length;
    for (let i = 0; i < n - 1; i++) {
        for (let j = 0; j < n - i - 1; j++) {
            if (arr[j] > arr[j + 1]) {
                // Swap elements
                let temp = arr[j];
                arr[j] = arr[j + 1];
                arr[j + 1] = temp;
            }
        }
    }
    return arr;
}`;

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8 text-center">
        🧪 Test Images & Code Blocks
      </h1>
      
      {/* Test Code Block */}
      <div className="mb-12">
        <h2 className="text-xl font-semibold mb-4">📋 Code Block Test</h2>
        <CodeBlock 
          code={testCode}
          language="javascript"
          className="mb-4"
        />
      </div>

      {/* Test Images */}
      <div className="mb-12">
        <h2 className="text-xl font-semibold mb-4">🖼️ Image Display Test</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {testImages.map((imageName) => (
            <div key={imageName} className="border rounded-lg p-4">
              <h3 className="text-sm font-medium mb-2 text-gray-600">
                {imageName}
              </h3>
              <ImageDisplay 
                imageUrl={imageName}
                altText={`Test image ${imageName}`}
                className="w-full"
              />
            </div>
          ))}
        </div>
      </div>

      {/* URL Test */}
      <div className="bg-gray-100 p-4 rounded-lg">
        <h3 className="font-semibold mb-2">🔗 Expected Image URLs:</h3>
        <ul className="text-sm space-y-1">
          {testImages.map((img) => (
            <li key={img} className="font-mono">
              /images/quiz-images/{img}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default TestImageDisplay;