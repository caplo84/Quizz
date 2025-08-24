import apiClient from './apiClient.js';

export const healthApi = {
  check: () => apiClient.get('/health'),
};