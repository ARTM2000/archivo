import {
  Grid,
  MenuItem,
  TextField,
} from '@mui/material';
import { ChartRange } from '../../utils/types';

export const CHART_RANGE: ChartRange[] = [
  {
    title: '30 seconds',
    duration: 30 * 1000,
  },
  {
    title: '1 minute',
    duration: 60 * 1000,
  },
  {
    title: '2 minutes',
    duration: 2 * 60 * 1000,
  },
  {
    title: '5 minutes',
    duration: 5 * 60 * 1000,
  },
  {
    title: '10 minutes',
    duration: 10 * 60 * 1000,
  },
  {
    title: '30 minutes',
    duration: 30 * 60 * 1000,
  },
  {
    title: '1 hour',
    duration: 60 * 60 * 1000,
  },
  {
    title: '2 hours',
    duration: 2 * 60 * 60 * 1000,
  },
  {
    title: '6 hours',
    duration: 6 * 60 * 60 * 1000,
  },
];

export const RangePicker = (props: {
  currentRange: ChartRange;
  setCurrentRange: (cIntvl: ChartRange) => void;
}) => {
  const onIntervalPickerChange = (e: any) => {
    props.setCurrentRange(CHART_RANGE[CHART_RANGE.findIndex(ch => ch.duration === +e.target.value)]);
  };

  return (
    <Grid item xs={12} md={4}>
      <TextField select
        // labelId="interval-picker-select"
        id="interval-picker-select"
        value={props.currentRange.duration}
        onChange={onIntervalPickerChange}
        label="Chart Range"
        fullWidth 
      >
        {CHART_RANGE.map((chIntrvl, i) => (
          <MenuItem key={i} value={chIntrvl.duration}>{chIntrvl.title}</MenuItem>
        ))}
      </TextField>
    </Grid>
  );
};
