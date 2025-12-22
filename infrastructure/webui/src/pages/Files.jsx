// Files page - With multi-select, batch operations, and search filter

import { useEffect, useState, useRef, useCallback, useMemo } from 'react';
import { Loader2, Search, X, Download, Trash2, CheckSquare, Square } from 'lucide-react';

// Hooks
import { useFileStorage } from '../hooks/useFileStorage';
import { useFilePreview } from '../hooks/useFilePreview';
import { useDragAndDrop } from '../hooks/useDragAndDrop';
import { useFileSelection } from '../hooks/useFileSelection';

// Utils
import { getBreadcrumbs, joinPath } from '../utils/fileUtils';

// Components
import { GlassCard } from '../components/ui/GlassCard';
import { FileToolbar } from '../components/FileToolbar';
import { FileGridView } from '../components/FileGridView';
import { FileListView } from '../components/FileListView';
import { TrashView } from '../components/TrashView';
import { DragDropOverlay } from '../components/DragDropOverlay';
import { NewFolderModal } from '../components/NewFolderModal';
import { FilePreviewModal } from '../components/FilePreviewModal';
import { ContextMenu } from '../components/ContextMenu';

import { VaultModal } from '../components/Vault/VaultModal';
import { useVault } from '../context/VaultContext';

export default function Files({ initialPath = '/' }) {
  // Vault context
  const { isUnlocked, unlock, setup, vaultConfig, key } = useVault();
  const [showVaultModal, setShowVaultModal] = useState(false);

  // File storage operations hook
  const {
    files,
    trashedFiles,
    path,
    loading,
    uploading,
    error,
    loadFiles,
    loadTrash,
    uploadFiles,
    downloadFile,
    deleteFile,
    restoreFile,
    renameFile,
    createFolder,
    navigateTo,
    batchDownload,
    downloadFolderAsZip,
    batchDelete,
    moveFile,
  } = useFileStorage(initialPath, isUnlocked ? key : null);

  // Check vault access
  useEffect(() => {
    if (path.startsWith('vault') || path.startsWith('/vault')) {
      if (!isUnlocked) {
        setShowVaultModal(true);
      }
    }
  }, [path, isUnlocked]);

  const handleVaultUnlock = async (password) => {
    await unlock(password);
    setShowVaultModal(false);
    loadFiles(path); // Reload files now that we have keys (if we were using listing API that needed keys, but we are using raw listing. However, maybe we trigger re-decryption?)
  };

  const handleVaultSetup = async (data) => {
    await setup(data);
    setShowVaultModal(false);
    loadFiles(path);
  };


  // File preview hook
  const {
    previewItem,
    previewContent,
    previewLoading,
    openPreview,
    closePreview,
  } = useFilePreview();

  // Search/Filter state
  const [searchQuery, setSearchQuery] = useState('');

  // Filtered files based on search
  const filteredFiles = useMemo(() => {
    if (!searchQuery.trim()) return files;
    const query = searchQuery.toLowerCase();
    return files.filter(f => f.name.toLowerCase().includes(query));
  }, [files, searchQuery]);

  // Multi-selection hook
  const {
    selectedCount,
    selectedItems,
    allSelected,
    toggleSelect,
    clearSelection,
    toggleSelectAll,
    isSelected,
  } = useFileSelection(filteredFiles);

  // View state
  const [viewMode, setViewMode] = useState('list');
  const [showTrash, setShowTrash] = useState(false);
  const [showNewFolderModal, setShowNewFolderModal] = useState(false);

  // Context menu state
  const [contextMenu, setContextMenu] = useState({ isOpen: false, position: { x: 0, y: 0 }, item: null });

  // File input ref for upload button
  const fileInputRef = useRef(null);

  // Drag and drop
  const handleFilesDropped = useCallback((droppedFiles) => {
    uploadFiles(droppedFiles, path);
  }, [uploadFiles, path]);

  const { isDragging, dragProps } = useDragAndDrop(handleFilesDropped);

  // Load files on mount and path change
  useEffect(() => {
    loadFiles('/');
    loadTrash();
  }, []);

  // Clear selection when path changes
  useEffect(() => {
    clearSelection();
    setSearchQuery('');
  }, [path, clearSelection]);

  // Handlers
  const handleNavigate = useCallback((item) => {
    if (!item.isDir) return;
    const nextPath = joinPath(path, item.name);
    navigateTo(nextPath);
  }, [path, navigateTo]);

  const handlePreview = useCallback((item) => {
    openPreview(item, path);
  }, [openPreview, path]);

  const handleRename = useCallback(async (item, newName) => {
    return await renameFile(item, newName, path);
  }, [renameFile, path]);

  const handleDownload = useCallback((item) => {
    if (item.isDir) {
      // Download folder as ZIP
      downloadFolderAsZip(joinPath(path, item.name));
    } else {
      downloadFile(item, path);
    }
  }, [downloadFile, downloadFolderAsZip, path]);

  const handleDelete = useCallback((item) => {
    deleteFile(item, path);
  }, [deleteFile, path]);

  const handleRestore = useCallback((item) => {
    restoreFile(item, path);
  }, [restoreFile, path]);

  // Context menu handlers
  const handleContextMenu = useCallback((e, item) => {
    e.preventDefault();
    setContextMenu({
      isOpen: true,
      position: { x: e.clientX, y: e.clientY },
      item
    });
  }, []);

  const closeContextMenu = useCallback(() => {
    setContextMenu({ isOpen: false, position: { x: 0, y: 0 }, item: null });
  }, []);

  // Rename trigger from context menu (needs to pass to child component)
  const [renameTarget, setRenameTarget] = useState(null);
  const handleRenameFromContext = useCallback((item) => {
    setRenameTarget(item);
  }, []);

  // Handle file move (drag & drop)
  const handleMoveFile = useCallback(async (sourceItem, targetFolder) => {
    const sourcePath = joinPath(path, sourceItem.name);
    const destinationPath = joinPath(path, targetFolder.name, sourceItem.name);
    return await moveFile(sourcePath, destinationPath, path);
  }, [path, moveFile]);

  const handleCreateFolder = useCallback(async (folderName) => {
    return await createFolder(folderName, path);
  }, [createFolder, path]);

  const handleUploadClick = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const handleFileInputChange = useCallback((e) => {
    const selectedFiles = e.target.files;
    if (selectedFiles && selectedFiles.length > 0) {
      uploadFiles(Array.from(selectedFiles), path);
    }
    e.target.value = '';
  }, [uploadFiles, path]);

  const handleToggleTrash = useCallback(() => {
    setShowTrash(!showTrash);
    if (!showTrash) loadTrash();
  }, [showTrash, loadTrash]);

  // Batch action handlers
  const handleBatchDownload = useCallback(() => {
    const names = selectedItems.map(f => f.name);
    batchDownload(names, path);
  }, [selectedItems, batchDownload, path]);

  const handleBatchDelete = useCallback(() => {
    const names = selectedItems.map(f => f.name);
    batchDelete(names, path);
    clearSelection();
  }, [selectedItems, batchDelete, path, clearSelection]);

  const handleToggleSelectAll = useCallback(() => {
    const allNames = filteredFiles.map(f => f.name);
    toggleSelectAll(allNames);
  }, [filteredFiles, toggleSelectAll]);

  const breadcrumbs = getBreadcrumbs(path);

  return (
    <div className="space-y-6">
      {/* Error Display */}
      {error && (
        <div className="rounded-xl border border-rose-500/30 bg-rose-500/10 p-4">
          <p className="text-rose-400 text-sm font-medium">{error}</p>
        </div>
      )}

      {/* Hidden File Input */}
      <input
        ref={fileInputRef}
        type="file"
        multiple
        onChange={handleFileInputChange}
        className="hidden"
      />

      {/* Vault Modal */}
      <VaultModal
        isOpen={showVaultModal}
        onClose={() => {
          // If closed without unlocking, redirect to root?
          // setShowVaultModal(false);
          // navigateTo('/');
          // For now just close, but content remains hidden/locked via logic if I implement it
          // Actually, if we don't unlock, we should probably not show list?
          // The file list of /vault contains encrypted blobs. It's safe to show filenames if they are random.
          // But usually we want to obscure them.
          // Let's keep it simple: modal blocking interaction.
          setShowVaultModal(false);
        }}
        onUnlock={handleVaultUnlock}
        onSetup={handleVaultSetup}
      />

      {/* New Folder Modal */}
      <NewFolderModal
        isOpen={showNewFolderModal}
        onClose={() => setShowNewFolderModal(false)}
        onCreateFolder={handleCreateFolder}
        currentPath={path}
      />

      {/* File Preview Modal */}
      <FilePreviewModal
        previewItem={previewItem}
        previewContent={previewContent}
        previewLoading={previewLoading}
        onClose={closePreview}
        onDownload={handleDownload}
      />

      {/* Context Menu */}
      <ContextMenu
        isOpen={contextMenu.isOpen}
        position={contextMenu.position}
        item={contextMenu.item}
        onClose={closeContextMenu}
        onOpen={handleNavigate}
        onPreview={handlePreview}
        onRename={handleRenameFromContext}
        onDownload={handleDownload}
        onDelete={handleDelete}
      />

      {/* Main Content */}
      {showTrash ? (
        <TrashView
          trashedFiles={trashedFiles}
          onRefresh={loadTrash}
          onRestore={handleRestore}
          onEmptyTrash={async () => {
            if (!window.confirm('Papierkorb endgültig leeren? Diese Aktion kann nicht rückgängig gemacht werden.')) return;
            for (const item of trashedFiles) {
              try {
                await fetch(`${window.location.origin}/api/v1/storage/trash/delete?path=${encodeURIComponent(item.name)}`, {
                  method: 'DELETE',
                  credentials: 'include',
                  headers: {
                    'Authorization': `Bearer ${localStorage.getItem('accessToken')}`,
                    'X-CSRF-Token': localStorage.getItem('csrfToken') || '',
                  },
                });
              } catch (err) {
                console.error('Failed to delete:', item.name, err);
              }
            }
            loadTrash();
          }}
        />
      ) : (
        <div className="relative" {...dragProps}>
          {/* Drag & Drop Overlay */}
          <DragDropOverlay isDragging={isDragging} />

          <GlassCard className="!p-0">
            {/* Toolbar */}
            <FileToolbar
              breadcrumbs={breadcrumbs}
              fileCount={filteredFiles.length}
              trashedCount={trashedFiles.length}
              viewMode={viewMode}
              showTrash={showTrash}
              uploading={uploading}
              onNavigate={navigateTo}
              onUploadClick={handleUploadClick}
              onNewFolder={() => setShowNewFolderModal(true)}
              onRefresh={() => loadFiles(path)}
              onViewModeChange={setViewMode}
              onToggleTrash={handleToggleTrash}
            />

            {/* Search & Selection Bar */}
            <div className="px-4 py-3 border-b border-white/5 flex flex-wrap items-center gap-3">
              {/* Search Input */}
              <div className="relative flex-1 min-w-[200px] max-w-md">
                <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                <input
                  type="text"
                  placeholder="Dateien filtern..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full pl-10 pr-10 py-2 bg-slate-800/50 border border-white/10 rounded-lg text-white text-sm placeholder-slate-500 focus:outline-none focus:border-blue-500/50 transition-all"
                />
                {searchQuery && (
                  <button
                    onClick={() => setSearchQuery('')}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                  >
                    <X size={14} />
                  </button>
                )}
              </div>

              {/* Select All Toggle */}
              <button
                onClick={handleToggleSelectAll}
                className={`flex items-center gap-2 px-3 py-2 rounded-lg border transition-all text-sm ${allSelected
                  ? 'bg-blue-500/20 text-blue-400 border-blue-500/30'
                  : 'bg-white/5 text-slate-400 border-white/10 hover:bg-white/10 hover:text-white'
                  }`}
              >
                {allSelected ? <CheckSquare size={16} /> : <Square size={16} />}
                <span className="hidden sm:inline">Alle auswählen</span>
              </button>

              {/* Selection Actions */}
              {selectedCount > 0 && (
                <div className="flex items-center gap-2 ml-auto">
                  <span className="text-sm text-slate-400">
                    {selectedCount} ausgewählt
                  </span>
                  <button
                    onClick={handleBatchDownload}
                    className="flex items-center gap-2 px-3 py-2 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all text-sm"
                    title="Als ZIP herunterladen"
                  >
                    <Download size={16} />
                    <span className="hidden sm:inline">ZIP</span>
                  </button>
                  <button
                    onClick={handleBatchDelete}
                    className="flex items-center gap-2 px-3 py-2 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 transition-all text-sm"
                    title="Ausgewählte löschen"
                  >
                    <Trash2 size={16} />
                    <span className="hidden sm:inline">Löschen</span>
                  </button>
                  <button
                    onClick={clearSelection}
                    className="p-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white border border-white/10 transition-all"
                    title="Auswahl aufheben"
                  >
                    <X size={16} />
                  </button>
                </div>
              )}
            </div>

            {/* Files Content */}
            <div className="p-6">
              {loading ? (
                <div className="flex flex-col items-center justify-center py-12">
                  <Loader2 size={32} className="text-blue-400 animate-spin mb-3" />
                  <p className="text-slate-400 text-sm">Loading files...</p>
                </div>
              ) : viewMode === 'grid' ? (
                <FileGridView
                  files={filteredFiles}
                  onNavigate={handleNavigate}
                  onPreview={handlePreview}
                  onRename={handleRename}
                  onDownload={handleDownload}
                  onDelete={handleDelete}
                  selectedItems={selectedItems}
                  onToggleSelect={toggleSelect}
                  isSelected={isSelected}
                  onContextMenu={handleContextMenu}
                  renameTarget={renameTarget}
                  onRenameComplete={() => setRenameTarget(null)}
                  onMoveFile={handleMoveFile}
                />
              ) : (
                <FileListView
                  files={filteredFiles}
                  onNavigate={handleNavigate}
                  onPreview={handlePreview}
                  onRename={handleRename}
                  onDownload={handleDownload}
                  onDelete={handleDelete}
                  selectedItems={selectedItems}
                  onToggleSelect={toggleSelect}
                  isSelected={isSelected}
                  onContextMenu={handleContextMenu}
                  renameTarget={renameTarget}
                  onRenameComplete={() => setRenameTarget(null)}
                  onMoveFile={handleMoveFile}
                />
              )}
            </div>
          </GlassCard>
        </div>
      )}
    </div>
  );
}
