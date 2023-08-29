import { AxiosError } from 'axios';
import {
  Create,
  PasswordInput,
  SimpleForm,
  TextInput,
  required,
} from 'react-admin';
import { ArchiveResponse } from '../../utils/types';
import { toast } from 'react-toastify';

export const RegisterUser = () => {
  const onError = (err: Error, _: any, __: any) => {
    const error = err as AxiosError<any, ArchiveResponse>;
    const message = error.response?.data.message;
    toast(message, { type: 'error', position: toast.POSITION.BOTTOM_CENTER });
  };

  return (
    <Create
      mutationOptions={{
        onError,
      }}
    >
      <SimpleForm>
        <TextInput
          source="email"
          label="Email"
          validate={[required('Email is required')]}
        />
        <TextInput
          source="username"
          label="Username"
          validate={[required('Username is required')]}
        />
        <PasswordInput
          source="password"
          label="Password"
          validate={[required('Password is required')]}
        />
      </SimpleForm>
    </Create>
  );
};
