import axios from 'axios';

export const HttpAgent = axios.create({
  baseURL: import.meta.env.VITE_ARCHIVO_API_PANEL_BASE_URL,
  timeout: 30000,
  withCredentials: true,
});
