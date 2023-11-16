import { Box, CircularProgress, Grid, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import { HttpAgent } from '../../utils/http-agent';
import { toast } from 'react-toastify';
import { Title } from 'react-admin';
import {
  ActivityChart,
  CHART_INTERVAL,
} from '../activity-chart/activity-chart';

export const Dashboard = () => {
  const [backupFileCount, setBackupFileCount] = useState<number>(0);
  const [snapshotsSize, setSnapshotsSize] = useState<string>('0 B');
  const [sourceServersCount, setSourceServersCount] = useState<number>(0);

  const [commonLoading, setCommonLoading] = useState<boolean>(true);

  const getCommonMetrics = () => {
    setCommonLoading(true);
    HttpAgent.get('/dashboard/metrics/common')
      .then((res) => {
        const data = res.data as {
          data: {
            backup_files_count: number;
            snapshot_occupied_size: string;
            source_servers_count: number;
          };
        };

        setBackupFileCount(data.data.backup_files_count);
        setSnapshotsSize(data.data.snapshot_occupied_size);
        setSourceServersCount(data.data.source_servers_count);
        setCommonLoading(false);
      })
      .catch((err) => {
        console.log('error in gathering common metrics', err);
        setCommonLoading(false);
        toast('Something went wrong :(', {
          type: 'error',
          position: toast.POSITION.BOTTOM_CENTER,
        });
      });
  };

  useEffect(() => {
    getCommonMetrics();
  }, []);

  return (
    <Box sx={{ marginTop: '20px' }}>
      <Title title={'Dashboard'} />
      <Box sx={{ flexGrow: 1, margin: '0 20px' }}>
        <Grid container justifyContent={'left'} alignItems={'end'}></Grid>
        <ActivityChart currentChartInterval={CHART_INTERVAL[4]} />
        <Grid container spacing={4} justifyContent={'center'}>
          <MetricInfo
            title="Total source servers"
            value={sourceServersCount}
            loading={commonLoading}
          />
          <MetricInfo
            title="Total files for backup"
            value={backupFileCount}
            loading={commonLoading}
          />
          <MetricInfo
            title="Total snapshots size"
            value={snapshotsSize}
            loading={commonLoading}
          />
        </Grid>
      </Box>
    </Box>
  );
};

const MetricInfo = (props: { title: string; value: any; loading: boolean }) => {
  return (
    <Grid item md={4} xs={12}>
      <Box
        component={'div'}
        width={'auto'}
        height={'250px'}
        mx={'auto'}
        justifyContent={'center'}
        alignContent={'center'}
        sx={{
          backgroundColor: 'white',
          boxShadow: '0 3px 8px #bcd2f5',
          borderRadius: '8px',
          padding: '15px',
        }}
      >
        <Typography variant="h6">{props.title}</Typography>
        {props.loading ? (
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'center',
              paddingTop: '22%',
            }}
          >
            <CircularProgress
              sx={{ justifyContent: 'center', alignContent: 'center' }}
            />
          </Box>
        ) : (
          <Box sx={{ paddingTop: '15%', textAlign: 'center' }}>
            <Typography variant="h3">{props.value}</Typography>
          </Box>
        )}
      </Box>
    </Grid>
  );
};
