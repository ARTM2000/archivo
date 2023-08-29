import './App.css';
import 'react-toastify/dist/ReactToastify.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { NewUserChangePass } from './components/auth/new-user-change-pass';
import { AdminPanel } from './Admin';
import { ToastContainer } from 'react-toastify';

function App() {
  return (
    <BrowserRouter basename="/panel">
      <Routes>
        <Route
          path="/pre-auth/change-password"
          element={<NewUserChangePass />}
        />
        <Route path="/*" element={<AdminPanel />} />
      </Routes>
      <ToastContainer />
    </BrowserRouter>
  );
}

export default App;
