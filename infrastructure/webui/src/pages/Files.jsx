// Files page - Refactored with extracted components and hooks

import { useEffect, useState, useRef, useCallback } from 'react';
import { Loader2 } from 'lucide-react';

// Hooks
import { useFileStorage } from '../hooks/useFileStorage';
import { useFilePreview } from '../hooks/useFilePreview';
import { useDragAndDrop } from '../hooks/useDragAndDrop';

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

export default function Files() {
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
  } = useFileStorage();

  // File preview hook
  const {
    previewItem,
    previewContent,
    previewLoading,
    openPreview,
    closePreview,
  } = useFilePreview();

  // View state
  const [viewMode, setViewMode] = useState('list');
  const [showTrash, setShowTrash] = useState(false);
  const [showNewFolderModal, setShowNewFolderModal] = useState(false);

  // File input ref for upload button
  const fileInputRef = useRef(null);

  // Drag and drop
  const handleFilesDropped = useCallback((droppedFiles) => {
    uploadFiles(droppedFiles, path);
  }, [uploadFiles, path]);

  const { isDragging, dragProps } = useDragAndDrop(handleFilesDropped);

  // Load files on mount
  useEffect(() => {
    loadFiles('/');
    loadTrash();
  }, []);

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
    downloadFile(item, path);
  }, [downloadFile, path]);

  const handleDelete = useCallback((item) => {
    deleteFile(item, path);
  }, [deleteFile, path]);

  const handleRestore = useCallback((item) => {
    restoreFile(item, path);
  }, [restoreFile, path]);

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

      {/* Main Content */}
      {showTrash ? (
        <TrashView
          trashedFiles={trashedFiles}
          onRefresh={loadTrash}
          onRestore={handleRestore}
        />
      ) : (
        <div className="relative" {...dragProps}>
          {/* Drag & Drop Overlay */}
          <DragDropOverlay isDragging={isDragging} />

          <GlassCard className="!p-0">
            {/* Toolbar */}
            <FileToolbar
              breadcrumbs={breadcrumbs}
              fileCount={files.length}
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

            {/* Files Content */}
            <div className="p-6">
              {loading ? (
                <div className="flex flex-col items-center justify-center py-12">
                  <Loader2 size={32} className="text-blue-400 animate-spin mb-3" />
                  <p className="text-slate-400 text-sm">Loading files...</p>
                </div>
              ) : viewMode === 'grid' ? (
                <FileGridView
                  files={files}
                  onNavigate={handleNavigate}
                  onPreview={handlePreview}
                  onRename={handleRename}
                  onDownload={handleDownload}
                  onDelete={handleDelete}
                />
              ) : (
                <FileListView
                  files={files}
                  onNavigate={handleNavigate}
                  onPreview={handlePreview}
                  onRename={handleRename}
                  onDownload={handleDownload}
                  onDelete={handleDelete}
                />
              )}
            </div>
          </GlassCard>
        </div>
      )}
    </div>
  );
}
