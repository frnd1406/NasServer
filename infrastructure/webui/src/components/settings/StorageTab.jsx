import { useState, useEffect } from "react";
import {
    HardDrive,
    Trash2,
    Folder,
    File,
    AlertTriangle,
    Loader2,
    RefreshCw,
    Plus,
    X,
    Database,
    Save
} from "lucide-react";
import { useToast } from "../ui/Toast";
import { apiRequest } from "../../lib/api";

// Glass Card Component
const GlassCard = ({ children, className = "" }) => (
    <div className={`relative overflow-hidden rounded-2xl border border-white/10 bg-slate-900/40 backdrop-blur-xl shadow-2xl ${className}`}>
        <div className="absolute top-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-white/20 to-transparent opacity-50"></div>
        <div className="p-6 h-full flex flex-col">
            {children}
        </div>
    </div>
);

// Section Header Component
const SectionHeader = ({ icon: Icon, title, description }) => (
    <div className="flex items-start gap-4 mb-6">
        <div className="p-3 rounded-xl bg-blue-500/20 border border-blue-500/30">
            <Icon size={24} className="text-blue-400" />
        </div>
        <div>
            <h2 className="text-xl font-bold text-white tracking-tight">{title}</h2>
            <p className="text-slate-400 text-sm mt-1">{description}</p>
        </div>
    </div>
);

// Progress Bar for Disk Usage
const DiskUsageBar = ({ used, total, label, mountPoint, fsType }) => {
    const percentage = total > 0 ? (used / total) * 100 : 0;

    // Color based on usage
    let barColor = "bg-emerald-500";
    let bgColor = "bg-emerald-500/20";
    if (percentage > 90) {
        barColor = "bg-rose-500";
        bgColor = "bg-rose-500/20";
    } else if (percentage > 70) {
        barColor = "bg-amber-500";
        bgColor = "bg-amber-500/20";
    }

    const formatSize = (bytes) => {
        if (bytes >= 1024 * 1024 * 1024 * 1024) {
            return (bytes / (1024 * 1024 * 1024 * 1024)).toFixed(2) + " TB";
        } else if (bytes >= 1024 * 1024 * 1024) {
            return (bytes / (1024 * 1024 * 1024)).toFixed(2) + " GB";
        } else if (bytes >= 1024 * 1024) {
            return (bytes / (1024 * 1024)).toFixed(2) + " MB";
        }
        return bytes + " B";
    };

    return (
        <div className="space-y-3 p-4 bg-slate-800/30 rounded-xl border border-white/5">
            <div className="flex justify-between items-start">
                <div>
                    <div className="flex items-center gap-2">
                        <span className="text-white font-medium">{mountPoint}</span>
                        <span className="text-xs px-2 py-0.5 rounded-full bg-slate-700 text-slate-300">
                            {fsType}
                        </span>
                    </div>
                </div>
                <span className="text-white font-mono text-sm">
                    {formatSize(used)} / {formatSize(total)}
                </span>
            </div>

            <div className={`h-4 rounded-full ${bgColor} overflow-hidden`}>
                <div
                    className={`h-full ${barColor} transition-all duration-500 rounded-full`}
                    style={{ width: `${Math.min(percentage, 100)}%` }}
                />
            </div>

            <div className="flex justify-between text-xs text-slate-500">
                <span>{label || "Physical Disk"}</span>
                <span>{percentage.toFixed(1)}% belegt</span>
            </div>
        </div>
    );
};

export default function StorageTab() {
    const toast = useToast();
    const [loading, setLoading] = useState(true);
    const [emptyingTrash, setEmptyingTrash] = useState(false);
    const [savingPaths, setSavingPaths] = useState(false);

    // Data States
    const [disks, setDisks] = useState([]);
    const [trashData, setTrashData] = useState({ size: 0, count: 0 });
    const [settings, setSettings] = useState({
        warningThreshold: 80,
        criticalThreshold: 90,
        autoCleanup: false,
        cleanupAgeDays: 30,
        aiIndexPaths: []
    });

    // New Path Input
    const [newPath, setNewPath] = useState("");

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        try {
            const [hwData, settingsData, trashRes] = await Promise.all([
                apiRequest("/api/v1/system/hardware/storage", { method: "GET" }).catch(() => []),
                apiRequest("/api/v1/storage/settings", { method: "GET" }),
                apiRequest("/api/v1/storage/trash", { method: "GET" }).catch(() => ({ items: [] }))
            ]);

            setDisks(hwData || []);
            setSettings({
                ...settingsData,
                aiIndexPaths: settingsData?.aiIndexPaths || []
            });

            const trashItems = trashRes?.items || [];
            setTrashData({
                count: trashItems.length,
                size: trashItems.reduce((sum, item) => sum + (item.size || 0), 0)
            });

        } catch (err) {
            console.error(err);
            toast.error("Speicherdaten konnten nicht geladen werden");
        } finally {
            setLoading(false);
        }
    };

    const handleEmptyTrash = async () => {
        if (!confirm("Papierkorb endgültig leeren? Diese Aktion kann nicht rückgängig gemacht werden.")) {
            return;
        }

        setEmptyingTrash(true);
        try {
            await apiRequest("/api/v1/storage/trash/empty", { method: "POST" });
            toast.success("Papierkorb geleert");
            setTrashData({ size: 0, count: 0 });
        } catch (err) {
            toast.error("Fehler beim Leeren des Papierkorbs");
        } finally {
            setEmptyingTrash(false);
        }
    };

    const handleAddPath = () => {
        if (!newPath.trim()) return;
        if (settings.aiIndexPaths.includes(newPath.trim())) {
            toast.error("Pfad existiert bereits");
            return;
        }

        const updatedPaths = [...settings.aiIndexPaths, newPath.trim()];
        savePaths(updatedPaths);
        setNewPath("");
    };

    const handleRemovePath = (pathToRemove) => {
        const updatedPaths = settings.aiIndexPaths.filter(p => p !== pathToRemove);
        savePaths(updatedPaths);
    };

    const savePaths = async (paths) => {
        setSavingPaths(true);
        try {
            // We only update the AI Index Paths here, keeping other settings as is
            await apiRequest("/api/v1/storage/settings", {
                method: "PUT",
                body: JSON.stringify({
                    ...settings,
                    aiIndexPaths: paths
                })
            });
            setSettings(prev => ({ ...prev, aiIndexPaths: paths }));
            toast.success("Pfade aktualisiert");
        } catch (err) {
            toast.error("Fehler beim Speichern der Pfade");
        } finally {
            setSavingPaths(false);
        }
    };

    const formatSize = (bytes) => {
        if (bytes >= 1024 * 1024 * 1024) {
            return (bytes / (1024 * 1024 * 1024)).toFixed(2) + " GB";
        } else if (bytes >= 1024 * 1024) {
            return (bytes / (1024 * 1024)).toFixed(2) + " MB";
        }
        return bytes + " B";
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <Loader2 size={48} className="text-blue-400 animate-spin" />
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Physical Storage */}
            <GlassCard>
                <SectionHeader
                    icon={HardDrive}
                    title="Physische Speicher"
                    description="Status der verbundenen Laufwerke"
                />

                <div className="space-y-4">
                    {disks.map((disk, idx) => (
                        <DiskUsageBar
                            key={idx}
                            mountPoint={disk.mount_point}
                            fsType={disk.filesystem}
                            used={disk.used}
                            total={disk.total}
                            label={disk.mount_point === '/' ? "System Root" : "Externer Speicher"}
                        />
                    ))}
                    {disks.length === 0 && (
                        <p className="text-slate-500 italic text-center py-4">Keine Laufwerke gefunden</p>
                    )}
                </div>

                <button
                    onClick={loadData}
                    className="flex items-center gap-2 px-4 py-2 mt-4 text-slate-400 hover:text-white transition-colors"
                >
                    <RefreshCw size={16} />
                    <span className="text-sm">Aktualisieren</span>
                </button>
            </GlassCard>

            {/* AI Index Locations */}
            <GlassCard>
                <SectionHeader
                    icon={Database}
                    title="AI Index Locations"
                    description="Verzeichnisse, die automatisch vom AI Agent überwacht werden"
                />

                <div className="space-y-4">
                    <div className="flex gap-2">
                        <input
                            type="text"
                            value={newPath}
                            onChange={(e) => setNewPath(e.target.value)}
                            placeholder="/mnt/data/documents"
                            className="flex-1 px-4 py-2 bg-slate-800/50 border border-white/10 rounded-xl text-white placeholder:text-slate-500 focus:border-blue-500/50 focus:outline-none"
                            onKeyDown={(e) => e.key === 'Enter' && handleAddPath()}
                        />
                        <button
                            onClick={handleAddPath}
                            disabled={savingPaths || !newPath.trim()}
                            className="px-4 py-2 bg-blue-500/20 hover:bg-blue-500/30 text-blue-400 rounded-xl border border-blue-500/30 transition-all flex items-center gap-2 disabled:opacity-50"
                        >
                            {savingPaths ? <Loader2 size={18} className="animate-spin" /> : <Plus size={18} />}
                            Hinzufügen
                        </button>
                    </div>

                    <div className="space-y-2">
                        {settings.aiIndexPaths.map((path) => (
                            <div key={path} className="flex items-center justify-between p-3 bg-slate-800/30 rounded-xl border border-white/5 group">
                                <div className="flex items-center gap-3">
                                    <Folder size={18} className="text-slate-400" />
                                    <span className="text-slate-200 font-mono text-sm">{path}</span>
                                </div>
                                <button
                                    onClick={() => handleRemovePath(path)}
                                    disabled={savingPaths}
                                    className="p-2 text-slate-500 hover:text-rose-400 hover:bg-rose-500/10 rounded-lg transition-all opacity-0 group-hover:opacity-100"
                                >
                                    <X size={16} />
                                </button>
                            </div>
                        ))}
                        {settings.aiIndexPaths.length === 0 && (
                            <p className="text-slate-500 italic text-sm p-2 text-center">Keine Pfade konfiguriert</p>
                        )}
                    </div>
                </div>
            </GlassCard>

            {/* Trash Bin */}
            <GlassCard>
                <SectionHeader
                    icon={Trash2}
                    title="Papierkorb"
                    description="Gelöschte Dateien, die wiederhergestellt werden können"
                />

                <div className="flex items-center justify-between p-4 bg-slate-800/30 rounded-xl">
                    <div>
                        <p className="text-white font-medium">
                            {trashData.count} Elemente
                        </p>
                        <p className="text-slate-400 text-sm">
                            Größe: {formatSize(trashData.size)}
                        </p>
                    </div>

                    <button
                        onClick={handleEmptyTrash}
                        disabled={emptyingTrash || trashData.count === 0}
                        className="flex items-center gap-2 px-4 py-2 bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 rounded-xl border border-rose-500/30 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {emptyingTrash ? (
                            <Loader2 size={18} className="animate-spin" />
                        ) : (
                            <Trash2 size={18} />
                        )}
                        <span className="font-medium">Leeren</span>
                    </button>
                </div>
            </GlassCard>
        </div>
    );
}
