import { Datagrid, DateField, List, TextField } from 'react-admin';
import { useParams } from 'react-router-dom';

export const FilesList = () => {
  const params = useParams();

  return (
    <List
      resource="files"
      queryOptions={{ meta: { serverId: params.serverId } }}
    >
      <Datagrid>
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
      </Datagrid>
    </List>
  );
};
