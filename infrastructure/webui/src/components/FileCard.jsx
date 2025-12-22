// File Card component for Grid View with selection support

import { useState, useEffect, useCallback } from 'react';
import { Edit3, Check, X, Eye, Download, Trash2, Archive, CheckSquare, Square } from 'lucide-react';
import { FileIcon } from './FileIcon';
import { formatFileSize, canPreview } from '../utils/fileUtils';
import { useLongPress } from '../hooks/useLongPress';

export function FileCard({
    item,
    onNavigate,
    onPreview,
    onRename,
    onDownload,
    onDelete,
    onToggleSelect,
    isSelected,
    onContextMenu,
    renameTarget,
    onRenameComplete,
    onMoveFile,
}) {
    const [isRenaming, setIsRenaming] = useState(false);
    const [newName, setNewName] = useState('');
    const [isDragOver, setIsDragOver] = useState(false);

    // Long-press handler for mobile preview
    const handleLongPress = useCallback(() => {
        if (!item.isDir && canPreview(item.name)) {
            onPreview?.(item);
        } else {
            // Open context menu on long-press for folders or non-previewable files
            onContextMenu?.({ preventDefault: () => { }, clientX: 100, clientY: 100 }, item);
        }
    }, [item, onPreview, onContextMenu]);

    const longPressProps = useLongPress(handleLongPress, {
        delay: 500,
    });

    const selected = isSelected?.(item.name) || false;

    // Trigger rename from context menu
    useEffect(() => {
        if (renameTarget && renameTarget.name === item.name) {
            startRename();
            onRenameComplete?.();
        }
    }, [renameTarget]);

    const startRename = () => {
        setIsRenaming(true);
        setNewName(item.name);
    };

    const handleRename = async () => {
        if (newName && newName !== item.name) {
            await onRename(item, newName);
        }
        setIsRenaming(false);
        setNewName('');
    };

    const cancelRename = () => {
        setIsRenaming(false);
        setNewName('');
    };

    // Drag handlers
    const handleDragStart = (e) => {
        e.dataTransfer.setData('application/json', JSON.stringify(item));
        e.dataTransfer.effectAllowed = 'move';
    };

    const handleDragOver = (e) => {
        if (!item.isDir) return;
        e.preventDefault();
        e.dataTransfer.dropEffect = 'move';
        setIsDragOver(true);
    };

    const handleDragLeave = () => {
        setIsDragOver(false);
    };

    const handleDrop = async (e) => {
        if (!item.isDir) return;
        e.preventDefault();
        setIsDragOver(false);

        try {
            const sourceItem = JSON.parse(e.dataTransfer.getData('application/json'));
            if (sourceItem.name === item.name) return; // Can't drop on itself
            await onMoveFile?.(sourceItem, item);
        } catch (err) {
            console.error('Drop failed:', err);
        }
    };

    return (
        <div
            className={`group relative overflow-hidden rounded-xl border transition-all cursor-pointer ${selected
                ? 'border-blue-500/50 bg-blue-500/10'
                : isDragOver
                    ? 'border-emerald-500/50 bg-emerald-500/10'
                    : 'border-white/10 bg-slate-900/40 dark:bg-slate-900/40 hover:bg-white/5'
                }`}
            onClick={() => item.isDir && onNavigate(item)}
            onContextMenu={(e) => onContextMenu?.(e, item)}
            draggable={!isRenaming}
            onDragStart={handleDragStart}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            {...longPressProps}
        >
            {/* Selection Checkbox */}
            <div
                className="absolute top-2 left-2 z-10"
                onClick={(e) => e.stopPropagation()}
            >
                <button
                    onClick={() => onToggleSelect?.(item.name)}
                    className={`p-1 rounded transition-all ${selected
                        ? 'text-blue-400 bg-blue-500/20'
                        : 'text-slate-500 hover:text-slate-300 opacity-0 group-hover:opacity-100 bg-slate-800/80'
                        }`}
                >
                    {selected ? <CheckSquare size={16} /> : <Square size={16} />}
                </button>
            </div>

            <div className="p-4 flex flex-col items-center text-center">
                {/* Icon */}
                <div className={`p-4 rounded-xl mb-3 ${item.isDir ? 'bg-blue-500/20 text-blue-400' : 'bg-slate-800/50 text-slate-400'} group-hover:scale-110 transition-transform`}>
                    <FileIcon name={item.name} isDir={item.isDir} size={32} />
                </div>

                {/* Name */}
                {isRenaming ? (
                    <div className="w-full flex items-center gap-1" onClick={(e) => e.stopPropagation()}>
                        <input
                            type="text"
                            value={newName}
                            onChange={(e) => setNewName(e.target.value)}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter') handleRename();
                                if (e.key === 'Escape') cancelRename();
                            }}
                            className="flex-1 px-2 py-1 text-xs bg-slate-800 border border-white/10 rounded text-white focus:outline-none focus:border-blue-500"
                            autoFocus
                        />
                        <button onClick={handleRename} className="p-1 rounded bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30">
                            <Check size={14} />
                        </button>
                        <button onClick={cancelRename} className="p-1 rounded bg-rose-500/20 text-rose-400 hover:bg-rose-500/30">
                            <X size={14} />
                        </button>
                    </div>
                ) : (
                    <p className="text-sm font-medium text-white truncate w-full px-2 group-hover:text-blue-400 transition-colors">
                        {item.name}
                    </p>
                )}

                {/* Size */}
                {!item.isDir && (
                    <p className="text-xs text-slate-500 mt-1">{formatFileSize(item.size)}</p>
                )}

                {/* Actions */}
                <div className="flex items-center gap-1 mt-3 opacity-0 group-hover:opacity-100 transition-opacity" onClick={(e) => e.stopPropagation()}>
                    {!item.isDir && canPreview(item.name) && (
                        <button
                            onClick={() => onPreview(item)}
                            className="p-1.5 rounded-lg bg-violet-500/10 hover:bg-violet-500/20 text-violet-400 border border-violet-500/20 transition-all"
                            title="Preview"
                        >
                            <Eye size={12} />
                        </button>
                    )}
                    <button
                        onClick={startRename}
                        className="p-1.5 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 border border-blue-500/20 transition-all"
                        title="Rename"
                    >
                        <Edit3 size={12} />
                    </button>
                    <button
                        onClick={() => onDownload(item)}
                        className="p-1.5 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all"
                        title={item.isDir ? "Download as ZIP" : "Download"}
                    >
                        {item.isDir ? <Archive size={12} /> : <Download size={12} />}
                    </button>
                    <button
                        onClick={() => onDelete(item)}
                        className="p-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 transition-all"
                        title="Delete"
                    >
                        <Trash2 size={12} />
                    </button>
                </div>
            </div>
        </div>
    );
}
