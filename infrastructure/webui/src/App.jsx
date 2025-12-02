import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { ErrorBoundary } from "react-error-boundary";
import Layout from "./Layout";
import Files from "./pages/Files";
import Metrics from "./pages/Metrics";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Backup from "./pages/Backup";
import Dashboard from "./pages/Dashboard";
import Search from "./pages/Search";
import ErrorFallback from "./components/ErrorFallback";

export default function App() {
  return (
    <ErrorBoundary
      FallbackComponent={ErrorFallback}
      onReset={() => {
        // Reset application state here if needed
        window.location.href = "/dashboard";
      }}
      onError={(error, errorInfo) => {
        // Log error to monitoring service (e.g., Sentry)
        console.error("Uncaught error:", error, errorInfo);
      }}
    >
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
    </ErrorBoundary>
  );
}
