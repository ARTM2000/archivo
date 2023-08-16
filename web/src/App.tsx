import {
  Admin,
  Resource,
  useGetPermissions,
  usePermissions,
} from 'react-admin';
import './App.css';
import { DataProvider } from './data-provider';
import { AuthProvider, PERMISSIONS } from './auth-provider';
import { LoginWrapper } from './components/auth/login-wrapper';
import { Dashboard } from './components/home/dashboard';
import { SourceServerList } from './components/sourceservers/list';
import StorageSharpIcon from '@mui/icons-material/StorageSharp';
import { SourceServerCreate } from './components/sourceservers/create';
import { Route } from 'react-router-dom';
import { FilesList } from './components/files/list';
import { FileSnapshotsShow } from './components/files/show';
import { useEffect, useState } from 'react';
import { UserList } from './components/users/list';

function App() {
  const [perm, setPerm] = useState<PERMISSIONS>(PERMISSIONS.USER);
  useEffect(() => {
    AuthProvider.getPermissions('').then((perm: PERMISSIONS) => {
      setPerm(perm);
    });
  }, []);

  return (
    <>
      <Admin
        dataProvider={DataProvider as any}
        authProvider={AuthProvider}
        dashboard={Dashboard}
        loginPage={LoginWrapper}
        darkTheme={{ palette: { mode: 'dark' } }}
        requireAuth
      >
        <Resource
          name="servers"
          list={SourceServerList}
          create={SourceServerCreate}
          hasEdit={false}
          icon={StorageSharpIcon}
        >
          <Route path=":serverId/files" element={<FilesList />} />
          <Route
            path=":serverId/files/:filename"
            element={<FileSnapshotsShow />}
          />
        </Resource>
        {perm === PERMISSIONS.ADMIN && (
          <Resource name="users" list={UserList} />
        )}
      </Admin>
    </>
  );
}

export default App;
