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
import Settings from "./pages/Settings";
import Unlock from "./pages/Unlock";
import Setup from "./pages/Setup";
import FilePreviewDemo from "./pages/FilePreviewDemo";
import ErrorFallback from "./components/ErrorFallback";
import ConnectionFallbackModal from "./components/ConnectionFallbackModal";
import logger from "./utils/logger";
import { VaultProvider } from "./context/VaultContext";

export default function App() {
  return (
    <ErrorBoundary
      FallbackComponent={ErrorFallback}
      onReset={() => {
        // Reset application state here if needed
        window.location.href = "/dashboard";
      }}
      onError={(error, errorInfo) => {
        // FIX [BUG-JS-010]: Use production-safe logger
        logger.error("Uncaught error:", error, errorInfo);
      }}
    >
      <VaultProvider>
        {/* Global connection fallback modal for offline detection */}
        <ConnectionFallbackModal />
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/demo" element={<FilePreviewDemo />} />
            <Route path="/unlock" element={<Unlock />} />
            <Route path="/setup" element={<Setup />} />
            <Route path="/" element={<Layout />}>
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="dashboard" element={<Dashboard />} />
              <Route path="files" element={<Files />} />
              <Route path="files/vault" element={<Files initialPath="vault" />} />
              <Route path="search" element={<Search />} />
              <Route path="metrics" element={<Metrics />} />
              <Route path="backups" element={<Backup />} />
              <Route path="settings" element={<Settings />} />
              <Route path="ai" element={<Navigate to="/search" replace />} />
              <Route path="*" element={<Navigate to="/dashboard" replace />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </VaultProvider>
    </ErrorBoundary>
  );
}
