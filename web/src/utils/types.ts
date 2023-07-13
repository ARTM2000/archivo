export type ArchiveResponse<T = any> = {
  track_id: string;
  error: boolean;
  message: string;
  data: T;
};
