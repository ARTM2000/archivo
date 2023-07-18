import { TableCell, TableRow, Checkbox } from '@mui/material';
import React from 'react';
import {
  Datagrid,
  DatagridBody,
  RecordContextProvider,
  List,
  TextField,
  DateField,
} from 'react-admin';
import { useNavigate } from 'react-router-dom';

const MyDatagridRow = (props: {
  record: { id: number; name: string };
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
      <TableRow
        onClick={() => {
          console.log('record > ', record);
          history(`${id}/files`);
        }}
      >
        <TableCell padding="checkbox">
          {selectable && (
            <Checkbox
              checked={selected}
              onClick={(event) => onToggleItem(id, event)}
            />
          )}
        </TableCell>
        {React.Children.map(children, (field) => (
          <TableCell key={`${id}-${field.props.source}`}>{field}</TableCell>
        ))}
      </TableRow>
    </RecordContextProvider>
  );
};

const MyDatagridBody = (props: any) => (
  <DatagridBody {...props} row={<MyDatagridRow {...props} />} />
);
const MyDatagrid = (props: any) => (
  <Datagrid {...props} body={<MyDatagridBody />} />
);

export const SourceServerList = () => {
  return (
    <List resource="servers">
      <MyDatagrid>
        <TextField source="id" label="ID" />
        <TextField source="name" label="Name" />
        <DateField
          source="created_at"
          label="Created at"
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
