import { Checkbox, IconButton, TableCell, TableRow } from '@mui/material';
import {
  Button,
  Datagrid,
  DatagridBody,
  DateField,
  ExportButton,
  List,
  RecordContextProvider,
  TextField,
  TopToolbar,
} from 'react-admin';
import { useParams } from 'react-router-dom';
import DownloadIcon from '@mui/icons-material/Download';
import React from 'react';
import ArrowBackIosNewSharpIcon from '@mui/icons-material/ArrowBackIosNewSharp';
import { toast } from 'react-toastify';

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

  const downloadSnapshot = () => {
    const url = `${import.meta.env.VITE_ARCHIVO_API_PANEL_BASE_URL}/servers/${
      params.serverId
    }/files/${params.filename}/${record.name}/download`;
    window.location.href = url;
    toast('Download started', {
      type: 'success',
      position: toast.POSITION.BOTTOM_CENTER,
    });
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

const ActionsList = (props: {
  serverId: string;
  serverName: string;
  filename: string;
}) => (
  <TopToolbar>
    <Button
      label="Back to Files"
      onClick={() => {
        window.location =
          `/panel/servers/${props.serverId}/${props.serverName}/files` as any;
      }}
    >
      <ArrowBackIosNewSharpIcon />
    </Button>
    <ExportButton />
  </TopToolbar>
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
      actions={
        <ActionsList
          serverId={params.serverId as string}
          serverName={params.serverName as string}
          filename={params.filename as string}
        />
      }
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
