import { useEffect, useState } from 'react';
import { AxiosResponse } from 'axios';
import { RegisterAdmin } from './register-admin';
import { Typography } from '@mui/material';
import { HttpAgent } from '../../utils/http-agent';
import { useNotify } from 'react-admin';

export const LoginWrapper = () => {
  const [loading, setLoading] = useState<boolean>(true);
  const [adminExist, setAdminExist] = useState<boolean>(false);
  const notify = useNotify();

  useEffect(() => {
    HttpAgent.get('/auth/admin/existence')
      .then(
        (
          res: AxiosResponse<
            { data: { admin_exist: boolean }; error: boolean },
            any
          >,
        ) => {
          if (!res.data.error) {
            setAdminExist(res.data.data.admin_exist);
            setLoading(false);
          }
        },
      )
      .catch((err) => {
        notify('Something went wrong :(', { type: 'error' });
        setLoading(false);
        console.log(err);
      });
  }, []);

  return loading ? (
    <Typography variant="h4" style={{ textAlign: 'center', marginTop: '20vh' }}>
      Just a moment...
    </Typography>
  ) : (
    <>
      {!adminExist && <RegisterAdmin />}
      {adminExist && <h1>login user</h1>}
    </>
  );
};
