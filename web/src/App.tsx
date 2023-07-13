import { Admin, Resource } from 'react-admin';
import './App.css';
import { DataProvider } from './data-provider';
import { AuthProvider } from './auth-provider';
import { LoginWrapper } from './components/auth/login-wrapper';
import { Dashboard } from './components/home/dashboard';
import { SourceServerList } from './components/sourceservers/list';

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
        <Resource name="servers" list={SourceServerList}></Resource>
      </Admin>
    </>
  );
}

export default App;
