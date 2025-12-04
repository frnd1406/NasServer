import { useState, useEffect } from "react";
import {
    User,
    Mail,
    Lock,
    Palette,
    Monitor,
    Info,
    Server,
    Cpu,
    HardDrive,
    Shield,
    Save,
    Eye,
    EyeOff,
    Check,
    AlertCircle,
    Loader2,
    ChevronRight
} from "lucide-react";
import { useTheme } from "../components/ThemeToggle";
import { useToast } from "../components/Toast";
import { apiRequest } from "../lib/api";
import { getAuth } from "../utils/auth";

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

// Settings Item Row
const SettingsRow = ({ label, description, children }) => (
    <div className="flex items-center justify-between py-4 border-b border-white/5 last:border-0">
        <div className="flex-1">
            <p className="text-white font-medium">{label}</p>
            {description && <p className="text-slate-500 text-sm mt-0.5">{description}</p>}
        </div>
        <div className="ml-4">
            {children}
        </div>
    </div>
);

export default function Settings() {
    const { theme, toggleTheme } = useTheme();
    const toast = useToast();
    const [activeTab, setActiveTab] = useState("profile");

    // Profile state
    const [profile, setProfile] = useState({
        username: "",
        email: ""
    });
    const [profileLoading, setProfileLoading] = useState(false);

    // Password change state
    const [passwords, setPasswords] = useState({
        current: "",
        new: "",
        confirm: ""
    });
    const [showPasswords, setShowPasswords] = useState({
        current: false,
        new: false,
        confirm: false
    });
    const [passwordSaving, setPasswordSaving] = useState(false);

    // System info state
    const [systemInfo, setSystemInfo] = useState(null);
    const [systemLoading, setSystemLoading] = useState(false);

    // Load user profile
    useEffect(() => {
        const auth = getAuth();
        if (auth.accessToken) {
            // For now, use placeholder - would come from API
            setProfile({
                username: "Admin User",
                email: "admin@nas.local"
            });
        }
    }, []);

    // Load system info for About section
    useEffect(() => {
        if (activeTab === "about") {
            loadSystemInfo();
        }
    }, [activeTab]);

    const loadSystemInfo = async () => {
        setSystemLoading(true);
        try {
            const data = await apiRequest("/api/v1/system/health", { method: "GET" });
            setSystemInfo(data);
        } catch (err) {
            toast.error("Systeminfo konnte nicht geladen werden");
        } finally {
            setSystemLoading(false);
        }
    };

    const handlePasswordChange = async () => {
        if (passwords.new !== passwords.confirm) {
            toast.error("Passwörter stimmen nicht überein");
            return;
        }
        if (passwords.new.length < 8) {
            toast.error("Passwort muss mindestens 8 Zeichen haben");
            return;
        }

        setPasswordSaving(true);
        try {
            await apiRequest("/api/v1/user/password", {
                method: "POST",
                body: JSON.stringify({
                    current_password: passwords.current,
                    new_password: passwords.new
                })
            });
            toast.success("Passwort erfolgreich geändert");
            setPasswords({ current: "", new: "", confirm: "" });
        } catch (err) {
            toast.error(err.message || "Fehler beim Ändern des Passworts");
        } finally {
            setPasswordSaving(false);
        }
    };

    const tabs = [
        { id: "profile", label: "Profil", icon: User },
        { id: "appearance", label: "Erscheinung", icon: Palette },
        { id: "about", label: "Über", icon: Info }
    ];

    return (
        <div className="space-y-6">
            {/* Header */}
            <div>
                <h1 className="text-3xl font-bold text-white tracking-tight">Einstellungen</h1>
                <p className="text-slate-400 mt-2 text-sm">Konfiguriere dein NAS.AI System</p>
            </div>

            <div className="flex flex-col lg:flex-row gap-6">
                {/* Sidebar Navigation */}
                <div className="lg:w-64 flex-shrink-0">
                    <GlassCard className="!p-2">
                        <nav className="space-y-1">
                            {tabs.map((tab) => {
                                const Icon = tab.icon;
                                const isActive = activeTab === tab.id;
                                return (
                                    <button
                                        key={tab.id}
                                        onClick={() => setActiveTab(tab.id)}
                                        className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${isActive
                                                ? "bg-blue-500/20 text-blue-400 border border-blue-500/30"
                                                : "text-slate-400 hover:text-white hover:bg-white/5"
                                            }`}
                                    >
                                        <Icon size={20} />
                                        <span className="font-medium">{tab.label}</span>
                                        {isActive && <ChevronRight size={16} className="ml-auto" />}
                                    </button>
                                );
                            })}
                        </nav>
                    </GlassCard>
                </div>

                {/* Main Content */}
                <div className="flex-1">
                    {/* Profile Section */}
                    {activeTab === "profile" && (
                        <div className="space-y-6">
                            <GlassCard>
                                <SectionHeader
                                    icon={User}
                                    title="Profil Informationen"
                                    description="Deine persönlichen Daten und Kontoeinstellungen"
                                />

                                <div className="space-y-4">
                                    <div>
                                        <label className="block text-sm font-medium text-slate-300 mb-2">Benutzername</label>
                                        <div className="relative">
                                            <User size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
                                            <input
                                                type="text"
                                                value={profile.username}
                                                onChange={(e) => setProfile({ ...profile, username: e.target.value })}
                                                className="w-full pl-10 pr-4 py-3 bg-slate-800/50 border border-white/10 rounded-xl text-white focus:border-blue-500/50 focus:outline-none transition-all"
                                                placeholder="Dein Benutzername"
                                            />
                                        </div>
                                    </div>

                                    <div>
                                        <label className="block text-sm font-medium text-slate-300 mb-2">E-Mail</label>
                                        <div className="relative">
                                            <Mail size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
                                            <input
                                                type="email"
                                                value={profile.email}
                                                onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                                                className="w-full pl-10 pr-4 py-3 bg-slate-800/50 border border-white/10 rounded-xl text-white focus:border-blue-500/50 focus:outline-none transition-all"
                                                placeholder="name@example.com"
                                            />
                                        </div>
                                    </div>

                                    <button className="flex items-center gap-2 px-4 py-2.5 bg-blue-500/20 hover:bg-blue-500/30 text-blue-400 rounded-xl border border-blue-500/30 transition-all mt-4">
                                        <Save size={18} />
                                        <span className="font-medium">Änderungen speichern</span>
                                    </button>
                                </div>
                            </GlassCard>

                            {/* Password Change */}
                            <GlassCard>
                                <SectionHeader
                                    icon={Lock}
                                    title="Passwort ändern"
                                    description="Aktualisiere dein Passwort für mehr Sicherheit"
                                />

                                <div className="space-y-4">
                                    {["current", "new", "confirm"].map((field) => (
                                        <div key={field}>
                                            <label className="block text-sm font-medium text-slate-300 mb-2">
                                                {field === "current" ? "Aktuelles Passwort" : field === "new" ? "Neues Passwort" : "Passwort bestätigen"}
                                            </label>
                                            <div className="relative">
                                                <Lock size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
                                                <input
                                                    type={showPasswords[field] ? "text" : "password"}
                                                    value={passwords[field]}
                                                    onChange={(e) => setPasswords({ ...passwords, [field]: e.target.value })}
                                                    className="w-full pl-10 pr-12 py-3 bg-slate-800/50 border border-white/10 rounded-xl text-white focus:border-blue-500/50 focus:outline-none transition-all"
                                                    placeholder="••••••••"
                                                />
                                                <button
                                                    type="button"
                                                    onClick={() => setShowPasswords({ ...showPasswords, [field]: !showPasswords[field] })}
                                                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-white transition-colors"
                                                >
                                                    {showPasswords[field] ? <EyeOff size={18} /> : <Eye size={18} />}
                                                </button>
                                            </div>
                                        </div>
                                    ))}

                                    {passwords.new && passwords.new.length < 8 && (
                                        <div className="flex items-center gap-2 text-amber-400 text-sm">
                                            <AlertCircle size={16} />
                                            <span>Mindestens 8 Zeichen erforderlich</span>
                                        </div>
                                    )}

                                    {passwords.new && passwords.confirm && passwords.new !== passwords.confirm && (
                                        <div className="flex items-center gap-2 text-rose-400 text-sm">
                                            <AlertCircle size={16} />
                                            <span>Passwörter stimmen nicht überein</span>
                                        </div>
                                    )}

                                    <button
                                        onClick={handlePasswordChange}
                                        disabled={passwordSaving || !passwords.current || !passwords.new || passwords.new !== passwords.confirm}
                                        className="flex items-center gap-2 px-4 py-2.5 bg-violet-500/20 hover:bg-violet-500/30 text-violet-400 rounded-xl border border-violet-500/30 transition-all disabled:opacity-50 disabled:cursor-not-allowed mt-4"
                                    >
                                        {passwordSaving ? (
                                            <>
                                                <Loader2 size={18} className="animate-spin" />
                                                <span className="font-medium">Speichern...</span>
                                            </>
                                        ) : (
                                            <>
                                                <Shield size={18} />
                                                <span className="font-medium">Passwort ändern</span>
                                            </>
                                        )}
                                    </button>
                                </div>
                            </GlassCard>
                        </div>
                    )}

                    {/* Appearance Section */}
                    {activeTab === "appearance" && (
                        <GlassCard>
                            <SectionHeader
                                icon={Palette}
                                title="Erscheinungsbild"
                                description="Passe das Aussehen deiner Oberfläche an"
                            />

                            <div className="space-y-2">
                                <SettingsRow
                                    label="Theme"
                                    description={theme === "dark" ? "Dunkles Theme aktiv" : "Helles Theme aktiv"}
                                >
                                    <button
                                        onClick={toggleTheme}
                                        className={`relative w-14 h-8 rounded-full transition-colors ${theme === "dark" ? "bg-blue-500/30" : "bg-slate-700"
                                            }`}
                                    >
                                        <div className={`absolute top-1 w-6 h-6 rounded-full bg-white shadow-lg transition-all ${theme === "dark" ? "left-7" : "left-1"
                                            }`} />
                                    </button>
                                </SettingsRow>

                                <SettingsRow
                                    label="Sprache"
                                    description="Deutsch"
                                >
                                    <select className="px-3 py-2 bg-slate-800 border border-white/10 rounded-lg text-white focus:outline-none focus:border-blue-500/50">
                                        <option value="de">Deutsch</option>
                                        <option value="en">English</option>
                                    </select>
                                </SettingsRow>

                                <SettingsRow
                                    label="Kompakte Ansicht"
                                    description="Weniger Abstände für mehr Inhalt"
                                >
                                    <button className="relative w-14 h-8 rounded-full bg-slate-700 transition-colors">
                                        <div className="absolute top-1 left-1 w-6 h-6 rounded-full bg-white shadow-lg" />
                                    </button>
                                </SettingsRow>
                            </div>
                        </GlassCard>
                    )}

                    {/* About Section */}
                    {activeTab === "about" && (
                        <div className="space-y-6">
                            <GlassCard>
                                <SectionHeader
                                    icon={Info}
                                    title="Über NAS.AI"
                                    description="Systeminformationen und Version"
                                />

                                <div className="space-y-4">
                                    <div className="flex items-center gap-4 p-4 bg-gradient-to-br from-blue-500/10 to-violet-500/10 rounded-xl border border-blue-500/20">
                                        <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-500 to-violet-600 flex items-center justify-center shadow-lg">
                                            <Server size={32} className="text-white" />
                                        </div>
                                        <div>
                                            <h3 className="text-2xl font-bold text-white">NAS.AI</h3>
                                            <p className="text-slate-400">Version 1.0.0</p>
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-6">
                                        <div className="p-4 bg-slate-800/50 rounded-xl border border-white/5">
                                            <div className="flex items-center gap-2 text-slate-400 text-sm mb-1">
                                                <Monitor size={14} />
                                                <span>Frontend</span>
                                            </div>
                                            <p className="text-white font-medium">React + Vite</p>
                                        </div>

                                        <div className="p-4 bg-slate-800/50 rounded-xl border border-white/5">
                                            <div className="flex items-center gap-2 text-slate-400 text-sm mb-1">
                                                <Server size={14} />
                                                <span>Backend</span>
                                            </div>
                                            <p className="text-white font-medium">Go API</p>
                                        </div>
                                    </div>
                                </div>
                            </GlassCard>

                            <GlassCard>
                                <SectionHeader
                                    icon={Cpu}
                                    title="System Status"
                                    description="Live-Informationen aus dem Backend"
                                />

                                {systemLoading ? (
                                    <div className="flex items-center justify-center py-8">
                                        <Loader2 size={32} className="text-blue-400 animate-spin" />
                                    </div>
                                ) : systemInfo ? (
                                    <div className="space-y-3">
                                        <div className="flex items-center justify-between p-3 bg-slate-800/30 rounded-lg">
                                            <span className="text-slate-400">Status</span>
                                            <span className="flex items-center gap-2 text-emerald-400 font-medium">
                                                <div className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
                                                Online
                                            </span>
                                        </div>
                                        <div className="flex items-center justify-between p-3 bg-slate-800/30 rounded-lg">
                                            <span className="text-slate-400">Uptime</span>
                                            <span className="text-white font-medium">
                                                {systemInfo.uptime || "N/A"}
                                            </span>
                                        </div>
                                        <div className="flex items-center justify-between p-3 bg-slate-800/30 rounded-lg">
                                            <span className="text-slate-400">Hostname</span>
                                            <span className="text-white font-medium font-mono">
                                                {systemInfo.hostname || "nas.local"}
                                            </span>
                                        </div>
                                    </div>
                                ) : (
                                    <div className="text-center py-8 text-slate-500">
                                        <p>Keine Systemdaten verfügbar</p>
                                        <button
                                            onClick={loadSystemInfo}
                                            className="mt-2 text-blue-400 hover:underline"
                                        >
                                            Erneut versuchen
                                        </button>
                                    </div>
                                )}
                            </GlassCard>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
