import axios from 'axios';

export const HttpAgent = axios.create({
  baseURL: import.meta.env.VITE_ARCHIVE1_API_PANEL_BASE_URL,
  timeout: 30000,
});
