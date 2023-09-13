/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_ARCHIVO_API_PANEL_BASE_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
