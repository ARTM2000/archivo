import { Admin, Resource } from 'react-admin';
import './App.css';
import { DataProvider } from './data-provider';
import { AuthProvider } from './auth-provider';
import { LoginWrapper } from './components/auth/login-wrapper';
import { Dashboard } from './components/home/dashboard';
import { SourceServerList } from './components/sourceservers/list';
import StorageSharpIcon from '@mui/icons-material/StorageSharp';
import { SourceServerCreate } from './components/sourceservers/create';
import { Route, Routes } from 'react-router-dom';

function App() {
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
        ></Resource>
        <Routes>
          <Route></Route>
        </Routes>
      </Admin>
    </>
  );
}

export default App;