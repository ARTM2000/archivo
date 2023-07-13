import React, { useState } from 'react';
import {
  Typography,
  CircularProgress,
  Box,
  Grid,
  Button,
  TextField,
} from '@mui/material';
import { Logo } from '../branding/logo';
import { HttpAgent } from '../../utils/http-agent';
import { useNotify } from 'react-admin';
import { AxiosError } from 'axios';

export const RegisterAdmin = () => {
  const [username, setUsername] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const notify = useNotify();

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    HttpAgent.post('/auth/admin/register', {
      email,
      password,
      username,
    })
      .then(() => {
        notify('Admin registered', { type: 'success' });
        setTimeout(() => {
          window.location.reload();
        }, 1000);
      })
      .catch(
        (err: AxiosError<{ error: boolean; message: string; data: any }>) => {
          if (err.response?.status !== 500) {
            setLoading(false);
            notify(err.response?.data.message, { type: 'error' });
            return;
          }
          notify('Something went wrong :(', { type: 'error' });
          setLoading(false);
          console.log(err);
        },
      );
  };

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
        Register New Admin
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
                label="Email"
                value={email}
                onChange={(e: any) => setEmail(e.target.value)}
                fullWidth
              />
            </Grid>
            <Grid item>
              <TextField
                label="Username"
                value={username}
                onChange={(e: any) => setUsername(e.target.value)}
                fullWidth
              />
            </Grid>
            <Grid item>
              <TextField
                label="Password"
                type="password"
                value={password}
                onChange={(e: any) => setPassword(e.target.value)}
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
                  'Register Admin'
                )}
              </Button>
            </Grid>
          </Grid>
        </form>
      </Box>
    </>
  );
};
