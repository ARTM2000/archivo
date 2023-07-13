import { Admin, Resource } from 'react-admin';
import './App.css';
import { DataProvider } from './data-provider';
import { AuthProvider } from './auth-provider';
import { LoginWrapper } from './components/auth/login-wrapper';
import { Dashboard } from './components/home/dashboard';

function App() {
  return (
    <>
      <Admin
        dataProvider={DataProvider}
        authProvider={AuthProvider}
        dashboard={Dashboard}
        loginPage={LoginWrapper}
        requireAuth
      >
        <Resource name="srcsrv"></Resource>
      </Admin>
    </>
  );
}

export default App;
