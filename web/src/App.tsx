import { Admin, Resource, CustomRoutes } from "react-admin"
import './App.css'
import { DataProvider } from "./data-provider"
import { AuthProvider } from "./auth-provider"
import { LoginWrapper } from "./components/auth/login-wrapper"

function App() {
  return (
    <>
    <Admin dataProvider={DataProvider} authProvider={AuthProvider} loginPage={LoginWrapper} requireAuth>
      <CustomRoutes>

      </CustomRoutes>
      <Resource name="srcsrv"></Resource>
    </Admin>
    </>
  )
}

export default App
