import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Lock, Loader2, Unlock, Shield, KeyRound, AlertCircle, Settings } from "lucide-react";
import { apiRequest } from "../lib/api";

const API_BASE =
    import.meta.env.VITE_API_BASE_URL ||
    window.location.origin;

export default function VaultUnlock() {
    const [masterPassword, setMasterPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [vaultStatus, setVaultStatus] = useState(null);
    const [showSetup, setShowSetup] = useState(false);
    const [confirmPassword, setConfirmPassword] = useState("");
    const navigate = useNavigate();

    // Check vault status on mount
    useEffect(() => {
        checkVaultStatus();
    }, []);

    const checkVaultStatus = async () => {
        try {
            const res = await fetch(`${API_BASE}/api/v1/vault/status`, {
                credentials: "include",
            });
            if (res.ok) {
                const data = await res.json();
                setVaultStatus(data);
                // If already unlocked, redirect to dashboard
                if (!data.locked) {
                    navigate("/dashboard", { replace: true });
                }
                // If not configured, show setup
                if (!data.configured) {
                    setShowSetup(true);
                }
            }
        } catch (err) {
            console.error("Vault status check failed:", err);
        }
    };

    const handleUnlock = async (e) => {
        e.preventDefault();
        if (!masterPassword) return;

        setLoading(true);
        setError("");

        try {
            const { data, error: apiError } = await apiRequest("/api/v1/system/vault/unlock", {
                method: "POST",
                body: { masterPassword },
            });

            if (apiError) {
                throw new Error(apiError);
            }

            // Success - redirect to dashboard
            navigate("/dashboard", { replace: true });
        } catch (err) {
            setError(err.message || "Entsperren fehlgeschlagen");
        } finally {
            setLoading(false);
        }
    };

    const handleSetup = async (e) => {
        e.preventDefault();
        if (!masterPassword || masterPassword !== confirmPassword) {
            setError("Passwörter stimmen nicht überein");
            return;
        }
        if (masterPassword.length < 8) {
            setError("Passwort muss mindestens 8 Zeichen lang sein");
            return;
        }

        setLoading(true);
        setError("");

        try {
            const { data, error: apiError } = await apiRequest("/api/v1/system/vault/setup", {
                method: "POST",
                body: { masterPassword },
            });

            if (apiError) {
                throw new Error(apiError);
            }

            // Success - redirect to dashboard (auto-unlocked after setup)
            navigate("/dashboard", { replace: true });
        } catch (err) {
            setError(err.message || "Setup fehlgeschlagen");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-[#0a0a0c] text-slate-200 font-sans flex items-center justify-center p-4 relative overflow-hidden">

            {/* Animated Background Blobs */}
            <div className="fixed inset-0 z-0 pointer-events-none overflow-hidden">
                <div className="absolute top-[-10%] left-[-10%] w-[500px] h-[500px] bg-amber-600/20 rounded-full blur-[120px] animate-pulse-glow"></div>
                <div className="absolute bottom-[-10%] right-[-5%] w-[600px] h-[600px] bg-orange-600/10 rounded-full blur-[130px]"></div>
                <div className="absolute top-[40%] left-[30%] w-[300px] h-[300px] bg-yellow-500/10 rounded-full blur-[100px] opacity-60"></div>
            </div>

            {/* Unlock Card */}
            <div className="relative z-10 w-full max-w-md">

                {/* Glass Card */}
                <div className="relative overflow-hidden rounded-2xl border border-white/10 bg-slate-900/40 backdrop-blur-xl shadow-2xl">
                    <div className="absolute top-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-amber-500/30 to-transparent opacity-50"></div>

                    <div className="p-8">

                        {/* Header */}
                        <div className="text-center mb-8">
                            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-amber-500 to-orange-600 mb-4 shadow-lg shadow-amber-500/30">
                                {showSetup ? <KeyRound size={32} className="text-white" /> : <Lock size={32} className="text-white" />}
                            </div>
                            <h1 className="text-3xl font-bold text-white tracking-tight mb-2">
                                {showSetup ? "Vault einrichten" : "Vault entsperren"}
                            </h1>
                            <p className="text-slate-400 text-sm">
                                {showSetup
                                    ? "Erstellen Sie ein Master-Passwort für die Verschlüsselung"
                                    : "Geben Sie Ihr Master-Passwort ein, um das System zu entsperren"
                                }
                            </p>
                        </div>

                        {/* Status Badge */}
                        {vaultStatus && (
                            <div className="mb-6 flex items-center justify-center gap-2">
                                <div className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium ${vaultStatus.locked
                                        ? "bg-amber-500/10 text-amber-400 border border-amber-500/20"
                                        : "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
                                    }`}>
                                    <Shield size={14} />
                                    {vaultStatus.locked ? "Gesperrt" : "Entsperrt"}
                                </div>
                                {vaultStatus.vaultPath && (
                                    <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium bg-slate-500/10 text-slate-400 border border-slate-500/20">
                                        <Settings size={12} />
                                        {vaultStatus.vaultPath.split('/').pop()}
                                    </div>
                                )}
                            </div>
                        )}

                        {/* Error Message */}
                        {error && (
                            <div className="mb-6 p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 animate-in fade-in duration-300">
                                <p className="text-rose-400 text-sm font-medium flex items-center gap-2">
                                    <AlertCircle size={16} />
                                    {error}
                                </p>
                            </div>
                        )}

                        {/* Form */}
                        <form onSubmit={showSetup ? handleSetup : handleUnlock} className="space-y-5">

                            {/* Master Password Input */}
                            <div>
                                <label className="block text-sm font-medium text-slate-300 mb-2">
                                    Master-Passwort
                                </label>
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                        <KeyRound size={18} className="text-slate-400" />
                                    </div>
                                    <input
                                        type="password"
                                        value={masterPassword}
                                        onChange={(e) => setMasterPassword(e.target.value)}
                                        required
                                        minLength={showSetup ? 8 : 1}
                                        placeholder="••••••••••••"
                                        className="w-full pl-10 pr-4 py-3 bg-slate-900/50 border border-white/10 rounded-xl text-white placeholder:text-slate-500 focus:border-amber-500/50 focus:ring-2 focus:ring-amber-500/20 focus:outline-none transition-all"
                                        autoFocus
                                    />
                                </div>
                            </div>

                            {/* Confirm Password (Setup only) */}
                            {showSetup && (
                                <div>
                                    <label className="block text-sm font-medium text-slate-300 mb-2">
                                        Passwort bestätigen
                                    </label>
                                    <div className="relative">
                                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                            <Lock size={18} className="text-slate-400" />
                                        </div>
                                        <input
                                            type="password"
                                            value={confirmPassword}
                                            onChange={(e) => setConfirmPassword(e.target.value)}
                                            required
                                            minLength={8}
                                            placeholder="••••••••••••"
                                            className="w-full pl-10 pr-4 py-3 bg-slate-900/50 border border-white/10 rounded-xl text-white placeholder:text-slate-500 focus:border-amber-500/50 focus:ring-2 focus:ring-amber-500/20 focus:outline-none transition-all"
                                        />
                                    </div>
                                </div>
                            )}

                            {/* Submit Button */}
                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-amber-500/20 hover:bg-amber-500/30 text-amber-400 rounded-xl font-medium transition-all shadow-[0_0_20px_rgba(245,158,11,0.3)] hover:shadow-[0_0_30px_rgba(245,158,11,0.5)] disabled:opacity-50 disabled:cursor-not-allowed border border-amber-500/30 mt-6"
                            >
                                {loading ? (
                                    <>
                                        <Loader2 size={20} className="animate-spin" />
                                        <span>{showSetup ? "Einrichten..." : "Entsperren..."}</span>
                                    </>
                                ) : (
                                    <>
                                        <Unlock size={20} />
                                        <span>{showSetup ? "Vault einrichten" : "Entsperren"}</span>
                                    </>
                                )}
                            </button>
                        </form>

                        {/* Security Note */}
                        <div className="mt-6 pt-6 border-t border-white/5">
                            <p className="text-center text-xs text-slate-500">
                                🔒 Verschlüsselung: AES-256-GCM · Key Derivation: Argon2id
                            </p>
                        </div>
                    </div>
                </div>

                {/* Footer Info */}
                <div className="mt-6 text-center">
                    <p className="text-xs text-slate-500">
                        Zero-Knowledge Encryption · Nur Sie kennen das Passwort
                    </p>
                </div>
            </div>
        </div>
    );
}
