import { Box, Grid } from '@mui/material';
import { useEffect, useState } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { HttpAgent } from '../../utils/http-agent';
import { ChartData, ChartRange, ChartMetric } from '../../utils/types';
import { CHART_RANGE, RangePicker } from './interval-picker';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
);

export const options = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top' as const,
    },
    title: {
      display: false,
    },
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        stepSize: 1,
      },
    },
  },
};

export const ActivityChart = (props: { sourceServerName?: string }) => {
  const { sourceServerName } = props;

  const [chartData, setChartData] = useState<ChartData[]>([]);
  const [chartLabels, setChartLabels] = useState<string[]>([]);
  const [chartUpdateInterval, setChartUpdateInterval] = useState(5000);
  const [currentChartRange, setCurrentChartRange] =
    useState<ChartRange>(CHART_RANGE[4]);

  const formatChartData = (
    metrics: ChartMetric[],
  ): { data: ChartData[]; labels: string[] } => {
    const finalData: ChartData[] = [
      {
        label: 'Fail',
        borderColor: '#c40606',
        backgroundColor: '#c40606',
        data: [],
      },
      {
        label: 'Success',
        borderColor: '#16c406',
        backgroundColor: '#16c406',
        data: [],
      },
    ];

    const finalLabels: string[] = [];

    // find out current bucket size
    const firstBucket = metrics[0];
    const bucketSize =
      new Date(firstBucket.to).getTime() - new Date(firstBucket.from).getTime();
    setChartUpdateInterval(bucketSize);

    for (let i = 0; i < metrics.length; i++) {
      const metric = metrics[i];
      const metricTimeFrom = new Date(metric.from);

      finalLabels.push(
        `${metricTimeFrom.getFullYear()}-${
          metricTimeFrom.getMonth() + 1 < 10
            ? `0${metricTimeFrom.getMonth() + 1}`
            : metricTimeFrom.getMonth() + 1
        }-${
          metricTimeFrom.getDate() < 10
            ? `0${metricTimeFrom.getDate()}`
            : metricTimeFrom.getDate()
        } ${
          metricTimeFrom.getHours() < 10
            ? `0${metricTimeFrom.getHours()}`
            : metricTimeFrom.getHours()
        }:${
          metricTimeFrom.getMinutes() < 10
            ? `0${metricTimeFrom.getMinutes()}`
            : metricTimeFrom.getMinutes()
        }:${
          metricTimeFrom.getSeconds() < 10
            ? `0${metricTimeFrom.getSeconds()}`
            : metricTimeFrom.getSeconds()
        }`,
      );

      finalData[1].data.push(metric.total_success);
      finalData[0].data.push(metric.total_fail);
    }

    return {
      data: finalData,
      labels: finalLabels,
    };
  };

  const fetchAllChartData = (chartRange: ChartRange) => {
    const now = new Date();

    const params: {
      to: number;
      from: number;
      srv_name?: string;
    } = {
      to: now.valueOf(),
      from: new Date(now.valueOf() - chartRange.duration).valueOf(),
    };
    let url = '/dashboard/metrics/activities';

    if (sourceServerName) {
      params.srv_name = sourceServerName;
      url += '/single-server';
    }

    HttpAgent.get(url, {
      params,
    })
      .then((res) => res.data)
      .then(
        (data: {
          data: {
            metrics: ChartMetric[];
          };
        }) => {
          const formattedData = formatChartData(data.data.metrics);
          setChartData(formattedData.data);
          setChartLabels(formattedData.labels);
        },
      )
      .catch((err) => {
        console.log('got error in fetchAllChartData', err);
      });
  };

  useEffect(() => {
    const interval = setInterval(
      () => fetchAllChartData(currentChartRange),
      chartUpdateInterval,
    );
    return () => clearInterval(interval);
  }, [currentChartRange]);

  useEffect(() => {
    fetchAllChartData(currentChartRange);
  }, []);

  return (
    <Grid container justifyContent={'center'} spacing={4}>
      <RangePicker
        currentRange={currentChartRange}
        setCurrentRange={(intrvl: ChartRange) =>
          setCurrentChartRange(intrvl)
        }
      />
      <Grid item xs={12}>
        <Box
          component={'div'}
          width={'auto'}
          height={'250px'}
          minHeight={'50vh'}
          mx={'auto'}
          justifyContent={'center'}
          alignContent={'center'}
          sx={{
            backgroundColor: 'white',
            boxShadow: '0 3px 8px #bcd2f5',
            borderRadius: '8px',
            padding: '15px',
            marginBottom: '28px',
          }}
        >
          <Line
            options={options}
            data={{
              labels: chartLabels,
              datasets: chartData,
            }}
            width="100%"
          />
        </Box>
      </Grid>
    </Grid>
  );
};
