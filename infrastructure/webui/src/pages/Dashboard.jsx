import { useState } from "react";
import { useSystemHealth } from "../hooks/useSystemHealth";
import { useVault } from "../context/VaultContext";
import {
  ShieldCheck,
  ShieldAlert,
  Loader2,
  Clock,
  Archive,
  Server,
  Cpu,
  MemoryStick,
  HardDrive,
  Lock,
  AlertTriangle,
} from "lucide-react";
import { DashboardSkeleton } from "../components/Skeleton";
import SystemHealthCard from "../components/SystemHealthCard";

// Glass Card Component
const GlassCard = ({ children, className = "" }) => (
  <div className={`relative overflow-hidden rounded-2xl border border-white/10 bg-slate-900/40 backdrop-blur-xl shadow-2xl ${className}`}>
    <div className="absolute top-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-white/20 to-transparent opacity-50"></div>
    <div className="p-6 h-full flex flex-col">
      {children}
    </div>
  </div>
);

export default function Dashboard() {
  const { data, isLoading } = useSystemHealth();
  const { vaultConfig, isLoading: vaultLoading } = useVault();
  const settings = data?.settings;
  const lastBackup = data?.lastBackup;
  const snapshotCount = data?.snapshotCount || 0;
  const latestMetric = data?.latestMetric;
  const metricsError = data?.metricsError;

  // Calculate next backup time
  const getNextBackupTime = () => {
    if (!settings?.backup_schedule) return "Nicht geplant";

    const now = new Date();
    const [hours, minutes] = settings.backup_schedule.split(':');
    const nextRun = new Date();
    nextRun.setHours(parseInt(hours), parseInt(minutes), 0, 0);

    // If the time has passed today, set it to tomorrow
    if (nextRun < now) {
      nextRun.setDate(nextRun.getDate() + 1);
    }

    // Check if it's today or tomorrow
    const isToday = nextRun.toDateString() === now.toDateString();
    const timeStr = settings.backup_schedule;

    return isToday ? `Heute, ${timeStr} Uhr` : `Morgen, ${timeStr} Uhr`;
  };

  const formatLastBackupTime = () => {
    if (!lastBackup) return "Kein Backup vorhanden";

    const date = new Date(lastBackup.modTime || lastBackup.created_at);
    const now = new Date();
    const diffMs = now - date;
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffHours / 24);

    if (diffDays > 0) {
      return `vor ${diffDays} Tag${diffDays > 1 ? 'en' : ''}`;
    } else if (diffHours > 0) {
      return `vor ${diffHours} Stunde${diffHours > 1 ? 'n' : ''}`;
    } else {
      return "Kürzlich";
    }
  };

  const isBackupActive = settings?.auto_backup_enabled ?? true;
  const formatPct = (v) => (typeof v === "number" ? `${Math.round(v)}%` : "—");

  // Show skeleton during initial load
  if (isLoading && !data) {
    return <DashboardSkeleton />;
  }

  return (
    <div className="space-y-6">

      {/* Welcome Header */}
      <div>
        <h1 className="text-3xl font-bold text-white tracking-tight">
          Dashboard
        </h1>
        <p className="text-slate-400 mt-2 text-sm">
          Systemübersicht und Status
        </p>
      </div>

      {/* Vault Not Configured Warning */}
      {!vaultLoading && !vaultConfig && (
        <div className="relative overflow-hidden rounded-2xl border border-amber-500/30 bg-amber-500/10 backdrop-blur-xl shadow-lg">
          <div className="absolute top-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-amber-400/40 to-transparent"></div>
          <div className="p-4 flex items-center gap-4">
            <div className="p-3 rounded-xl bg-amber-500/20 border border-amber-500/30 shrink-0">
              <Lock size={24} className="text-amber-400" />
            </div>
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <AlertTriangle size={16} className="text-amber-400" />
                <p className="text-amber-300 font-semibold text-sm uppercase tracking-wider">Vault nicht eingerichtet</p>
              </div>
              <p className="text-amber-200/80 text-sm">
                Für verschlüsselte Dateispeicherung gehe zu <a href="/files" className="text-amber-400 hover:text-amber-300 underline font-medium">Dateien → vault</a> und richte die Ende-zu-Ende-Verschlüsselung ein.
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Main Grid - 3 Cards */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

        {/* 1. System Ressourcen Card (Größer - Spans 2 columns on large screens) */}
        <GlassCard className="lg:col-span-2 hover:bg-blue-500/5 transition-colors">
          <div className="flex items-start justify-between mb-6">
            <div className="flex items-center gap-3">
              <div className="p-3 rounded-xl bg-blue-500/20 border border-blue-500/30">
                <Server size={24} className="text-blue-400" />
              </div>
              <div>
                <p className="text-slate-400 text-xs uppercase tracking-wider">System Ressourcen</p>
                <p className="text-lg font-semibold text-white mt-0.5">Hardware Status</p>
              </div>
            </div>

            {/* Live Indicator */}
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20">
              <div className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></div>
              <span className="text-emerald-400 text-xs font-medium uppercase tracking-wider">Live</span>
            </div>
          </div>

          {/* CPU, RAM, Storage - Clean Bars without timestamps */}
          <div className="space-y-4 flex-1">

            {metricsError && (
              <div className="p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-center">
                <ShieldAlert size={20} className="text-red-400 mx-auto mb-2" />
                <p className="text-red-400 font-semibold">Metriken konnten nicht geladen werden</p>
                <p className="text-slate-400 text-xs mt-1">{metricsError}</p>
              </div>
            )}

            {/* CPU */}
            <div className="p-4 rounded-xl bg-white/5 border border-white/5 hover:bg-white/10 transition-colors">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-cyan-500/20">
                    <Cpu size={20} className="text-cyan-400" />
                  </div>
                  <div>
                    <p className="text-slate-300 text-sm font-medium">CPU Auslastung</p>
                    <p className="text-xs text-slate-500">Prozessor</p>
                  </div>
                </div>
                <p className="text-white font-mono text-2xl font-bold">{formatPct(latestMetric?.cpu_usage)}</p>
              </div>
              {/* Progress Bar */}
              <div className="w-full bg-slate-800/50 rounded-full h-2 overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-cyan-500 to-cyan-400 transition-all duration-500 rounded-full"
                  style={{ width: `${latestMetric?.cpu_usage || 0}%` }}
                ></div>
              </div>
            </div>

            {/* RAM */}
            <div className="p-4 rounded-xl bg-white/5 border border-white/5 hover:bg-white/10 transition-colors">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-violet-500/20">
                    <MemoryStick size={20} className="text-violet-400" />
                  </div>
                  <div>
                    <p className="text-slate-300 text-sm font-medium">Arbeitsspeicher</p>
                    <p className="text-xs text-slate-500">RAM</p>
                  </div>
                </div>
                <p className="text-white font-mono text-2xl font-bold">{formatPct(latestMetric?.ram_usage)}</p>
              </div>
              {/* Progress Bar */}
              <div className="w-full bg-slate-800/50 rounded-full h-2 overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-violet-500 to-violet-400 transition-all duration-500 rounded-full"
                  style={{ width: `${latestMetric?.ram_usage || 0}%` }}
                ></div>
              </div>
            </div>

            {/* Storage */}
            <div className="p-4 rounded-xl bg-white/5 border border-white/5 hover:bg-white/10 transition-colors">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-amber-500/20">
                    <HardDrive size={20} className="text-amber-400" />
                  </div>
                  <div>
                    <p className="text-slate-300 text-sm font-medium">Speicher</p>
                    <p className="text-xs text-slate-500">Disk Space</p>
                  </div>
                </div>
                <p className="text-white font-mono text-2xl font-bold">{formatPct(latestMetric?.disk_usage)}</p>
              </div>
              {/* Progress Bar */}
              <div className="w-full bg-slate-800/50 rounded-full h-2 overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-amber-500 to-amber-400 transition-all duration-500 rounded-full"
                  style={{ width: `${latestMetric?.disk_usage || 0}%` }}
                ></div>
              </div>
            </div>

          </div>
        </GlassCard>

        {/* 2. Security & Backup Card (Combined with Snapshots) */}
        <GlassCard className={`${isBackupActive ? 'hover:bg-emerald-500/5' : 'hover:bg-slate-500/5'} transition-colors`}>
          <div className="flex items-start justify-between mb-4">
            <div className={`p-3 rounded-xl ${isBackupActive ? 'bg-emerald-500/20 border-emerald-500/30' : 'bg-slate-700/50 border-slate-500/30'} border`}>
              {isBackupActive ? (
                <ShieldCheck size={24} className="text-emerald-400" />
              ) : (
                <ShieldAlert size={24} className="text-slate-400" />
              )}
            </div>
            {isLoading && (
              <Loader2 size={16} className="animate-spin text-slate-400" />
            )}
          </div>

          <div className="flex-1 space-y-4">
            <div>
              <p className="text-slate-400 text-xs uppercase tracking-wider">Datensicherheit</p>
              <p className={`text-2xl font-bold mt-2 ${isBackupActive ? 'text-emerald-400' : 'text-slate-400'}`}>
                {isBackupActive ? 'Auto-Backup' : 'Manuell'}
              </p>
            </div>

            {/* Backup Schedule */}
            <div className="p-3 rounded-lg bg-white/5 border border-white/10">
              <div className="flex items-center gap-2 mb-2">
                <Clock size={14} className={isBackupActive ? 'text-emerald-400' : 'text-slate-500'} />
                <span className="text-slate-400 text-xs font-medium">Nächster Lauf</span>
              </div>
              <p className={`text-sm font-semibold ${isBackupActive ? 'text-emerald-400' : 'text-slate-400'}`}>
                {getNextBackupTime()}
              </p>
            </div>

            {/* Last Snapshot */}
            <div className="p-3 rounded-lg bg-white/5 border border-white/10">
              <div className="flex items-center gap-2 mb-2">
                <Archive size={14} className="text-blue-400" />
                <span className="text-slate-400 text-xs font-medium">Letzter Snapshot</span>
              </div>
              <p className="text-sm font-semibold text-blue-400">
                {formatLastBackupTime()}
              </p>
            </div>

            {/* Snapshot Count */}
            <div className="p-3 rounded-lg bg-blue-500/10 border border-blue-500/20">
              <div className="flex items-center justify-between">
                <span className="text-slate-400 text-xs font-medium">Gespeicherte Snapshots</span>
                <span className="text-blue-400 font-bold text-lg">{snapshotCount}</span>
              </div>
              {settings?.backup_retention && (
                <p className="text-slate-500 text-xs mt-1">
                  {settings.backup_retention} Tage Aufbewahrung
                </p>
              )}
            </div>
          </div>
        </GlassCard>

      </div>

      {/* 3. Live System Monitor (Phase 8) */}
      <div className="mt-8">
        <SystemHealthCard />
      </div>
    </div>
  );
}
