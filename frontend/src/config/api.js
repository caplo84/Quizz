const config = {
  BASE_URL: import.meta.env.VITE_API_URL || import.meta.env.REACT_APP_API_URL || 'http://localhost:8080',
  API_VERSION: 'v1',
  TIMEOUT: 10000,
  
  get API_BASE() {
    return `${this.BASE_URL}/api/${this.API_VERSION}`;
  }
};

export default config;