import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Layout from "./Layout";
import Files from "./pages/Files";
import Metrics from "./pages/Metrics";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Backup from "./pages/Backup";
import Dashboard from "./pages/Dashboard";
import Search from "./pages/Search";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<Layout />}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="files" element={<Files />} />
          <Route path="search" element={<Search />} />
          <Route path="metrics" element={<Metrics />} />
          <Route path="backups" element={<Backup />} />
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
