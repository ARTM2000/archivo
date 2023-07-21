import { Checkbox, TableCell, TableRow } from '@mui/material';
import React from 'react';
import {
  Datagrid,
  DatagridBody,
  DateField,
  List,
  RecordContextProvider,
  TextField,
} from 'react-admin';
import { useNavigate, useParams } from 'react-router-dom';

const MyDatagridRow = (props: {
  record: { id: number; filename: string; snapshots: number; updated_at: Date };
  id: number;
  onToggleItem: Function;
  children: any;
  selected: boolean;
  selectable: boolean;
}) => {
  const { record, id, onToggleItem, children, selected, selectable } = props;
  const history = useNavigate();
  // const params = useParams();

  return (
    <RecordContextProvider value={record}>
      <TableRow
        onClick={() => {
          history(`${record.filename}`);
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
  <Datagrid
    {...props}
    style={{ cursor: 'pointer' }}
    body={<MyDatagridBody />}
  />
);

export const FilesList = () => {
  const params = useParams();

  return (
    <List
      resource="files"
      queryOptions={{ meta: { serverId: params.serverId } }}
    >
      <MyDatagrid>
        <TextField source="id" label="ID" />
        <TextField source="filename" label="Filename" />
        <TextField source="snapshots" label="Snapshots" />
        <DateField
          source="updated_at"
          label="Updated at"
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
