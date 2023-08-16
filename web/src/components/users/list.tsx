import { BooleanField, Datagrid, DateField, List, TextField } from "react-admin";

export const UserList = () => {
    return (
      <List resource="users">
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
        </Datagrid>
      </List>
    );
  };