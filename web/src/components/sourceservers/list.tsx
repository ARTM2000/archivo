import { TableCell, TableRow, Checkbox } from '@mui/material';
import React from 'react';
import {
  Datagrid,
  DatagridBody,
  RecordContextProvider,
  List,
  TextField,
} from 'react-admin';

const MyDatagridRow = (props: {
  record: any;
  id: any;
  onToggleItem: Function;
  children: any;
  selected: boolean;
  selectable: boolean;
}) => {
  const { record, id, onToggleItem, children, selected, selectable } = props;
  return (
    <RecordContextProvider value={record}>
      <TableRow onClick={() => console.log('record > ', record)}>
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
        <TextField source="id" title="ID" />
        <TextField source="name" title="Name" />
      </MyDatagrid>
    </List>
  );
};
