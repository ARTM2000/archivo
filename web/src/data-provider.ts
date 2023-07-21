import {
  CreateParams,
  GetListParams,
  GetListResult,
  DataProvider as IDataProvider,
} from 'react-admin';
import { HttpAgent } from './utils/http-agent';
import { ArchiveResponse } from './utils/types';
import { TOKEN_KEY } from './auth-provider';

export const DataProvider: Partial<IDataProvider> = {
  getList: async (
    resource: string,
    params: GetListParams,
  ): Promise<GetListResult<any>> => {
    let url = `/${resource}`;
    if (resource === 'files') {
      url = `/servers/${params.meta.serverId}/files`;
    }
    if (resource === 'snapshot') {
      url = `/servers/${params.meta.serverId}/files/${params.meta.filename}`;
      if (params.sort.field === 'id')
        params.sort.field = params.meta.sort.DefaultBy;
    }

    const { page, perPage } = params.pagination;
    const { field, order } = params.sort;

    const response = await HttpAgent.get<
      ArchiveResponse<{ list: any[]; total: number }>
    >(url, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
      },
      params: {
        sort_by: field,
        sort_order: order,
        start: (page - 1) * perPage,
        end: page * perPage,
        filter: JSON.stringify(params.filter),
      },
    });

    const data = response.data.data;
    return {
      data: data.list || [],
      total: data.total,
    };
  },

  create: async (resource: string, params: CreateParams<any>) => {
    const response = await HttpAgent.post<
      ArchiveResponse<{ id: number; [key: string]: any }>
    >(`/${resource}/new`, params.data, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
      },
    });

    const data = response.data.data;
    return {
      data: {
        ...params.data,
        id: data.id,
      },
    };
  },
};
