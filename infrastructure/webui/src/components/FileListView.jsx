// File List View component (table view) with multi-select support

import { useState, useEffect } from 'react';
import { Edit3, Check, X, Eye, Download, Trash2, FolderOpen, Archive, CheckSquare, Square, Lock, Unlock } from 'lucide-react';
import { FileIcon } from './FileIcon';
import { formatFileSize, canPreview } from '../utils/fileUtils';
import { useVault } from '../context/VaultContext';

export function FileListView({
    files,
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
}) {
    const [renamingItem, setRenamingItem] = useState(null);
    const [newName, setNewName] = useState('');

    // Vault state for encryption indicator
    const { isUnlocked } = useVault();

    // Trigger rename from context menu
    useEffect(() => {
        if (renameTarget) {
            startRename(renameTarget);
            onRenameComplete?.();
        }
    }, [renameTarget]);

    const startRename = (item) => {
        setRenamingItem(item);
        setNewName(item.name);
    };

    const handleRename = async () => {
        if (renamingItem && newName && newName !== renamingItem.name) {
            await onRename(renamingItem, newName);
        }
        setRenamingItem(null);
        setNewName('');
    };

    const cancelRename = () => {
        setRenamingItem(null);
        setNewName('');
    };

    if (files.length === 0) {
        return (
            <div className="py-12 text-center text-slate-400">
                <FolderOpen size={48} className="mx-auto mb-3 opacity-30" />
                <p className="text-sm">No files or folders</p>
            </div>
        );
    }

    return (
        <div className="overflow-x-auto -mx-6 px-6">
            <table className="w-full text-left border-collapse">
                <thead>
                    <tr className="text-xs text-slate-500 border-b border-white/5">
                        <th className="py-3 px-2 font-medium uppercase tracking-wider w-10"></th>
                        <th className="py-3 px-2 font-medium uppercase tracking-wider">Name</th>
                        <th className="py-3 px-2 font-medium uppercase tracking-wider">Size</th>
                        <th className="py-3 px-2 font-medium uppercase tracking-wider">Modified</th>
                        <th className="py-3 px-2 font-medium uppercase tracking-wider text-right">Actions</th>
                    </tr>
                </thead>
                <tbody className="text-sm">
                    {files.map((item) => {
                        const isRenaming = renamingItem?.name === item.name;
                        const selected = isSelected?.(item.name) || false;
                        const isEncrypted = item?.name?.endsWith('.enc') || false;

                        return (
                            <tr
                                key={item.name}
                                className={`group border-b border-white/5 last:border-0 hover:bg-white/5 transition-colors cursor-pointer ${selected ? 'bg-blue-500/10' : ''
                                    }`}
                                onClick={() => item.isDir && onNavigate(item)}
                                onContextMenu={(e) => onContextMenu?.(e, item)}
                            >
                                {/* Checkbox */}
                                <td className="py-4 px-2" onClick={(e) => e.stopPropagation()}>
                                    <button
                                        onClick={() => onToggleSelect?.(item.name)}
                                        className={`p-1 rounded transition-all ${selected
                                            ? 'text-blue-400'
                                            : 'text-slate-500 hover:text-slate-300'
                                            }`}
                                    >
                                        {selected ? <CheckSquare size={18} /> : <Square size={18} />}
                                    </button>
                                </td>

                                <td className="py-4 px-2 font-medium text-white">
                                    {isRenaming ? (
                                        <div className="flex items-center gap-2" onClick={(e) => e.stopPropagation()}>
                                            <div className={`p-2 rounded-lg ${item.isDir ? 'bg-blue-500/20 text-blue-400' : 'bg-slate-800 text-slate-400'}`}>
                                                <FileIcon name={item.name} isDir={item.isDir} size={16} />
                                            </div>
                                            <input
                                                type="text"
                                                value={newName}
                                                onChange={(e) => setNewName(e.target.value)}
                                                onKeyDown={(e) => {
                                                    if (e.key === 'Enter') handleRename();
                                                    if (e.key === 'Escape') cancelRename();
                                                }}
                                                className="flex-1 px-3 py-1.5 bg-slate-800 border border-white/10 rounded-lg text-white focus:outline-none focus:border-blue-500"
                                                autoFocus
                                            />
                                            <button onClick={handleRename} className="p-2 rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30">
                                                <Check size={14} />
                                            </button>
                                            <button onClick={cancelRename} className="p-2 rounded-lg bg-rose-500/20 text-rose-400 hover:bg-rose-500/30">
                                                <X size={14} />
                                            </button>
                                        </div>
                                    ) : (
                                        <div className="flex items-center gap-3">
                                            <div className={`p-2 rounded-lg ${item.isDir ? 'bg-blue-500/20 text-blue-400' : 'bg-slate-800 text-slate-400'}`}>
                                                <FileIcon name={item.name} isDir={item.isDir} size={16} />
                                            </div>
                                            {/* Encryption Status Indicator */}
                                            {isEncrypted && (
                                                <span className={`flex items-center ${isUnlocked ? 'text-emerald-400' : 'text-rose-400'}`}>
                                                    {isUnlocked ? <Unlock size={14} /> : <Lock size={14} />}
                                                </span>
                                            )}
                                            <span className={item.isDir ? "hover:text-blue-400 transition-colors" : ""}>
                                                {item.name}
                                            </span>
                                        </div>
                                    )}
                                </td>
                                <td className="py-4 px-2 text-slate-400 font-mono text-xs">
                                    {item.isDir ? "-" : formatFileSize(item.size)}
                                </td>
                                <td className="py-4 px-2 text-slate-400 text-xs">
                                    {new Date(item.modTime).toLocaleString()}
                                </td>
                                <td className="py-4 px-2 text-right">
                                    <div className="flex items-center justify-end gap-2" onClick={(e) => e.stopPropagation()}>
                                        {!item.isDir && canPreview(item.name) && (
                                            <button
                                                onClick={() => onPreview(item)}
                                                className="p-2 rounded-lg bg-violet-500/10 hover:bg-violet-500/20 text-violet-400 border border-violet-500/20 transition-all"
                                                title="Preview"
                                            >
                                                <Eye size={14} />
                                            </button>
                                        )}
                                        <button
                                            onClick={() => startRename(item)}
                                            className="p-2 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 border border-blue-500/20 transition-all"
                                            title="Rename"
                                        >
                                            <Edit3 size={14} />
                                        </button>
                                        <button
                                            onClick={() => onDownload(item)}
                                            className="p-2 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all"
                                            title={item.isDir ? "Download as ZIP" : "Download"}
                                        >
                                            {item.isDir ? <Archive size={14} /> : <Download size={14} />}
                                        </button>
                                        <button
                                            onClick={() => onDelete(item)}
                                            className="p-2 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 transition-all"
                                            title="Delete"
                                        >
                                            <Trash2 size={14} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        );
                    })}
                </tbody>
            </table>
        </div>
    );
}
