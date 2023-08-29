import {
  Box,
  Button,
  CircularProgress,
  Grid,
  TextField,
  Typography,
} from '@mui/material';
import { Logo } from '../branding/logo';
import { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { HttpAgent } from '../../utils/http-agent';
import { AxiosError } from 'axios';

export const NewUserChangePass = () => {
  const [oldPass, setOldPass] = useState<string>('');
  const [newPass, setNewPass] = useState<string>('');
  const [newPassRepeat, setNewPassRepeat] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    console.log(oldPass, newPass, newPassRepeat);

    if (newPass.trim() !== newPassRepeat.trim()) {
      toast('new password and password confirm are not same', {
        type: 'error',
        position: toast.POSITION.BOTTOM_CENTER,
      });
      return;
    }

    setLoading(true);
    HttpAgent.post('/pre-auth/change-user-initial-pass', {
      initial_password: oldPass,
      new_password: newPass,
    })
      .then((res) => {
        console.log(res.data.message);
        toast(res.data.message, {
          type: 'success',
          position: toast.POSITION.BOTTOM_CENTER,
        });
        window.location = '/panel/login' as any;
      })
      .catch((err: AxiosError) => {
        console.log(err.response?.data);
        if (err.status?.toString().match(/^5/)) {
          toast('Internal server error', {
            type: 'error',
            position: toast.POSITION.BOTTOM_CENTER,
          });
          return;
        }
        toast((err.response?.data as any).message, {
          type: 'error',
          position: toast.POSITION.BOTTOM_CENTER,
        });
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {}, []);

  return (
    <>
      <Box maxWidth={'210px'} margin={'10vh auto 0 auto'}>
        <Logo width="200px" />
      </Box>
      <Typography
        margin={'1vh auto 2vh auto'}
        variant="h4"
        align="center"
        gutterBottom
      >
        Change Your Password First
      </Typography>
      <Box height={'100px'} margin={'10px auto'} maxWidth={'500px'}>
        <form
          onSubmit={handleFormSubmit}
          style={{
            border: '2px solid lightgray',
            borderRadius: '10px',
            padding: '15px',
          }}
        >
          <Grid container spacing={2} direction={'column'} minWidth={'400px'}>
            <Grid item>
              <TextField
                label="Current Password"
                type="password"
                value={oldPass}
                onChange={(e: any) => setOldPass(e.target.value)}
                fullWidth
              />
            </Grid>
            <Grid item>
              <TextField
                label="New Password"
                type="password"
                value={newPass}
                onChange={(e: any) => setNewPass(e.target.value)}
                fullWidth
              />
            </Grid>
            <Grid item>
              <TextField
                label="New Password Confirm"
                type="password"
                value={newPassRepeat}
                onChange={(e: any) => setNewPassRepeat(e.target.value)}
                fullWidth
              />
            </Grid>
            <Grid item>
              <Button
                type="submit"
                variant="contained"
                fullWidth
                disabled={loading}
              >
                {loading ? (
                  <CircularProgress size={24} color="info" />
                ) : (
                  'Change Password'
                )}
              </Button>
            </Grid>
          </Grid>
        </form>
      </Box>
    </>
  );
};
