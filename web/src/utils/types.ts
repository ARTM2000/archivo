export type ArchiveResponse<T = any> = {
  track_id: string;
  error: boolean;
  message: string;
  data: T;
};

export type ChartInterval = {
  title: string;
  duration: number;
};

export type ChartMetric = {
  from: string;
  to: string;
  total_success: number;
  total_fail: number;
  details: {
    [key: string]: {
      success_count: number;
      fail_count: number;
    };
  };
};

export type ChartData = {
  label: string;
  data: number[];
  borderColor: string;
  backgroundColor?: string;
};
