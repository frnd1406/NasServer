import { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { deriveKey, arrayBufferToBase64, base64ToUint8Array } from '../lib/crypto';
import { apiRequest } from '../lib/api';

const VaultContext = createContext(null);

export function VaultProvider({ children }) {
    const [key, setKey] = useState(null); // CryptoKey (Raw) - Never stored in localStorage
    const [isUnlocked, setIsUnlocked] = useState(false);
    const [vaultConfig, setVaultConfig] = useState(null); // Salt, config loaded from server
    const [isLoading, setIsLoading] = useState(true);

    // Load vault metadata on mount
    useEffect(() => {
        checkVaultStatus();
    }, []);

    const checkVaultStatus = async () => {
        try {
            // Try to fetch hidden meta file
            // We assume it's at /vault/.meta.json
            // Using standard download endpoint which might fail if file doesn't exist
            const res = await fetch('/api/v1/storage/download?path=vault/.meta.json');
            if (res.ok) {
                const meta = await res.json();
                setVaultConfig(meta);
            } else {
                setVaultConfig(null); // Not setup
            }
        } catch (e) {
            console.warn("Vault metadata check failed:", e);
        } finally {
            setIsLoading(false);
        }
    };

    const lock = useCallback(() => {
        setKey(null);
        setIsUnlocked(false);
    }, []);

    const unlock = useCallback(async (password) => {
        if (!vaultConfig?.salt) {
            throw new Error("Vault configuration missing");
        }

        const salt = base64ToUint8Array(vaultConfig.salt);
        const derivedKey = await deriveKey(password, salt);

        // Verify key? We can't verify unless we try to decrypt something or verify a hash.
        // For now, we accept it. If it's wrong, decryption will fail later.
        // Ideally validationHash is stored in meta.

        setKey(derivedKey);
        setIsUnlocked(true);
    }, [vaultConfig]);

    const setup = useCallback(async (data) => {
        // data = { key: CryptoKey, salt: string (base64) }

        // 1. Create vault directory
        try {
            await apiRequest('/api/v1/storage/mkdir', {
                method: 'POST',
                body: JSON.stringify({ path: 'vault' })
            });
        } catch (e) {
            // Ignore if exists
        }

        // 2. Upload metadata
        const meta = {
            salt: data.salt,
            created: new Date().toISOString(),
            version: 1
        };

        const blob = new Blob([JSON.stringify(meta, null, 2)], { type: 'application/json' });
        const formData = new FormData();
        formData.append('file', blob, '.meta.json'); // Filename

        // Upload to /vault/
        // Note: Access token is now sent via HttpOnly cookie
        await fetch('/api/v1/storage/upload?path=vault', {
            method: 'POST',
            body: formData,
            credentials: 'include', // Send auth cookie
            headers: {
                'X-CSRF-Token': localStorage.getItem('csrfToken') || ''
            }
        });

        setKey(data.key);
        setVaultConfig(meta);
        setIsUnlocked(true);
    }, []);

    return (
        <VaultContext.Provider value={{
            isUnlocked,
            key,
            vaultConfig,
            isLoading,
            lock,
            unlock,
            setup
        }}>
            {children}
        </VaultContext.Provider>
    );
}

export const useVault = () => useContext(VaultContext);
