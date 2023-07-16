import { AuthProvider as IAuthProvider } from 'ra-core';
import { HttpAgent } from './utils/http-agent';
import { ArchiveResponse } from './utils/types';

export const TOKEN_KEY = 'tkn';

export const AuthProvider: IAuthProvider = {
  login: async (params: { email: string; password: string }) => {
    try {
      const res = await HttpAgent.post<
        ArchiveResponse<{ access_token: string }>
      >('/auth/login', {
        email: params.email,
        password: params.password,
      });
      const token = res.data.data.access_token;
      localStorage.setItem(TOKEN_KEY, token);
    } catch (err) {
      console.log('login error', err);
      throw err;
    }
  },
  checkAuth: async (_) => {
    const token = localStorage.getItem(TOKEN_KEY) || '';
    try {
      await HttpAgent.get<
        ArchiveResponse<{
          user: {
            id: number;
            username: string;
            email: string;
            is_admin: boolean;
            created_at: string;
            updated_at: string;
          };
        }>
      >('/auth/me', {
        headers: { authorization: `Bearer ${token}` },
      });
    } catch (err) {
      console.log('check auth error: ', err);
      throw err;
    }
  },
  checkError: async (error) => {
    const status = error.status;
    if (status === 401 || status === 403) {
      localStorage.removeItem(TOKEN_KEY);
      return Promise.reject();
    }
    // other error code (404, 500, etc): no need to log out
    return Promise.resolve();
  },
  logout: async () => {
    localStorage.removeItem(TOKEN_KEY);
    console.log('user logged out');
  },
  getIdentity: async () => {
    const token = localStorage.getItem(TOKEN_KEY) || '';
    try {
      const res = await HttpAgent.get<
        ArchiveResponse<{
          user: {
            id: number;
            username: string;
            email: string;
            is_admin: boolean;
            created_at: string;
            updated_at: string;
          };
        }>
      >('/auth/me', {
        headers: { authorization: `Bearer ${token}` },
      });
      console.log('user id', res.data.data.user.id);
      return {
        id: 0,
        fullName: res.data.data.user.username,
        // avatar: ''
      };
    } catch (err) {
      console.log('check auth error: ', err);
      throw err;
    }
  },
  getPermissions: async () => {},
};
