import { Admin, Resource, CustomRoutes } from "react-admin"
import './App.css'
import { DataProvider } from "./data-provider"

function App() {
  return (
    <>
    <Admin dataProvider={DataProvider}>
      <CustomRoutes>
        
      </CustomRoutes>
      <Resource name="srcsrv"></Resource>
    </Admin>
    </>
  )
}

export default App
