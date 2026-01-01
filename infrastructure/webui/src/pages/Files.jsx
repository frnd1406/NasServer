// Files page - Cleaned up and refactored
import { useEffect, useState, useRef, useCallback, useMemo } from 'react';
import { Loader2 } from 'lucide-react';

// Hooks
import { useFileStorage } from '../hooks/useFileStorage';
import { useFilePreview } from '../hooks/useFilePreview';
import { useDragAndDrop } from '../hooks/useDragAndDrop';
import { useFileSelection } from '../hooks/useFileSelection';

// Utils
import { getBreadcrumbs, joinPath } from '../utils/fileUtils';

// Components
import { GlassCard } from '../components/ui/GlassCard';
import { FileHeader } from '../components/Files/FileHeader';
import { FileGridView } from '../components/FileGridView';
import { FileListView } from '../components/FileListView';
import { TrashView } from '../components/TrashView';
import { DragDropOverlay } from '../components/DragDropOverlay';
import { NewFolderModal } from '../components/NewFolderModal';
import { FilePreviewModal } from '../components/FilePreviewModal';
import { ContextMenu } from '../components/ContextMenu';

import { useVault } from '../context/VaultContext';

export default function Files({ initialPath = '/' }) {
  // Vault context
  const { isUnlocked, password } = useVault();

  // Storage Hook
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
    emptyTrash,
    restoreFile,
    renameFile,
    createFolder,
    navigateTo,
    batchDownload,
    downloadFolderAsZip,
    batchDelete,
    moveFile,
  } = useFileStorage(initialPath, isUnlocked ? password : null);

  // Preview Hook
  const {
    previewItem,
    previewContent,
    previewLoading,
    openPreview,
    closePreview,
  } = useFilePreview();

  // Local UI State
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState('list');
  const [showTrash, setShowTrash] = useState(false);
  const [showNewFolderModal, setShowNewFolderModal] = useState(false);
  const [encryptionMode, setEncryptionMode] = useState('auto');
  const [contextMenu, setContextMenu] = useState({ isOpen: false, position: { x: 0, y: 0 }, item: null });
  const [renameTarget, setRenameTarget] = useState(null);

  // File Input Ref
  const fileInputRef = useRef(null);

  // Filter Logic
  const filteredFiles = useMemo(() => {
    if (!searchQuery.trim()) return files;
    const query = searchQuery.toLowerCase();
    return files.filter(f => f.name.toLowerCase().includes(query));
  }, [files, searchQuery]);

  // Selection Hook
  const {
    selectedCount,
    selectedItems,
    allSelected,
    toggleSelect,
    clearSelection,
    toggleSelectAll,
    isSelected,
  } = useFileSelection(filteredFiles);

  // Drag & Drop Hook
  const handleFilesDropped = useCallback((droppedFiles) => {
    uploadFiles(droppedFiles, path);
  }, [uploadFiles, path]);

  const { isDragging, dragProps } = useDragAndDrop(handleFilesDropped);

  // Initial Load
  useEffect(() => {
    loadFiles('/');
    loadTrash();
  }, [loadFiles, loadTrash]);

  // Path Change Logic
  useEffect(() => {
    clearSelection();
    setSearchQuery('');
  }, [path, clearSelection]);

  // --- Handlers ---

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
      downloadFolderAsZip(joinPath(path, item.name));
    } else {
      downloadFile(item, path);
    }
  }, [downloadFile, downloadFolderAsZip, path]);

  const handleDelete = useCallback((item) => {
    deleteFile(item, path);
  }, [deleteFile, path]);

  const handleContextMenu = useCallback((e, item) => {
    e.preventDefault();
    setContextMenu({
      isOpen: true,
      position: { x: e.clientX, y: e.clientY },
      item
    });
  }, []);

  const handleRenameFromContext = useCallback((item) => {
    setRenameTarget(item);
  }, []);

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
      uploadFiles(Array.from(selectedFiles), path, { encryptionMode });
    }
    e.target.value = '';
  }, [uploadFiles, path, encryptionMode]);

  const handleToggleTrash = useCallback(() => {
    setShowTrash(!showTrash);
    if (!showTrash) loadTrash();
  }, [showTrash, loadTrash]);

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

  // --- Render ---

  return (
    <div className="space-y-6">
      {/* Error Banner */}
      {error && (
        <div className="rounded-xl border border-rose-500/30 bg-rose-500/10 p-4">
          <p className="text-rose-400 text-sm font-medium">{error}</p>
        </div>
      )}

      {/* Hidden Upload Input */}
      <input
        ref={fileInputRef}
        type="file"
        multiple
        onChange={handleFileInputChange}
        className="hidden"
      />

      {/* Modals & Overlays */}
      <NewFolderModal
        isOpen={showNewFolderModal}
        onClose={() => setShowNewFolderModal(false)}
        onCreateFolder={handleCreateFolder}
        currentPath={path}
      />

      <FilePreviewModal
        previewItem={previewItem}
        previewContent={previewContent}
        previewLoading={previewLoading}
        onClose={closePreview}
        onDownload={handleDownload}
      />

      <ContextMenu
        isOpen={contextMenu.isOpen}
        position={contextMenu.position}
        item={contextMenu.item}
        onClose={() => setContextMenu({ ...contextMenu, isOpen: false })}
        onOpen={handleNavigate}
        onPreview={handlePreview}
        onRename={handleRenameFromContext}
        onDownload={handleDownload}
        onDelete={handleDelete}
      />

      {/* Main View */}
      {showTrash ? (
        <TrashView
          trashedFiles={trashedFiles}
          onRefresh={loadTrash}
          onRestore={(item) => restoreFile(item, path)}
          onEmptyTrash={emptyTrash}
        />
      ) : (
        <div className="relative" {...dragProps}>
          <DragDropOverlay isDragging={isDragging} />

          <GlassCard className="!p-0">
            <FileHeader
              breadcrumbs={breadcrumbs}
              fileCount={filteredFiles.length}
              trashedCount={trashedFiles.length}
              viewMode={viewMode}
              showTrash={showTrash}
              uploading={uploading}
              encryptionMode={encryptionMode}
              searchQuery={searchQuery}
              selectedItems={selectedItems}
              selectedCount={selectedCount}
              allSelected={allSelected}
              onModeChange={setEncryptionMode}
              onNavigate={navigateTo}
              onUploadClick={handleUploadClick}
              onNewFolder={() => setShowNewFolderModal(true)}
              onRefresh={() => loadFiles(path)}
              onViewModeChange={setViewMode}
              onToggleTrash={handleToggleTrash}
              setSearchQuery={setSearchQuery}
              onToggleSelectAll={handleToggleSelectAll}
              onBatchDownload={handleBatchDownload}
              onBatchDelete={handleBatchDelete}
              onClearSelection={clearSelection}
            />

            <div className="p-6">
              {loading ? (
                <div className="flex flex-col items-center justify-center py-12">
                  <Loader2 size={32} className="text-blue-400 animate-spin mb-3" />
                  <p className="text-slate-400 text-sm">Dateien werden geladen...</p>
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
