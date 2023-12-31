import { AuthProvider as IAuthProvider } from 'ra-core';
import { HttpAgent } from './utils/http-agent';
import { ArchiveResponse } from './utils/types';
import { toast } from 'react-toastify';

export enum PERMISSIONS {
  ADMIN = 'admin',
  USER = 'user',
}

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

      // check that user should redirect to pre-auth change password or not
      const res = await HttpAgent.get<
        ArchiveResponse<{
          user: {
            id: number;
            username: string;
            email: string;
            is_admin: boolean;
            change_initial_password: boolean;
            created_at: string;
            updated_at: string;
          };
        }>
      >('/auth/me');

      if (res.data.data.user.change_initial_password) {
        console.log('here in redirect');
        window.location = '/panel/pre-auth/change-password' as any;
        return {};
      }
    } catch (err) {
      console.log('login error', err);
      throw err;
    }
  },
  checkAuth: async (_) => {
    try {
      const res = await HttpAgent.get<
        ArchiveResponse<{
          user: {
            id: number;
            username: string;
            email: string;
            is_admin: boolean;
            change_initial_password: boolean;
            created_at: string;
            updated_at: string;
          };
        }>
      >('/auth/me');

      if (res.data.data.user.change_initial_password) {
        toast('please change your initial password', {
          type: 'error',
          position: toast.POSITION.BOTTOM_CENTER,
        });
        throw new Error('user should change initial password');
      }
    } catch (err) {
      console.log('check auth error: ', err);
      throw {};
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
    try {
      await HttpAgent.post<
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
      >('/auth/logout');
    } catch (err) {
      console.log('logout error: ', err);
      throw {};
    }
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
  getPermissions: async () => {
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
      if (res.data.data.user.is_admin) {
        console.log('was admin');
        return PERMISSIONS.ADMIN;
      }
      console.log('was user');
      return PERMISSIONS.USER;
    } catch (err) {
      console.log('check auth error: ', err);
      throw err;
    }
  },
};
