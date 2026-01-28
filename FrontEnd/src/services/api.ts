import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api', // Vite proxy will handle this
  headers: {
    'Content-Type': 'application/json',
  },
});

export default {
  
  getStocks(params: any) {
    return apiClient.get('/stocks', { params });
  },

  getRecommendations(params: any) {
      return apiClient.get('/recommend', { params });
  }
};
