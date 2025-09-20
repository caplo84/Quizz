// Real quiz data extracted from backend database migration
// This is a temporary solution until backend API includes choices in the response

export const QUIZ_DATA = {
  'javascript-basics': {
    1: {
      choices: [
        'Declares a block-scoped variable',
        'Declares a function-scoped variable', 
        'Declares a constant variable',
        'Declares a global variable'
      ],
      correctAnswer: 'Declares a block-scoped variable'
    },
    2: {
      choices: [
        'string',
        'boolean',
        'float',
        'object'
      ],
      correctAnswer: 'float'
    },
    3: {
      choices: [
        'Stops the execution of a function and returns a value',
        'Restarts the function',
        'Throws an error', 
        'Prints a value to the console'
      ],
      correctAnswer: 'Stops the execution of a function and returns a value'
    },
    4: {
      choices: [
        'append()',
        'push()',
        'add()', 
        'insert()'
      ],
      correctAnswer: 'push()'
    },
    5: {
      choices: [
        'string',
        'object',
        'null',
        'undefined'
      ],
      correctAnswer: 'object'
    },
    6: {
      choices: [
        '==',
        '!=',
        '===',
        '='
      ],
      correctAnswer: '==='
    },
    7: {
      choices: [
        'call myFunction',
        'myFunction()',
        'execute myFunction',
        'run myFunction()'
      ],
      correctAnswer: 'myFunction()'
    },
    8: {
      choices: [
        'pop()',
        'remove()',
        'delete()',
        'shift()'
      ],
      correctAnswer: 'pop()'
    },
    9: {
      choices: [
        'var colors = ["red", "green", "blue"]',
        'var colors = (1:"red", 2:"green", 3:"blue")',
        'var colors = "red", "green", "blue"',
        'var colors = 1 = ("red"), 2 = ("green"), 3 = ("blue")'
      ],
      correctAnswer: 'var colors = ["red", "green", "blue"]'
    },
    10: {
      choices: [
        'length()',
        'length',
        'size()',
        'count()'
      ],
      correctAnswer: 'length'
    }
  },
  'html-basics': {
    // Add HTML quiz data here if needed
  },
  'css-basics': {
    // Add CSS quiz data here if needed
  },
  'accessibility-basics': {
    // Add Accessibility quiz data here if needed
  }
};
