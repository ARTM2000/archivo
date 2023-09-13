import { Box } from '@mui/material';
import { Datagrid, DateField, List, TextField } from 'react-admin';
import { useParams } from 'react-router-dom';

export const ShowUserActivities = () => {
  const params = useParams();
  return (
    <>
      <Box>
        <p style={{ marginLeft: '10px', marginTop: '20px', marginBottom: 0 }}>
          Username: {params.username}
        </p>
      </Box>
      <List
        resource="user_activities"
        queryOptions={{
          meta: {
            userId: params.userId,
          },
        }}
      >
        <Datagrid size="medium">
          <TextField source="id" label="ID" />
          <TextField source="act" label="Activity" />
          <DateField
            source="created_at"
            label="Created At"
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
    </>
  );
};
