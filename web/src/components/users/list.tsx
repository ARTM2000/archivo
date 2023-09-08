import { Checkbox, TableCell, TableRow } from '@mui/material';
import React from 'react';
import {
  BooleanField,
  CreateButton,
  Datagrid,
  DatagridBody,
  DateField,
  ExportButton,
  List,
  RecordContextProvider,
  TextField,
  TopToolbar,
} from 'react-admin';
import { useNavigate } from 'react-router-dom';

const MyDatagridRow = (props: {
  record: { id: number; username: string; email: string; is_admin: boolean };
  id: number;
  onToggleItem: Function;
  children: any;
  selected: boolean;
  selectable: boolean;
}) => {
  const { record, id, onToggleItem, children, selected, selectable } = props;
  const history = useNavigate();
  return (
    <RecordContextProvider value={record}>
      <TableRow>
        <TableCell padding="checkbox">
          {selectable && (
            <Checkbox
              checked={selected}
              onClick={(event) => onToggleItem(id, event)}
            />
          )}
        </TableCell>
        {React.Children.map(children, (field) => (
          <TableCell
            key={`${id}-${field.props.source}`}
            onClick={() => {
              history(`${id}/${record.username}/activities`);
            }}
          >
            {field}
          </TableCell>
        ))}
      </TableRow>
    </RecordContextProvider>
  );
};

const MyDatagridBody = (props: any) => (
  <DatagridBody {...props} row={<MyDatagridRow {...props} />} />
);
const MyDatagrid = (props: any) => (
  <Datagrid
    {...props}
    style={{ cursor: 'pointer' }}
    body={<MyDatagridBody />}
    size="medium"
  />
);

const ActionsList = () => (
  <TopToolbar>
    <CreateButton label="Register New User" />
    <ExportButton />
  </TopToolbar>
);

export const UserList = () => {
  return (
    <List resource="users" actions={<ActionsList />}>
      <MyDatagrid size="medium">
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
      </MyDatagrid>
    </List>
  );
};
