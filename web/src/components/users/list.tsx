import {
  BooleanField,
  CreateButton,
  Datagrid,
  DateField,
  ExportButton,
  List,
  TextField,
  TopToolbar,
} from 'react-admin';

const ActionsList = () => (
  <TopToolbar>
    <CreateButton label="Register New User" />
    <ExportButton />
  </TopToolbar>
);

export const UserList = () => {
  return (
    <List resource="users" actions={<ActionsList />}>
      <Datagrid size="medium">
        <TextField source="id" label="ID" />
        <TextField source="username" label="Username" />
        <TextField source="email" label="Email" />
        <BooleanField source="is_admin" label="Admin" />
        <DateField
          source="created_at"
          label="Joined at"
          options={{
            hour: '2-digit',
            minute: '2-digit',
            second: 'numeric',
            weekday: 'short',
            year: 'numeric',
            month: 'numeric',
            day: 'numeric',
          }}
        />
        <DateField
          source="last_login_at"
          label="Last Login"
          options={{
            hour: '2-digit',
            minute: '2-digit',
            second: 'numeric',
            weekday: 'short',
            year: 'numeric',
            month: 'numeric',
            day: 'numeric',
          }}
        />
      </Datagrid>
    </List>
  );
};
