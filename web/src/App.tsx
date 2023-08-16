import { Admin, Resource } from 'react-admin';
import './App.css';
import { DataProvider } from './data-provider';
import { AuthProvider, PERMISSIONS } from './auth-provider';
import { LoginWrapper } from './components/auth/login-wrapper';
import { Dashboard } from './components/home/dashboard';
import { SourceServerList } from './components/sourceservers/list';
import StorageSharpIcon from '@mui/icons-material/StorageSharp';
import PeopleAltSharpIcon from '@mui/icons-material/PeopleAltSharp';
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
  });

  return (
    <>
      <Admin
        dataProvider={DataProvider as any}
        authProvider={AuthProvider}
        dashboard={Dashboard}
        loginPage={LoginWrapper}
        requireAuth
      >
        <Resource
          name="servers"
          list={SourceServerList}
          create={SourceServerCreate}
          hasEdit={false}
          icon={StorageSharpIcon}
        >
          <Route path=":serverId/:serverName/files" element={<FilesList />} />
          <Route
            path=":serverId/:serverName/files/:filename"
            element={<FileSnapshotsShow />}
          />
        </Resource>
        {perm === PERMISSIONS.ADMIN && (
          <Resource name="users" list={UserList} icon={PeopleAltSharpIcon} />
        )}
      </Admin>
    </>
  );
}

export default App;
