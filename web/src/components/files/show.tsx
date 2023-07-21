import { Datagrid, DateField, List, TextField } from 'react-admin';
import { useParams } from 'react-router-dom';

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
    >
      <Datagrid>
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
      </Datagrid>
    </List>
  );
};
