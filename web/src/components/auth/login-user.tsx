import {
  Box,
  CircularProgress,
  Grid,
  TextField,
  Typography,
  Button,
} from '@mui/material';
import { Logo } from '../branding/logo';
import React, { useState } from 'react';
import { useLogin, useNotify } from 'react-admin';
import { AxiosError } from 'axios';
import { ArchiveResponse } from '../../utils/types';

export const LoginUser = () => {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const notify = useNotify();
  const login = useLogin();

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    login({ email, password })
      .then(() => {
        notify('Welcome!', { type: 'success', autoHideDuration: 3000 });
      })
      .catch((err: AxiosError<ArchiveResponse>) => {
        if (err.response?.status !== 500) {
          notify(err.response?.data.message, { type: 'error' });
          return;
        }
        notify('Something went wrong :(', { type: 'error' });
      })
      .finally(() => {
        setLoading(false);
      });
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
        Login
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
                  'Login'
                )}
              </Button>
            </Grid>
          </Grid>
        </form>
      </Box>
    </>
  );
};
