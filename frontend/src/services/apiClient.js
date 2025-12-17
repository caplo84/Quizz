import config from '../config/api.js';

class ApiClient {
  constructor() {
    this.baseURL = config.API_BASE;
    this.fallbackBaseURL = '/api/v1';
  }

  async request(endpoint, options = {}) {
    const requestOptions = {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      };

    const isAbsoluteApiPath = endpoint.startsWith('/api/');
    const candidates = isAbsoluteApiPath
      ? [`${config.BASE_URL}${endpoint}`, endpoint]
      : [`${this.baseURL}${endpoint}`, `${this.fallbackBaseURL}${endpoint}`];
    let lastError = null;

    for (const url of candidates) {

      try {
        const response = await fetch(url, requestOptions);

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json().catch(() => ({}));
        return data;
      } catch (error) {
        lastError = error;

        // Retry with fallback for network-level failures only
        if (error instanceof TypeError) {
          continue;
        }

        throw error;
      }
    }

    console.error('💥 Request failed:', lastError);
    throw lastError;
  }

  // HTTP Methods
  get(endpoint, options = {}) {
    return this.request(endpoint, { ...options, method: 'GET' });
  }

  post(endpoint, data, options = {}) {
    return this.request(endpoint, {
      ...options,
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  put(endpoint, data, options = {}) {
    return this.request(endpoint, {
      ...options,
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  delete(endpoint, options = {}) {
    return this.request(endpoint, { ...options, method: 'DELETE' });
  }
}

export default new ApiClient();