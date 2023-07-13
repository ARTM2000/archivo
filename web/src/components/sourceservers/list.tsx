import { Datagrid, List, TextField } from "react-admin"

export const SourceServerList = () => {

    return <List resource="servers">
        <Datagrid>
            <TextField source='id' />
            <TextField source='name' />
        </Datagrid>
    </List>
}
