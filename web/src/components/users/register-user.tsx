import { AxiosError } from 'axios';
import {
  Create,
  PasswordInput,
  SimpleForm,
  TextInput,
  required,
  useNotify,
} from 'react-admin';
import { ArchiveResponse } from '../../utils/types';

export const RegisterUser = () => {
  const notify = useNotify();
  const onError = (err: Error, _: any, __: any) => {
    const error = err as AxiosError<any, ArchiveResponse>;
    const message = error.response?.data.message;
    notify(message, { type: 'error' });
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
