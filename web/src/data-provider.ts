import jsonServerProvider from 'ra-data-json-server';
import { Options, fetchUtils } from 'react-admin';

const httpClient = (url: string, options: Options | undefined) => {
  console.log('request data-provider >>', { url, options });
  return fetchUtils.fetchJson(url, options);
};

export const DataProvider = jsonServerProvider(
  import.meta.env.VITE_ARCHIVE1_API_PANEL_BASE_URL,
  httpClient,
);
