// Custom hook for file storage operations

import { useState, useCallback } from 'react';
import { authHeaders } from '../utils/auth';
import { joinPath } from '../utils/fileUtils';

const API_BASE = import.meta.env.VITE_API_BASE_URL || window.location.origin;

export function useFileStorage() {
    const [files, setFiles] = useState([]);
    const [trashedFiles, setTrashedFiles] = useState([]);
    const [path, setPath] = useState('/');
    const [loading, setLoading] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState('');

    // Load files from a directory
    const loadFiles = useCallback(async (target = path) => {
        if (!authHeaders().Authorization) {
            setError('Bitte zuerst einloggen.');
            return;
        }
        setLoading(true);
        setError('');
        try {
            const res = await fetch(
                `${API_BASE}/api/v1/storage/files?path=${encodeURIComponent(target)}`,
                {
                    credentials: 'include',
                    headers: authHeaders()
                }
            );
            if (res.status === 401) {
                setError('Authentifizierung fehlgeschlagen. Bitte neu einloggen.');
                return;
            }
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}: Failed to load files`);
            }
            const data = await res.json();
            // Filter out .trash directory
            const filteredFiles = (data.items || []).filter(item => item.name !== '.trash');
            setFiles(filteredFiles);
        } catch (err) {
            setError(err.message || 'Unknown error');
        } finally {
            setLoading(false);
        }
    }, [path]);

    // Load trashed files
    const loadTrash = useCallback(async () => {
        try {
            const res = await fetch(
                `${API_BASE}/api/v1/storage/trash`,
                {
                    credentials: 'include',
                    headers: authHeaders()
                }
            );
            if (res.ok) {
                const data = await res.json();
                setTrashedFiles(data.items || []);
            }
        } catch (err) {
            console.error('Failed to load trash:', err);
        }
    }, []);

    // Upload files
    const uploadFiles = useCallback(async (filesToUpload, currentPath) => {
        if (!filesToUpload || filesToUpload.length === 0) return;

        setUploading(true);
        setError('');

        try {
            for (let i = 0; i < filesToUpload.length; i++) {
                const file = filesToUpload[i];
                console.log(`Uploading file ${i + 1}/${filesToUpload.length}:`, file.name);

                const form = new FormData();
                form.append('file', file);
                form.append('path', currentPath);

                const headers = authHeaders();
                delete headers['Content-Type'];

                const csrfToken = localStorage.getItem('csrfToken') || localStorage.getItem('csrf_token');
                if (csrfToken) {
                    headers['X-CSRF-Token'] = csrfToken;
                }

                const res = await fetch(`${API_BASE}/api/v1/storage/upload`, {
                    method: 'POST',
                    body: form,
                    credentials: 'include',
                    headers: headers,
                });

                if (res.status === 401) {
                    setError('Authentifizierung fehlgeschlagen.');
                    return;
                }
                if (!res.ok) {
                    const errorText = await res.text().catch(() => 'No error details');
                    throw new Error(`Upload failed for ${file.name}: HTTP ${res.status} - ${errorText}`);
                }
            }

            await loadFiles(currentPath);
        } catch (err) {
            console.error('Upload error:', err);
            setError(err.message || 'Unknown error');
        } finally {
            setUploading(false);
        }
    }, [loadFiles]);

    // Download file
    const downloadFile = useCallback(async (item, currentPath) => {
        const target = joinPath(currentPath, item.name);
        try {
            const res = await fetch(
                `${API_BASE}/api/v1/storage/download?path=${encodeURIComponent(target)}`,
                {
                    credentials: 'include',
                    headers: authHeaders(),
                }
            );
            if (!res.ok) {
                setError(`Download failed: HTTP ${res.status}`);
                return;
            }
            const blob = await res.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = item.name;
            a.click();
            window.URL.revokeObjectURL(url);
        } catch (err) {
            setError(err.message);
        }
    }, []);

    // Delete file
    const deleteFile = useCallback(async (item, currentPath) => {
        if (item.name === '.trash') {
            setError('Der Papierkorb kann nicht gelöscht werden');
            return;
        }

        if (!window.confirm(`"${item.name}" wirklich löschen?`)) return;

        const target = joinPath(currentPath, item.name);
        try {
            const res = await fetch(
                `${API_BASE}/api/v1/storage/delete?path=${encodeURIComponent(target)}`,
                {
                    method: 'DELETE',
                    credentials: 'include',
                    headers: authHeaders(),
                }
            );
            if (!res.ok) {
                setError(`Delete failed: HTTP ${res.status}`);
                return;
            }
            await loadFiles(currentPath);
            await loadTrash();
        } catch (err) {
            setError(err.message);
        }
    }, [loadFiles, loadTrash]);

    // Restore file from trash
    const restoreFile = useCallback(async (item, currentPath) => {
        try {
            const res = await fetch(
                `${API_BASE}/api/v1/storage/restore?path=${encodeURIComponent(item.originalPath || item.name)}`,
                {
                    method: 'POST',
                    credentials: 'include',
                    headers: authHeaders(),
                }
            );
            if (!res.ok) {
                setError(`Restore failed: HTTP ${res.status}`);
                return;
            }
            await loadTrash();
            await loadFiles(currentPath);
        } catch (err) {
            setError(err.message);
        }
    }, [loadFiles, loadTrash]);

    // Rename file
    const renameFile = useCallback(async (item, newName, currentPath) => {
        if (!newName || newName === item.name) return false;

        const oldPath = joinPath(currentPath, item.name);
        const newPath = joinPath(currentPath, newName);

        try {
            const res = await fetch(
                `${API_BASE}/api/v1/storage/rename`,
                {
                    method: 'POST',
                    credentials: 'include',
                    headers: {
                        ...authHeaders(),
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ oldPath, newPath }),
                }
            );
            if (!res.ok) {
                throw new Error(`Rename failed: HTTP ${res.status}`);
            }
            await loadFiles(currentPath);
            return true;
        } catch (err) {
            setError(err.message);
            return false;
        }
    }, [loadFiles]);

    // Create folder
    const createFolder = useCallback(async (folderName, currentPath) => {
        if (!folderName || folderName.trim() === '') {
            setError('Ordnername darf nicht leer sein');
            return false;
        }

        try {
            const folderPath = joinPath(currentPath, folderName.trim());
            const res = await fetch(`${API_BASE}/api/v1/storage/mkdir`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    ...authHeaders(),
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ path: folderPath }),
            });

            if (!res.ok) {
                throw new Error(`Ordner erstellen fehlgeschlagen: HTTP ${res.status}`);
            }

            await loadFiles(currentPath);
            return true;
        } catch (err) {
            setError(err.message || 'Fehler beim Erstellen des Ordners');
            return false;
        }
    }, [loadFiles]);

    // Navigate to directory
    const navigateTo = useCallback((newPath) => {
        setPath(newPath);
        loadFiles(newPath);
    }, [loadFiles]);

    // Go up one directory
    const goUp = useCallback(() => {
        if (path === '/') return;
        const parts = path.split('/').filter(Boolean);
        parts.pop();
        const parent = parts.length ? `/${parts.join('/')}` : '/';
        navigateTo(parent);
    }, [path, navigateTo]);

    // Clear error
    const clearError = useCallback(() => setError(''), []);

    // Batch download multiple files as ZIP
    const batchDownload = useCallback(async (selectedNames, currentPath) => {
        if (!selectedNames || selectedNames.length === 0) {
            setError('Keine Dateien ausgewählt');
            return;
        }

        try {
            const paths = selectedNames.map(name => joinPath(currentPath, name));

            const res = await fetch(`${API_BASE}/api/v1/storage/batch-download`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    ...authHeaders(),
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ paths }),
            });

            if (!res.ok) {
                throw new Error(`Batch download failed: HTTP ${res.status}`);
            }

            const blob = await res.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'download.zip';
            a.click();
            window.URL.revokeObjectURL(url);
        } catch (err) {
            setError(err.message);
        }
    }, []);

    // Download folder as ZIP
    const downloadFolderAsZip = useCallback(async (folderPath) => {
        try {
            const res = await fetch(
                `${API_BASE}/api/v1/storage/download-zip?path=${encodeURIComponent(folderPath)}`,
                {
                    credentials: 'include',
                    headers: authHeaders(),
                }
            );

            if (!res.ok) {
                throw new Error(`ZIP download failed: HTTP ${res.status}`);
            }

            const blob = await res.blob();
            const folderName = folderPath.split('/').filter(Boolean).pop() || 'folder';
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${folderName}.zip`;
            a.click();
            window.URL.revokeObjectURL(url);
        } catch (err) {
            setError(err.message);
        }
    }, []);

    // Delete multiple files
    const batchDelete = useCallback(async (selectedNames, currentPath) => {
        if (!selectedNames || selectedNames.length === 0) return;

        const count = selectedNames.length;
        if (!window.confirm(`${count} Dateien/Ordner wirklich löschen?`)) return;

        try {
            for (const name of selectedNames) {
                const target = joinPath(currentPath, name);
                const res = await fetch(
                    `${API_BASE}/api/v1/storage/delete?path=${encodeURIComponent(target)}`,
                    {
                        method: 'DELETE',
                        credentials: 'include',
                        headers: authHeaders(),
                    }
                );
                if (!res.ok) {
                    console.warn(`Failed to delete ${name}: HTTP ${res.status}`);
                }
            }
            await loadFiles(currentPath);
            await loadTrash();
        } catch (err) {
            setError(err.message);
        }
    }, [loadFiles, loadTrash]);

    return {
        // State
        files,
        trashedFiles,
        path,
        loading,
        uploading,
        error,
        // Actions
        loadFiles,
        loadTrash,
        uploadFiles,
        downloadFile,
        deleteFile,
        restoreFile,
        renameFile,
        createFolder,
        navigateTo,
        goUp,
        clearError,
        // Batch Actions
        batchDownload,
        downloadFolderAsZip,
        batchDelete,
    };
}

