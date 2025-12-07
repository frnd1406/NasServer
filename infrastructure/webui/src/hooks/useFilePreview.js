// Custom hook for file preview functionality

import { useState, useCallback } from 'react';
import { authHeaders } from '../utils/auth';
import { joinPath, isImage, isText } from '../utils/fileUtils';

const API_BASE = import.meta.env.VITE_API_BASE_URL || window.location.origin;

export function useFilePreview() {
    const [previewItem, setPreviewItem] = useState(null);
    const [previewContent, setPreviewContent] = useState(null);
    const [previewLoading, setPreviewLoading] = useState(false);

    const openPreview = useCallback(async (item, currentPath) => {
        if (item.isDir) return;

        setPreviewItem(item);
        setPreviewLoading(true);
        setPreviewContent(null);

        const target = joinPath(currentPath, item.name);

        try {
            if (isImage(item.name)) {
                const res = await fetch(
                    `${API_BASE}/api/v1/storage/download?path=${encodeURIComponent(target)}`,
                    {
                        credentials: 'include',
                        headers: authHeaders(),
                    }
                );
                if (res.ok) {
                    const blob = await res.blob();
                    const url = window.URL.createObjectURL(blob);
                    setPreviewContent({ type: 'image', url });
                }
            } else if (isText(item.name)) {
                const res = await fetch(
                    `${API_BASE}/api/v1/storage/download?path=${encodeURIComponent(target)}`,
                    {
                        credentials: 'include',
                        headers: authHeaders(),
                    }
                );
                if (res.ok) {
                    const text = await res.text();
                    setPreviewContent({ type: 'text', content: text });
                }
            }
        } catch (err) {
            console.error('Preview failed:', err);
        } finally {
            setPreviewLoading(false);
        }
    }, []);

    const closePreview = useCallback(() => {
        if (previewContent?.type === 'image' && previewContent.url) {
            window.URL.revokeObjectURL(previewContent.url);
        }
        setPreviewItem(null);
        setPreviewContent(null);
    }, [previewContent]);

    return {
        previewItem,
        previewContent,
        previewLoading,
        openPreview,
        closePreview,
    };
}
