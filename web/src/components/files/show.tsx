import { Checkbox, IconButton, TableCell, TableRow } from '@mui/material';
import {
  Datagrid,
  DatagridBody,
  DateField,
  List,
  RecordContextProvider,
  TextField,
  useNotify,
} from 'react-admin';
import { useParams } from 'react-router-dom';
import DownloadIcon from '@mui/icons-material/Download';
import React from 'react';

const MyDatagridRow = (props: {
  record: {
    id: number;
    name: string;
    size: string;
    checksum: string;
    created_at: Date;
  };
  id: number;
  onToggleItem: Function;
  children: any;
  selected: boolean;
  selectable: boolean;
}) => {
  const { record, id, onToggleItem, children, selected, selectable } = props;
  // const history = useNavigate();
  const params = useParams();
  const notify = useNotify();

  const downloadSnapshot = () => {
    const url = `${import.meta.env.VITE_ARCHIVE1_API_PANEL_BASE_URL}/servers/${
      params.serverId
    }/files/${params.filename}/${record.name}/download`;
    window.location.href = url;
    notify('Download started', { type: 'success' });
  };

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
        {React.Children.map(children, (field) => {
          return field.props?.className === 'dl-icon' ? (
            <TableCell
              key={`${id}-${field.props.source}`}
              onClick={(_) => {
                downloadSnapshot();
              }}
            >
              {field}
            </TableCell>
          ) : (
            <TableCell key={`${id}-${field.props.source}`}>{field}</TableCell>
          );
        })}
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

export const FileSnapshotsShow = () => {
  const params = useParams();

  return (
    <List
      resource="snapshot"
      queryOptions={{
        meta: {
          serverId: params.serverId,
          filename: params.filename,
          sort: { DefaultBy: 'name' },
        },
      }}
      title={`Servers > ${params.serverName} > ${params.filename} (snapshots)`}
      exporter={false}
    >
      <MyDatagrid>
        <TextField source="id" label="ID" />
        <TextField source="name" label="Name" />
        <TextField source="size" label="Size" />
        <TextField source="checksum" label="Checksum (sha256)" />
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
        <IconButton className="dl-icon">
          <DownloadIcon />
        </IconButton>
      </MyDatagrid>
    </List>
  );
};
