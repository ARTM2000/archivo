import { Dialog, DialogContent, DialogTitle, Typography } from '@mui/material';
import { useState } from 'react';
import {
  Create,
  SimpleForm,
  TextInput,
  required,
  useRedirect,
} from 'react-admin';
import { HttpAgent } from '../../utils/http-agent';
import { ArchiveResponse } from '../../utils/types';
import { AxiosError } from 'axios';
import { toast } from 'react-toastify';

export const SourceServerCreate = () => {
  const [open, setOpen] = useState<boolean>(false);
  const [apiKey, setApiKey] = useState<string>('');
  const redirect = useRedirect();

  const handleSubmit = (data: any) => {
    HttpAgent.post<
      ArchiveResponse<{ id: number; name: string; api_key: string }>
    >(`/servers/new`, data)
      .then((res) => {
        toast('New Source Server Created', { type: 'success', position: toast.POSITION.BOTTOM_CENTER });
        setApiKey(res.data.data.api_key);
        setOpen(true);
      })
      .catch((err: AxiosError<ArchiveResponse>) => {
        if (err.response?.status !== 500) {
          toast(err.response?.data.message, { type: 'error', position: toast.POSITION.BOTTOM_CENTER });
          return;
        }
        toast('Internal Server Error', { type: 'error', position: toast.POSITION.BOTTOM_CENTER });
      });
  };

  const handleDialogClose = () => {
    setOpen(false);
    redirect('list', 'servers');
  };

  return (
    <>
      <Create title={'Create Source Server'}>
        <SimpleForm onSubmit={handleSubmit}>
          <TextInput source="name" validate={required('name is required')} />
        </SimpleForm>
      </Create>
      <Dialog open={open} onClose={handleDialogClose} maxWidth="sm">
        <DialogTitle>Credentials</DialogTitle>
        <DialogContent>
          <Typography variant="h6">API Key</Typography>
          <Typography
            variant="body2"
            component="code"
            style={{ display: 'block' }}
          >
            {apiKey}
          </Typography>
          <br />
          <Typography
            variant="body2"
            fontWeight="lightbold"
            component="small"
            style={{ color: 'red' }}
          >
            Warning: Keep this credential somewhere safe as it is not accessible
            after quitting this page
          </Typography>
        </DialogContent>
      </Dialog>
    </>
  );
};
