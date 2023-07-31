import { AuthProvider as IAuthProvider } from 'ra-core';
import { HttpAgent } from './utils/http-agent';
import { ArchiveResponse } from './utils/types';

export const AuthProvider: IAuthProvider = {
  login: async (params: { email: string; password: string }) => {
    try {
      await HttpAgent.post<ArchiveResponse<{ access_token: string }>>(
        '/auth/login',
        {
          email: params.email,
          password: params.password,
        },
      );
    } catch (err) {
      console.log('login error', err);
      throw err;
    }
  },
  checkAuth: async (_) => {
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
      >('/auth/me');
    } catch (err) {
      console.log('check auth error: ', err);
      throw err;
    }
  },
  checkError: async (error) => {
    const status = error.status;
    if (status === 401 || status === 403) {
      return Promise.reject();
    }
    return Promise.resolve();
  },
  logout: async () => {
    // todo: call logout route
    console.log('user logged out');
  },
  getIdentity: async () => {
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
      >('/auth/me');
      console.log('user id', res.data.data.user.id);
      return {
        id: 0,
        fullName: res.data.data.user.username,
      };
    } catch (err) {
      console.log('check auth error: ', err);
      throw err;
    }
  },
  getPermissions: async () => {},
};
