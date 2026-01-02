import { useState, useEffect } from "react";
import {
    HardDrive,
    Trash2,
    Folder,
    File,
    AlertTriangle,
    Loader2,
    RefreshCw
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
const DiskUsageBar = ({ used, total, label }) => {
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
        <div className="space-y-2">
            <div className="flex justify-between text-sm">
                <span className="text-slate-400">{label}</span>
                <span className="text-white font-mono">
                    {formatSize(used)} / {formatSize(total)}
                </span>
            </div>
            <div className={`h-4 rounded-full ${bgColor} overflow-hidden`}>
                <div
                    className={`h-full ${barColor} transition-all duration-500 rounded-full`}
                    style={{ width: `${Math.min(percentage, 100)}%` }}
                />
            </div>
            <div className="text-right text-xs text-slate-500">
                {percentage.toFixed(1)}% verwendet
            </div>
        </div>
    );
};

export default function StorageTab() {
    const toast = useToast();
    const [loading, setLoading] = useState(true);
    const [emptyingTrash, setEmptyingTrash] = useState(false);
    const [storageData, setStorageData] = useState({
        disk_total: 0,
        disk_used: 0,
        disk_free: 0,
        file_count: 0,
        folder_count: 0,
        trash_size: 0,
        trash_count: 0
    });

    useEffect(() => {
        loadStorageData();
    }, []);

    const loadStorageData = async () => {
        setLoading(true);
        try {
            // Get health data which includes disk stats
            const healthData = await apiRequest("/health", { method: "GET" });

            // Get trash info
            let trashData = { size: 0, count: 0 };
            try {
                const trashResponse = await apiRequest("/api/v1/storage/trash", { method: "GET" });
                if (trashResponse?.items) {
                    trashData.count = trashResponse.items.length;
                    trashData.size = trashResponse.items.reduce((sum, item) => sum + (item.size || 0), 0);
                }
            } catch (e) {
                console.log("Trash info not available");
            }

            setStorageData({
                disk_total: healthData?.disk_total || 500 * 1024 * 1024 * 1024,  // Fallback 500GB
                disk_used: healthData?.disk_used || 0,
                disk_free: healthData?.disk_free || 0,
                file_count: healthData?.file_count || 0,
                folder_count: healthData?.folder_count || 0,
                trash_size: trashData.size,
                trash_count: trashData.count
            });
        } catch (err) {
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
            setStorageData(prev => ({ ...prev, trash_size: 0, trash_count: 0 }));
        } catch (err) {
            toast.error("Fehler beim Leeren des Papierkorbs");
        } finally {
            setEmptyingTrash(false);
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

    const usagePercent = storageData.disk_total > 0
        ? (storageData.disk_used / storageData.disk_total) * 100
        : 0;

    return (
        <div className="space-y-6">
            {/* Disk Usage */}
            <GlassCard>
                <SectionHeader
                    icon={HardDrive}
                    title="Speichernutzung"
                    description="Übersicht über den verfügbaren Speicherplatz"
                />

                <div className="space-y-6">
                    {/* Main Usage Bar */}
                    <DiskUsageBar
                        used={storageData.disk_used}
                        total={storageData.disk_total}
                        label="Hauptspeicher"
                    />

                    {/* Warning if > 90% */}
                    {usagePercent > 90 && (
                        <div className="flex items-center gap-3 p-4 bg-rose-500/10 border border-rose-500/30 rounded-xl">
                            <AlertTriangle size={24} className="text-rose-400" />
                            <div>
                                <p className="text-rose-400 font-medium">Speicher kritisch!</p>
                                <p className="text-slate-400 text-sm">Weniger als 10% freier Speicher verfügbar.</p>
                            </div>
                        </div>
                    )}

                    {/* Stats Grid */}
                    <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-4">
                        <div className="p-4 bg-slate-800/30 rounded-xl border border-white/5">
                            <div className="flex items-center gap-2 text-slate-400 text-sm mb-1">
                                <HardDrive size={14} />
                                <span>Belegt</span>
                            </div>
                            <p className="text-white font-mono text-lg">{formatSize(storageData.disk_used)}</p>
                        </div>

                        <div className="p-4 bg-slate-800/30 rounded-xl border border-white/5">
                            <div className="flex items-center gap-2 text-slate-400 text-sm mb-1">
                                <HardDrive size={14} />
                                <span>Frei</span>
                            </div>
                            <p className="text-white font-mono text-lg">{formatSize(storageData.disk_free)}</p>
                        </div>

                        <div className="p-4 bg-slate-800/30 rounded-xl border border-white/5">
                            <div className="flex items-center gap-2 text-slate-400 text-sm mb-1">
                                <File size={14} />
                                <span>Dateien</span>
                            </div>
                            <p className="text-white font-mono text-lg">{storageData.file_count.toLocaleString()}</p>
                        </div>

                        <div className="p-4 bg-slate-800/30 rounded-xl border border-white/5">
                            <div className="flex items-center gap-2 text-slate-400 text-sm mb-1">
                                <Folder size={14} />
                                <span>Ordner</span>
                            </div>
                            <p className="text-white font-mono text-lg">{storageData.folder_count.toLocaleString()}</p>
                        </div>
                    </div>
                </div>

                <button
                    onClick={loadStorageData}
                    className="flex items-center gap-2 px-4 py-2 mt-4 text-slate-400 hover:text-white transition-colors"
                >
                    <RefreshCw size={16} />
                    <span className="text-sm">Aktualisieren</span>
                </button>
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
                            {storageData.trash_count} Elemente
                        </p>
                        <p className="text-slate-400 text-sm">
                            Größe: {formatSize(storageData.trash_size)}
                        </p>
                    </div>

                    <button
                        onClick={handleEmptyTrash}
                        disabled={emptyingTrash || storageData.trash_count === 0}
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
