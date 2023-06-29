/// <reference types="vite/client" />

interface ImportMetaEnv {
	readonly VITE_ARCHIVE1_PANEL_BASE_URL: string;
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
