import {
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
    const { page, perPage } = params.pagination;
    const { field, order } = params.sort;
    const response = await HttpAgent.get<ArchiveResponse<{ list: any[], total: number }>>(
      `/${resource}/list`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
        },
        params: {
          sort_by: field,
          sort_order: order,
          start: (page - 1) * perPage,
          end: (page * perPage - 1),
          filter: JSON.stringify(params.filter),
        },
      },
    );

    const data = response.data.data;
    return {
      data: data.list,
      total: data.total,
    };
  },
};
