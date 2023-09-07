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
import { AxiosError } from 'axios';
import { ArchiveResponse } from '../../utils/types';
import { toast } from 'react-toastify';

export const RegisterAdmin = () => {
  const [username, setUsername] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    HttpAgent.post('/auth/admin/register', {
      email,
      password,
      username,
    })
      .then(() => {
        toast('Admin registered', {
          type: 'success',
          position: toast.POSITION.BOTTOM_CENTER,
        });
        setTimeout(() => {
          window.location.reload();
        }, 1000);
      })
      .catch((err: AxiosError<ArchiveResponse<any>>) => {
        setLoading(false);
        if (err.response?.status !== 500) {
          toast(err.response?.data.message, {
            type: 'error',
            position: toast.POSITION.BOTTOM_CENTER,
          });
          return;
        }
        toast('Something went wrong :(', {
          type: 'error',
          position: toast.POSITION.BOTTOM_CENTER,
        });
        console.log(err);
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
