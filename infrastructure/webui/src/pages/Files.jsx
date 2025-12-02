import { useEffect, useState, useRef } from "react";
import { authHeaders } from "../utils/auth";
import {
  FolderOpen,
  File,
  Upload,
  Download,
  Trash2,
  ArrowUp,
  Loader2,
  Grid3x3,
  List,
  Image as ImageIcon,
  FileText,
  FileArchive,
  FileCode,
  FileVideo,
  FileAudio,
  FilePlus,
  Edit3,
  Check,
  X,
  Eye,
  RefreshCw,
  Trash,
  UploadCloud,
  FolderPlus,
  Home,
  ChevronRight,
} from "lucide-react";

const API_BASE =
  import.meta.env.VITE_API_BASE_URL ||
  window.location.origin;

function joinPath(base, name) {
  if (!base || base === "/") {
    return `/${name}`;
  }
  return `${base.replace(/\/+$/, "")}/${name}`;
}

// Glass Card Component
const GlassCard = ({ children, className = "" }) => (
  <div className={`relative overflow-hidden rounded-2xl border border-white/10 bg-slate-900/40 backdrop-blur-xl shadow-2xl ${className}`}>
    <div className="absolute top-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-white/20 to-transparent opacity-50"></div>
    <div className="p-6 h-full flex flex-col">
      {children}
    </div>
  </div>
);

// Get file icon based on extension
const getFileIcon = (name, isDir) => {
  if (isDir) return FolderOpen;

  const ext = name.split('.').pop()?.toLowerCase();

  // Images
  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp'].includes(ext)) {
    return ImageIcon;
  }
  // Text
  if (['txt', 'md', 'log', 'json', 'xml', 'yaml', 'yml', 'csv'].includes(ext)) {
    return FileText;
  }
  // Code
  if (['js', 'jsx', 'ts', 'tsx', 'py', 'go', 'rs', 'java', 'cpp', 'c', 'h', 'css', 'html', 'sh'].includes(ext)) {
    return FileCode;
  }
  // Archive
  if (['zip', 'tar', 'gz', 'rar', '7z', 'bz2'].includes(ext)) {
    return FileArchive;
  }
  // Video
  if (['mp4', 'avi', 'mkv', 'mov', 'wmv', 'flv', 'webm'].includes(ext)) {
    return FileVideo;
  }
  // Audio
  if (['mp3', 'wav', 'flac', 'ogg', 'aac', 'm4a'].includes(ext)) {
    return FileAudio;
  }

  return File;
};

// Check if file is image
const isImage = (name) => {
  const ext = name.split('.').pop()?.toLowerCase();
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp'].includes(ext);
};

// Check if file is text
const isText = (name) => {
  const ext = name.split('.').pop()?.toLowerCase();
  return ['txt', 'md', 'log', 'json', 'xml', 'yaml', 'yml', 'csv', 'js', 'jsx', 'ts', 'tsx', 'py', 'go', 'rs', 'java', 'cpp', 'c', 'h', 'css', 'html', 'sh'].includes(ext);
};

export default function Files() {
  const [path, setPath] = useState("/");
  const [files, setFiles] = useState([]);
  const [trashedFiles, setTrashedFiles] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [selectedFile, setSelectedFile] = useState(null);

  // View modes
  const [viewMode, setViewMode] = useState("list"); // "list" or "grid"
  const [showTrash, setShowTrash] = useState(false);

  // Rename state
  const [renamingItem, setRenamingItem] = useState(null);
  const [newName, setNewName] = useState("");

  // Preview state
  const [previewItem, setPreviewItem] = useState(null);
  const [previewContent, setPreviewContent] = useState(null);
  const [previewLoading, setPreviewLoading] = useState(false);

  // Drag & Drop state
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef(null);

  // New Folder modal state
  const [showNewFolderModal, setShowNewFolderModal] = useState(false);
  const [newFolderName, setNewFolderName] = useState("");
  const [creatingFolder, setCreatingFolder] = useState(false);

  const loadFiles = async (target = path) => {
    if (!authHeaders().Authorization) {
      setError("Bitte zuerst einloggen.");
      return;
    }
    setLoading(true);
    setError("");
    try {
      const res = await fetch(
        `${API_BASE}/api/v1/storage/files?path=${encodeURIComponent(target)}`,
        {
          credentials: "include",
          headers: authHeaders()
        }
      );
      if (res.status === 401) {
        setError("Authentifizierung fehlgeschlagen. Bitte neu einloggen.");
        return;
      }
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}: Failed to load files`);
      }
      const data = await res.json();

      // Filter out .trash directory from the list
      const filteredFiles = (data.items || []).filter(item => item.name !== ".trash");
      setFiles(filteredFiles);
    } catch (err) {
      setError(err.message || "Unknown error");
    } finally {
      setLoading(false);
    }
  };

  const loadTrash = async () => {
    try {
      const res = await fetch(
        `${API_BASE}/api/v1/storage/trash`,
        {
          credentials: "include",
          headers: authHeaders()
        }
      );
      if (res.ok) {
        const data = await res.json();
        setTrashedFiles(data.items || []);
      }
    } catch (err) {
      console.error("Failed to load trash:", err);
    }
  };

  useEffect(() => {
    loadFiles("/");
    loadTrash();
  }, []);

  const handleUpload = async (filesToUpload = null) => {
    // Use provided files or fall back to selectedFile
    const files = filesToUpload || (selectedFile ? [selectedFile] : null);
    if (!files || files.length === 0) return;

    setUploading(true);
    setError("");

    try {
      // Upload each file sequentially
      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        console.log(`Uploading file ${i + 1}/${files.length}:`, file.name, `Size: ${file.size} bytes`);

        const form = new FormData();
        form.append("file", file);
        form.append("path", path);

        // WICHTIG: Für FormData NICHT Content-Type Header setzen!
        // Browser setzt automatisch multipart/form-data mit boundary
        const headers = authHeaders();
        delete headers['Content-Type']; // Falls vorhanden, entfernen

        // SECURITY FIX: Explizit CSRF Token aus localStorage holen und setzen
        // Verlasse dich nicht nur auf authHeaders() - erzwinge Token-Präsenz
        const csrfToken = localStorage.getItem('csrfToken') || localStorage.getItem('csrf_token');
        if (csrfToken) {
          headers['X-CSRF-Token'] = csrfToken;
        } else {
          console.warn('⚠️ WARNING: No CSRF token found in localStorage');
        }

        const uploadUrl = `${API_BASE}/api/v1/storage/upload`;
        console.log('Upload URL:', uploadUrl);
        console.log('Upload Headers prepared:', {
          ...headers,
          Authorization: headers.Authorization ? 'REDACTED (present)' : 'MISSING',
          'X-CSRF-Token': headers['X-CSRF-Token'] ? 'PRESENT' : 'MISSING'
        });

        const res = await fetch(uploadUrl, {
          method: "POST",
          body: form,
          credentials: "include",
          headers: headers,
        });

        console.log('Upload response status:', res.status);

        if (res.status === 401) {
          setError("Authentifizierung fehlgeschlagen.");
          return;
        }
        if (!res.ok) {
          const errorText = await res.text().catch(() => 'No error details');
          console.error('Upload failed:', res.status, errorText);
          throw new Error(`Upload failed for ${file.name}: HTTP ${res.status} - ${errorText}`);
        }

        console.log(`✓ Successfully uploaded: ${file.name}`);
      }

      setSelectedFile(null);
      await loadFiles(path);
    } catch (err) {
      console.error('Upload error details:', err);
      setError(err.message || "Unknown error");
    } finally {
      setUploading(false);
    }
  };

  const handleDownload = async (item) => {
    const target = joinPath(path, item.name);
    try {
      const res = await fetch(
        `${API_BASE}/api/v1/storage/download?path=${encodeURIComponent(target)}`,
        {
          credentials: "include",
          headers: authHeaders(),
        }
      );
      if (!res.ok) {
        setError(`Download failed: HTTP ${res.status}`);
        return;
      }
      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = item.name;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setError(err.message);
    }
  };

  const handleDelete = async (item) => {
    // Prevent deletion of .trash directory
    if (item.name === ".trash") {
      setError("Der Papierkorb kann nicht gelöscht werden");
      return;
    }

    if (!window.confirm(`"${item.name}" wirklich löschen?`)) return;

    const target = joinPath(path, item.name);
    try {
      const res = await fetch(
        `${API_BASE}/api/v1/storage/delete?path=${encodeURIComponent(target)}`,
        {
          method: "DELETE",
          credentials: "include",
          headers: authHeaders(),
        }
      );
      if (!res.ok) {
        setError(`Delete failed: HTTP ${res.status}`);
        return;
      }
      await loadFiles(path);
      await loadTrash();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleRestore = async (item) => {
    try {
      const res = await fetch(
        `${API_BASE}/api/v1/storage/restore?path=${encodeURIComponent(item.originalPath || item.name)}`,
        {
          method: "POST",
          credentials: "include",
          headers: authHeaders(),
        }
      );
      if (!res.ok) {
        setError(`Restore failed: HTTP ${res.status}`);
        return;
      }
      await loadTrash();
      await loadFiles(path);
    } catch (err) {
      setError(err.message);
    }
  };

  const handleRename = async () => {
    if (!renamingItem || !newName || newName === renamingItem.name) {
      setRenamingItem(null);
      setNewName("");
      return;
    }

    const oldPath = joinPath(path, renamingItem.name);
    const newPath = joinPath(path, newName);

    try {
      const res = await fetch(
        `${API_BASE}/api/v1/storage/rename`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            ...authHeaders(),
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            oldPath,
            newPath,
          }),
        }
      );
      if (!res.ok) {
        throw new Error(`Rename failed: HTTP ${res.status}`);
      }
      setRenamingItem(null);
      setNewName("");
      await loadFiles(path);
    } catch (err) {
      setError(err.message);
    }
  };

  const handlePreview = async (item) => {
    if (item.isDir) return;

    setPreviewItem(item);
    setPreviewLoading(true);
    setPreviewContent(null);

    const target = joinPath(path, item.name);

    try {
      if (isImage(item.name)) {
        // For images, create blob URL
        const res = await fetch(
          `${API_BASE}/api/v1/storage/download?path=${encodeURIComponent(target)}`,
          {
            credentials: "include",
            headers: authHeaders(),
          }
        );
        if (res.ok) {
          const blob = await res.blob();
          const url = window.URL.createObjectURL(blob);
          setPreviewContent({ type: 'image', url });
        }
      } else if (isText(item.name)) {
        // For text, fetch content
        const res = await fetch(
          `${API_BASE}/api/v1/storage/download?path=${encodeURIComponent(target)}`,
          {
            credentials: "include",
            headers: authHeaders(),
          }
        );
        if (res.ok) {
          const text = await res.text();
          setPreviewContent({ type: 'text', content: text });
        }
      }
    } catch (err) {
      setError("Preview failed: " + err.message);
    } finally {
      setPreviewLoading(false);
    }
  };

  const closePreview = () => {
    if (previewContent?.type === 'image' && previewContent.url) {
      window.URL.revokeObjectURL(previewContent.url);
    }
    setPreviewItem(null);
    setPreviewContent(null);
  };

  const handleNavigate = (item) => {
    if (!item.isDir) return;
    const nextPath = joinPath(path, item.name);
    setPath(nextPath);
    loadFiles(nextPath);
  };

  const goUp = () => {
    if (path === "/") return;
    const parts = path.split("/").filter(Boolean);
    parts.pop();
    const parent = parts.length ? `/${parts.join("/")}` : "/";
    setPath(parent);
    loadFiles(parent);
  };

  const formatFileSize = (bytes) => {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Drag & Drop handlers
  const handleDragOver = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const droppedFiles = e.dataTransfer.files;
    if (droppedFiles && droppedFiles.length > 0) {
      // Auto-upload all dropped files (fire and forget)
      handleUpload(Array.from(droppedFiles)).catch(err => {
        console.error('Upload error:', err);
      });
    }
  };

  // New Folder creation
  const handleCreateFolder = async () => {
    if (!newFolderName || newFolderName.trim() === '') {
      setError('Ordnername darf nicht leer sein');
      return;
    }

    setCreatingFolder(true);
    setError("");

    try {
      const folderPath = joinPath(path, newFolderName.trim());
      const res = await fetch(`${API_BASE}/api/v1/storage/mkdir`, {
        method: "POST",
        credentials: "include",
        headers: {
          ...authHeaders(),
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ path: folderPath }),
      });

      if (!res.ok) {
        throw new Error(`Ordner erstellen fehlgeschlagen: HTTP ${res.status}`);
      }

      setShowNewFolderModal(false);
      setNewFolderName("");
      await loadFiles(path);
    } catch (err) {
      setError(err.message || "Fehler beim Erstellen des Ordners");
    } finally {
      setCreatingFolder(false);
    }
  };

  // Breadcrumb navigation
  const getBreadcrumbs = () => {
    if (path === '/') return [{ name: 'Home', path: '/' }];

    const parts = path.split('/').filter(Boolean);
    const breadcrumbs = [{ name: 'Home', path: '/' }];

    let currentPath = '';
    parts.forEach((part) => {
      currentPath += `/${part}`;
      breadcrumbs.push({ name: part, path: currentPath });
    });

    return breadcrumbs;
  };

  const navigateToBreadcrumb = (breadcrumbPath) => {
    setPath(breadcrumbPath);
    loadFiles(breadcrumbPath);
  };

  // File/Folder Card Component for Grid View
  const FileCard = ({ item }) => {
    const Icon = getFileIcon(item.name, item.isDir);
    const isRenaming = renamingItem?.name === item.name;

    return (
      <div
        className="group relative overflow-hidden rounded-xl border border-white/10 bg-slate-900/40 backdrop-blur-xl hover:bg-white/5 transition-all cursor-pointer"
        onClick={() => item.isDir && handleNavigate(item)}
      >
        <div className="p-4 flex flex-col items-center text-center">
          {/* Icon */}
          <div className={`p-4 rounded-xl mb-3 ${item.isDir ? 'bg-blue-500/20 text-blue-400' : 'bg-slate-800/50 text-slate-400'} group-hover:scale-110 transition-transform`}>
            <Icon size={32} />
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
                  if (e.key === 'Escape') {
                    setRenamingItem(null);
                    setNewName("");
                  }
                }}
                className="flex-1 px-2 py-1 text-xs bg-slate-800 border border-white/10 rounded text-white focus:outline-none focus:border-blue-500"
                autoFocus
              />
              <button onClick={handleRename} className="p-1 rounded bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30">
                <Check size={14} />
              </button>
              <button onClick={() => { setRenamingItem(null); setNewName(""); }} className="p-1 rounded bg-rose-500/20 text-rose-400 hover:bg-rose-500/30">
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
            {!item.isDir && (isImage(item.name) || isText(item.name)) && (
              <button
                onClick={() => handlePreview(item)}
                className="p-1.5 rounded-lg bg-violet-500/10 hover:bg-violet-500/20 text-violet-400 border border-violet-500/20 transition-all"
                title="Preview"
              >
                <Eye size={12} />
              </button>
            )}
            <button
              onClick={() => { setRenamingItem(item); setNewName(item.name); }}
              className="p-1.5 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 border border-blue-500/20 transition-all"
              title="Rename"
            >
              <Edit3 size={12} />
            </button>
            {!item.isDir && (
              <button
                onClick={() => handleDownload(item)}
                className="p-1.5 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all"
                title="Download"
              >
                <Download size={12} />
              </button>
            )}
            <button
              onClick={() => handleDelete(item)}
              className="p-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 transition-all"
              title="Delete"
            >
              <Trash2 size={12} />
            </button>
          </div>
        </div>
      </div>
    );
  };

  const breadcrumbs = getBreadcrumbs();

  return (
    <div className="space-y-6">

      {/* Error Display */}
      {error && (
        <div className="rounded-xl border border-rose-500/30 bg-rose-500/10 p-4">
          <p className="text-rose-400 text-sm font-medium">{error}</p>
        </div>
      )}

      {/* Hidden File Input for Upload */}
      <input
        ref={fileInputRef}
        type="file"
        multiple
        onChange={(e) => {
          const files = e.target.files;
          if (files && files.length > 0) {
            handleUpload(Array.from(files));
          }
          // Reset input so the same file can be selected again
          e.target.value = '';
        }}
        className="hidden"
      />

      {/* New Folder Modal */}
      {showNewFolderModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-200">
          <div className="w-full max-w-md animate-in zoom-in-95 duration-200">
            <GlassCard>
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="p-3 rounded-xl bg-blue-500/20 border border-blue-500/30">
                    <FolderPlus size={20} className="text-blue-400" />
                  </div>
                  <h2 className="text-xl font-bold text-white">Neuer Ordner</h2>
                </div>
                <button
                  onClick={() => { setShowNewFolderModal(false); setNewFolderName(""); }}
                  className="p-2 rounded-lg bg-slate-800/50 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 border border-white/10 hover:border-rose-500/30 transition-all"
                >
                  <X size={18} />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Ordnername
                  </label>
                  <input
                    type="text"
                    value={newFolderName}
                    onChange={(e) => setNewFolderName(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' && newFolderName) handleCreateFolder();
                      if (e.key === 'Escape') { setShowNewFolderModal(false); setNewFolderName(""); }
                    }}
                    placeholder="Mein Ordner"
                    className="w-full px-4 py-2.5 bg-slate-800/50 border border-white/10 rounded-lg text-white focus:border-blue-500/50 focus:bg-slate-800 focus:outline-none transition-all"
                    autoFocus
                  />
                  <p className="text-xs text-slate-500 mt-1.5">
                    Wird erstellt in: <span className="text-blue-400 font-mono">{path}</span>
                  </p>
                </div>

                <div className="flex items-center justify-end gap-3 pt-4 border-t border-white/5">
                  <button
                    onClick={() => { setShowNewFolderModal(false); setNewFolderName(""); }}
                    className="px-4 py-2 rounded-lg bg-slate-800/50 hover:bg-slate-800 text-slate-300 hover:text-white border border-white/10 transition-all"
                  >
                    Abbrechen
                  </button>
                  <button
                    onClick={handleCreateFolder}
                    disabled={creatingFolder || !newFolderName}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-blue-500/20 hover:bg-blue-500/30 text-blue-400 border border-blue-500/30 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-[0_0_15px_rgba(59,130,246,0.2)]"
                  >
                    {creatingFolder ? (
                      <>
                        <Loader2 size={16} className="animate-spin" />
                        <span>Erstelle...</span>
                      </>
                    ) : (
                      <>
                        <FolderPlus size={16} />
                        <span>Erstellen</span>
                      </>
                    )}
                  </button>
                </div>
              </div>
            </GlassCard>
          </div>
        </div>
      )}

      {/* Files Display */}
      {showTrash ? (
        /* Trash Bin View */
        <GlassCard>
          <div className="flex items-center justify-between mb-6">
            <div>
              <h3 className="text-white font-semibold text-lg tracking-tight flex items-center gap-2">
                <Trash className="text-rose-400" size={20} />
                Papierkorb
              </h3>
              <p className="text-slate-400 text-xs mt-1">{trashedFiles.length} gelöschte Dateien</p>
            </div>
            <button
              onClick={() => loadTrash()}
              className="p-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white transition-all"
              title="Refresh"
            >
              <RefreshCw size={18} />
            </button>
          </div>

          {trashedFiles.length === 0 ? (
            <div className="py-12 text-center text-slate-400">
              <Trash size={48} className="mx-auto mb-3 opacity-30" />
              <p className="text-sm">Papierkorb ist leer</p>
            </div>
          ) : (
            <div className="space-y-2">
              {trashedFiles.map((item, idx) => {
                const Icon = getFileIcon(item.name, item.isDir);
                return (
                  <div
                    key={idx}
                    className="flex items-center justify-between p-3 rounded-lg bg-white/5 border border-white/5 hover:bg-white/10 transition-all"
                  >
                    <div className="flex items-center gap-3">
                      <div className="p-2 rounded-lg bg-slate-800/50 text-slate-400">
                        <Icon size={16} />
                      </div>
                      <div>
                        <p className="text-white font-medium text-sm">{item.name}</p>
                        <p className="text-slate-500 text-xs">{item.originalPath || 'Unknown path'}</p>
                      </div>
                    </div>
                    <button
                      onClick={() => handleRestore(item)}
                      className="flex items-center gap-2 px-4 py-2 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all"
                    >
                      <RefreshCw size={14} />
                      <span className="text-sm font-medium">Wiederherstellen</span>
                    </button>
                  </div>
                );
              })}
            </div>
          )}
        </GlassCard>
      ) : (
        /* Files View with Unified Toolbar */
        <div
          className="relative"
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          {/* Drag & Drop Overlay */}
          {isDragging && (
            <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/70 backdrop-blur-md pointer-events-none animate-in fade-in duration-200">
              <div className="p-16 rounded-3xl border-4 border-dashed border-blue-400 bg-gradient-to-br from-blue-500/20 to-cyan-500/20 shadow-[0_0_60px_rgba(59,130,246,0.4)] transform scale-105 transition-transform">
                <UploadCloud size={80} className="text-blue-400 mx-auto mb-6 animate-bounce" />
                <p className="text-3xl font-bold text-blue-400 mb-2">Dateien hier ablegen</p>
                <p className="text-slate-200 text-base">Mehrere Dateien gleichzeitig möglich</p>
                <div className="mt-6 flex items-center justify-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-blue-400 animate-pulse"></div>
                  <div className="w-3 h-3 rounded-full bg-blue-400 animate-pulse delay-75"></div>
                  <div className="w-3 h-3 rounded-full bg-blue-400 animate-pulse delay-150"></div>
                </div>
              </div>
            </div>
          )}

          <GlassCard className="!p-0">
            {/* Unified Toolbar */}
            <div className="p-4 border-b border-white/5">
              <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                {/* Left: Breadcrumbs */}
                <div className="flex items-center gap-2 flex-wrap">
                  {breadcrumbs.map((crumb, index) => (
                    <div key={crumb.path} className="flex items-center gap-2">
                      {index === 0 ? (
                        <button
                          onClick={() => navigateToBreadcrumb(crumb.path)}
                          className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300 hover:text-white border border-white/10 transition-all"
                        >
                          <Home size={16} />
                          <span className="text-sm font-medium">{crumb.name}</span>
                        </button>
                      ) : (
                        <>
                          <ChevronRight size={16} className="text-slate-600" />
                          <button
                            onClick={() => navigateToBreadcrumb(crumb.path)}
                            className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
                              index === breadcrumbs.length - 1
                                ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
                                : 'bg-white/5 hover:bg-white/10 text-slate-300 hover:text-white border border-white/10'
                            }`}
                          >
                            {crumb.name}
                          </button>
                        </>
                      )}
                    </div>
                  ))}
                  <div className="ml-2 px-2 py-1 rounded-lg bg-slate-800/50 border border-white/5">
                    <span className="text-xs text-slate-400">{files.length} items</span>
                  </div>
                </div>

                {/* Right: Actions */}
                <div className="flex items-center gap-2 flex-wrap">
                  {/* Upload Button */}
                  <button
                    onClick={() => fileInputRef.current?.click()}
                    disabled={uploading}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-blue-500/20 hover:bg-blue-500/30 text-blue-400 border border-blue-500/30 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-[0_0_15px_rgba(59,130,246,0.2)]"
                    title="Upload Files (Multiple)"
                  >
                    {uploading ? (
                      <>
                        <Loader2 size={16} className="animate-spin" />
                        <span className="hidden sm:inline text-sm font-medium">Hochladen...</span>
                      </>
                    ) : (
                      <>
                        <UploadCloud size={16} />
                        <span className="hidden sm:inline text-sm font-medium">Hochladen</span>
                      </>
                    )}
                  </button>

                  {/* New Folder Button */}
                  <button
                    onClick={() => setShowNewFolderModal(true)}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all"
                    title="New Folder"
                  >
                    <FolderPlus size={16} />
                    <span className="hidden sm:inline text-sm font-medium">New Folder</span>
                  </button>

                  {/* Refresh Button */}
                  <button
                    onClick={() => loadFiles(path)}
                    className="p-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white border border-white/10 transition-all"
                    title="Refresh"
                  >
                    <RefreshCw size={16} />
                  </button>

                  {/* View Mode Toggle */}
                  <div className="flex items-center gap-1 p-1 rounded-lg bg-white/5 border border-white/10">
                    <button
                      onClick={() => setViewMode("list")}
                      className={`p-2 rounded ${viewMode === "list" ? 'bg-blue-500/20 text-blue-400' : 'text-slate-400 hover:text-white'} transition-all`}
                      title="List View"
                    >
                      <List size={16} />
                    </button>
                    <button
                      onClick={() => setViewMode("grid")}
                      className={`p-2 rounded ${viewMode === "grid" ? 'bg-blue-500/20 text-blue-400' : 'text-slate-400 hover:text-white'} transition-all`}
                      title="Grid View"
                    >
                      <Grid3x3 size={16} />
                    </button>
                  </div>

                  {/* Trash Button */}
                  <button
                    onClick={() => { setShowTrash(!showTrash); if (!showTrash) loadTrash(); }}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg ${showTrash ? 'bg-rose-500/20 text-rose-400 border-rose-500/30' : 'bg-white/5 text-slate-300 border-white/10'} hover:bg-white/10 hover:text-white border transition-all`}
                  >
                    <Trash size={16} />
                    <span className="hidden sm:inline text-sm font-medium">Trash</span>
                    {trashedFiles.length > 0 && (
                      <span className="px-2 py-0.5 rounded-full bg-rose-500 text-white text-xs font-bold">
                        {trashedFiles.length}
                      </span>
                    )}
                  </button>
                </div>
              </div>
            </div>

            {/* Files Content */}
            <div className="p-6">

          {loading ? (
            <div className="flex flex-col items-center justify-center py-12">
              <Loader2 size={32} className="text-blue-400 animate-spin mb-3" />
              <p className="text-slate-400 text-sm">Loading files...</p>
            </div>
          ) : viewMode === "grid" ? (
            /* Grid View */
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
              {files.map((item) => (
                <FileCard key={item.name} item={item} />
              ))}
              {files.length === 0 && (
                <div className="col-span-full py-12 text-center text-slate-400">
                  <FolderOpen size={48} className="mx-auto mb-3 opacity-30" />
                  <p className="text-sm">No files or folders</p>
                </div>
              )}
            </div>
          ) : (
            /* List View */
            <div className="overflow-x-auto -mx-6 px-6">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="text-xs text-slate-500 border-b border-white/5">
                    <th className="py-3 px-2 font-medium uppercase tracking-wider">Name</th>
                    <th className="py-3 px-2 font-medium uppercase tracking-wider">Size</th>
                    <th className="py-3 px-2 font-medium uppercase tracking-wider">Modified</th>
                    <th className="py-3 px-2 font-medium uppercase tracking-wider text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="text-sm">
                  {files.map((item) => {
                    const Icon = getFileIcon(item.name, item.isDir);
                    const isRenaming = renamingItem?.name === item.name;

                    return (
                      <tr
                        key={item.name}
                        className="group border-b border-white/5 last:border-0 hover:bg-white/5 transition-colors cursor-pointer"
                        onClick={() => item.isDir && handleNavigate(item)}
                      >
                        <td className="py-4 px-2 font-medium text-white">
                          {isRenaming ? (
                            <div className="flex items-center gap-2" onClick={(e) => e.stopPropagation()}>
                              <div className={`p-2 rounded-lg ${item.isDir ? 'bg-blue-500/20 text-blue-400' : 'bg-slate-800 text-slate-400'}`}>
                                <Icon size={16} />
                              </div>
                              <input
                                type="text"
                                value={newName}
                                onChange={(e) => setNewName(e.target.value)}
                                onKeyDown={(e) => {
                                  if (e.key === 'Enter') handleRename();
                                  if (e.key === 'Escape') {
                                    setRenamingItem(null);
                                    setNewName("");
                                  }
                                }}
                                className="flex-1 px-3 py-1.5 bg-slate-800 border border-white/10 rounded-lg text-white focus:outline-none focus:border-blue-500"
                                autoFocus
                              />
                              <button onClick={handleRename} className="p-2 rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30">
                                <Check size={14} />
                              </button>
                              <button onClick={() => { setRenamingItem(null); setNewName(""); }} className="p-2 rounded-lg bg-rose-500/20 text-rose-400 hover:bg-rose-500/30">
                                <X size={14} />
                              </button>
                            </div>
                          ) : (
                            <div className="flex items-center gap-3">
                              <div className={`p-2 rounded-lg ${item.isDir ? 'bg-blue-500/20 text-blue-400' : 'bg-slate-800 text-slate-400'}`}>
                                <Icon size={16} />
                              </div>
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
                            {!item.isDir && (isImage(item.name) || isText(item.name)) && (
                              <button
                                onClick={() => handlePreview(item)}
                                className="p-2 rounded-lg bg-violet-500/10 hover:bg-violet-500/20 text-violet-400 border border-violet-500/20 transition-all"
                                title="Preview"
                              >
                                <Eye size={14} />
                              </button>
                            )}
                            <button
                              onClick={() => { setRenamingItem(item); setNewName(item.name); }}
                              className="p-2 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 border border-blue-500/20 transition-all"
                              title="Rename"
                            >
                              <Edit3 size={14} />
                            </button>
                            {!item.isDir && (
                              <button
                                onClick={() => handleDownload(item)}
                                className="p-2 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all"
                                title="Download"
                              >
                                <Download size={14} />
                              </button>
                            )}
                            <button
                              onClick={() => handleDelete(item)}
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
                  {files.length === 0 && (
                    <tr>
                      <td colSpan="4" className="py-12 text-center text-slate-400">
                        <FolderOpen size={48} className="mx-auto mb-3 opacity-30" />
                        <p className="text-sm">No files or folders</p>
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
            </div>
          </GlassCard>
        </div>
      )}

      {/* Preview Modal */}
      {previewItem && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-md p-4 animate-in fade-in duration-200">
          <div className="w-full max-w-4xl max-h-[90vh] animate-in zoom-in-95 duration-200">
            <GlassCard>
              <div className="p-6">
                {/* Modal Header */}
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h2 className="text-xl font-bold text-white tracking-tight">{previewItem.name}</h2>
                    <p className="text-slate-400 text-sm mt-1">{formatFileSize(previewItem.size)}</p>
                  </div>
                  <button
                    onClick={closePreview}
                    className="p-2 rounded-lg bg-slate-800/50 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 border border-white/10 hover:border-rose-500/30 transition-all"
                  >
                    <X size={20} />
                  </button>
                </div>

                {/* Preview Content */}
                {previewLoading ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 size={32} className="text-blue-400 animate-spin" />
                  </div>
                ) : previewContent?.type === 'image' ? (
                  <div className="max-h-[60vh] overflow-auto rounded-lg bg-black/50 p-4">
                    <img
                      src={previewContent.url}
                      alt={previewItem.name}
                      className="max-w-full mx-auto rounded"
                    />
                  </div>
                ) : previewContent?.type === 'text' ? (
                  <div className="max-h-[60vh] overflow-auto rounded-lg bg-slate-900 p-4 border border-white/10">
                    <pre className="text-slate-300 text-sm font-mono whitespace-pre-wrap">
                      {previewContent.content}
                    </pre>
                  </div>
                ) : (
                  <div className="py-12 text-center text-slate-400">
                    <p>Preview not available</p>
                  </div>
                )}

                {/* Actions */}
                <div className="flex items-center justify-end gap-3 mt-4 pt-4 border-t border-white/5">
                  <button
                    onClick={() => handleDownload(previewItem)}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 transition-all"
                  >
                    <Download size={16} />
                    <span>Download</span>
                  </button>
                  <button
                    onClick={closePreview}
                    className="px-4 py-2 rounded-lg bg-slate-800/50 hover:bg-slate-800 text-slate-300 hover:text-white border border-white/10 transition-all"
                  >
                    Close
                  </button>
                </div>
              </div>
            </GlassCard>
          </div>
        </div>
      )}
    </div>
  );
}
