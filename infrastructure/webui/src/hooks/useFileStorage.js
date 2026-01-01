// Custom hook for file storage operations

import { useState, useCallback } from 'react';
import { authHeaders } from '../utils/auth';
import { joinPath } from '../utils/fileUtils';
import { uploadFile as apiUploadFile } from '../lib/api';

const API_BASE = import.meta.env.VITE_API_BASE_URL || window.location.origin;

import { calculateTotalChunks, fileChunkIterator } from '../utils/chunking';
import { encryptChunk, generateIV, arrayBufferToBase64 } from '../lib/crypto';

export function useFileStorage(initialPath = '/', vaultKey = null) {
    const [files, setFiles] = useState([]);
    const [trashedFiles, setTrashedFiles] = useState([]);
    const [path, setPath] = useState(initialPath);
    const [loading, setLoading] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState('');

    // Load files from a directory
    const loadFiles = useCallback(async (target = path) => {
        // Check if user might be logged in (CSRF token exists)
        // Note: Access token is now in HttpOnly cookie, can't check directly
        const csrfToken = localStorage.getItem('csrfToken') || localStorage.getItem('csrf_token');
        if (!csrfToken) {
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
                window.location.href = '/login';
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

    // Upload files with optional encryption mode override
    const uploadFiles = useCallback(async (filesToUpload, currentPath, options = {}) => {
        if (!filesToUpload || filesToUpload.length === 0) return;

        const { encryptionMode = 'auto' } = options;

        setUploading(true);
        setError('');

        try {
            for (let i = 0; i < filesToUpload.length; i++) {
                const file = filesToUpload[i];
                console.log(`Uploading file ${i + 1}/${filesToUpload.length}:`, file.name);

                // CHECK: Is this a vault upload?
                // We define vault upload if currentPath starts with "vault" or "/vault" AND we have a key
                const isVault = (currentPath.startsWith('vault') || currentPath.startsWith('/vault')) && vaultKey;

                if (isVault) {
                    // === ENCRYPTED CHUNKED UPLOAD ===
                    if (!vaultKey) throw new Error("Vault locked: Cannot upload");

                    // 1. Init Upload
                    const iv = generateIV(); // Unique IV per file
                    // We need to store IV. Usually prepend to file? Or separate metadata.
                    // Simple approach: Prepend 12 bytes IV to the first chunk.
                    // The backend just sees "data".

                    // Init backend session
                    // Filename must end in .enc
                    const encFilename = file.name + ".enc";
                    const initRes = await fetch(`${API_BASE}/api/v1/vault/upload/init`, {
                        method: 'POST',
                        headers: { ...authHeaders(), 'Content-Type': 'application/json' },
                        body: JSON.stringify({ filename: encFilename, total_size: file.size }) // Size is aprox, actually encrypted size is slightly larger + IV
                    });

                    if (!initRes.ok) throw new Error("Vault upload init failed");
                    const { upload_id } = await initRes.json();

                    // 2. Upload Chunks
                    let chunkIndex = 0;
                    for await (const chunkCtx of fileChunkIterator(file)) {
                        // Encrypt chunk
                        // deriveKey is expensive, we rely on cached `vaultKey` (CryptoKey)

                        // If first chunk, we should prepend IV?
                        // Actually, AES-GCM needs IV for decryption.
                        // Strategy: The IV is static for the file? Or random?
                        // If random IV per chunk -> we need to store IVs.
                        // If static IV per file -> we reuse IV? (Insecure if same key+IV used for different data blocks? yes)
                        // AES-GCM is tricky with streaming.
                        // Usually: Encrypt file as one stream? WebCrypto doesn't support streaming encryption easily.
                        // We must encrypt CHUNKS. Each chunk needs its own IV or counter?
                        // "True Zero Knowledge" with large files usually implies:
                        // 1. Generate random FILE KEY.
                        // 2. Encrypt FILE KEY with Vault Key (Key Wrapping).
                        // 3. Encrypt chunks with FILE KEY + Counter (or unique IVs).

                        // SIMPLIFICATION for this prototype:
                        // We use the SAME IV for the whole file? NO, catastrophic.
                        // We generate a unique IV for EACH chunk?
                        // Then we need to store the IV with the chunk.
                        // Output format: [IV 12b][Ciphertext][Tag] ... [IV 12b][Ciphertext][Tag] (Tag is inclusive in ciphertext usually in WebCrypto)

                        const chunkIV = generateIV();
                        const encryptedBuffer = await encryptChunk(chunkCtx.data, vaultKey, chunkIV);

                        // Append IV to start of buffer
                        const combined = new Uint8Array(chunkIV.length + encryptedBuffer.byteLength);
                        combined.set(chunkIV);
                        combined.set(new Uint8Array(encryptedBuffer), chunkIV.length);

                        // Upload
                        const chunkRes = await fetch(`${API_BASE}/api/v1/vault/upload/chunk/${upload_id}`, {
                            method: 'POST',
                            headers: {
                                ...authHeaders(),
                                'Content-Type': 'application/octet-stream'
                            },
                            body: combined
                        });

                        if (!chunkRes.ok) throw new Error(`Chunk ${chunkIndex} upload failed`);
                        chunkIndex++;
                    }

                    // 3. Finalize
                    await fetch(`${API_BASE}/api/v1/vault/upload/finalize/${upload_id}`, {
                        method: 'POST',
                        headers: { ...authHeaders(), 'Content-Type': 'application/json' },
                        body: JSON.stringify({ path: 'vault' }) // Force into vault dir
                    });

                } else {
                    // === STANDARD UPLOAD (Via API) ===
                    // Pass encryption mode override to API
                    await apiUploadFile(file, currentPath, encryptionMode);
                }
            }

            await loadFiles(currentPath);
        } catch (err) {
            console.error('Upload error:', err);
            setError(err.message || 'Unknown error');
            if (err.message === 'Unauthorized') {
                window.location.href = '/login';
            }
        } finally {
            setUploading(false);
        }
    }, [loadFiles, vaultKey]);

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

    // Move file or folder to a new location
    const moveFile = useCallback(async (sourcePath, destinationPath, currentPath) => {
        try {
            const res = await fetch(`${API_BASE}/api/v1/storage/move`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    ...authHeaders(),
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ sourcePath, destinationPath }),
            });

            if (!res.ok) {
                const data = await res.json().catch(() => ({}));
                throw new Error(data.error || `Move failed: HTTP ${res.status}`);
            }

            await loadFiles(currentPath);
            return true;
        } catch (err) {
            setError(err.message);
            return false;
        }
    }, [loadFiles]);

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
        moveFile,
    };
}

