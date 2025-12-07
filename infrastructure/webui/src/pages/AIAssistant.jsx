import { useState, useRef, useEffect } from 'react';
import { Search, Loader2, Sparkles, FileText, ArrowRight, Brain, Zap } from 'lucide-react';

/**
 * AIAssistant - Dedicated full-page AI chat interface
 * Features:
 * - Search bar at top
 * - Answer on left with formatted file references
 * - Files panel on right (appears when results arrive)
 */
export default function AIAssistant() {
    const [query, setQuery] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [hasResult, setHasResult] = useState(false);
    const [answer, setAnswer] = useState('');
    const [sources, setSources] = useState([]);
    const [confidence, setConfidence] = useState('');
    const inputRef = useRef(null);

    // Focus input on mount
    useEffect(() => {
        inputRef.current?.focus();
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!query.trim() || isLoading) return;

        setIsLoading(true);
        setHasResult(false);
        setAnswer('');
        setSources([]);

        try {
            const response = await fetch(`/api/v1/ask?q=${encodeURIComponent(query.trim())}`);
            if (!response.ok) throw new Error('API Error');

            const data = await response.json();

            setAnswer(data.answer || 'Keine Antwort verfügbar.');
            setSources(data.cited_sources || []);
            setConfidence(data.confidence || 'UNBEKANNT');
            setHasResult(true);

        } catch (error) {
            console.error('AI Error:', error);
            setAnswer('Entschuldigung, ich konnte keine Verbindung zum AI-Service herstellen.');
            setHasResult(true);
        } finally {
            setIsLoading(false);
        }
    };

    // Format answer text: highlight file references like [filename.txt]
    const formatAnswer = (text) => {
        if (!text) return null;

        // Split by file references pattern [filename]
        const parts = text.split(/(\[[^\]]+\.txt\]|\[[^\]]+\.json\]|\[[^\]]+\.md\])/g);

        return parts.map((part, index) => {
            if (part.match(/^\[.+\.(txt|json|md)\]$/)) {
                const filename = part.slice(1, -1);
                return (
                    <span
                        key={index}
                        className="inline-flex items-center gap-1 px-2 py-0.5 mx-1 bg-blue-500/20 text-blue-400 rounded-md text-sm font-medium border border-blue-500/30 hover:bg-blue-500/30 transition-colors cursor-pointer"
                    >
                        <FileText size={12} />
                        {filename}
                    </span>
                );
            }
            return <span key={index}>{part}</span>;
        });
    };

    return (
        <div className="min-h-[calc(100vh-120px)] flex flex-col">
            {/* Header with Search */}
            <div className="mb-8">
                <div className="flex items-center gap-3 mb-6">
                    <div className="p-3 bg-gradient-to-br from-blue-500 to-violet-600 rounded-xl shadow-lg shadow-blue-500/20">
                        <Brain size={28} className="text-white" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-white">NAS.AI Assistant</h1>
                        <p className="text-sm text-slate-400">Intelligente Dokumentensuche mit RAG</p>
                    </div>
                </div>

                {/* Search Bar */}
                <form onSubmit={handleSubmit} className="relative">
                    <div className="flex items-center gap-3 p-2 bg-slate-800/80 backdrop-blur-sm rounded-2xl border border-white/10 shadow-xl focus-within:border-blue-500/50 focus-within:shadow-blue-500/10 transition-all">
                        <div className="pl-4">
                            <Search size={22} className="text-slate-400" />
                        </div>
                        <input
                            ref={inputRef}
                            type="text"
                            value={query}
                            onChange={(e) => setQuery(e.target.value)}
                            placeholder="Stelle eine Frage über deine Dokumente..."
                            className="flex-1 bg-transparent text-white text-lg py-3 outline-none placeholder:text-slate-500"
                            disabled={isLoading}
                        />
                        <button
                            type="submit"
                            disabled={isLoading || !query.trim()}
                            className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-violet-600 text-white font-semibold rounded-xl hover:from-blue-500 hover:to-violet-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-blue-500/25"
                        >
                            {isLoading ? (
                                <>
                                    <Loader2 size={18} className="animate-spin" />
                                    <span>Suche...</span>
                                </>
                            ) : (
                                <>
                                    <Sparkles size={18} />
                                    <span>Fragen</span>
                                </>
                            )}
                        </button>
                    </div>
                </form>
            </div>

            {/* Results Area */}
            {(isLoading || hasResult) && (
                <div className={`flex-1 grid gap-6 animate-fade-in-scale ${sources.length > 0 ? 'grid-cols-1 lg:grid-cols-3' : 'grid-cols-1'}`}>

                    {/* Left: Answer Panel */}
                    <div className={`${sources.length > 0 ? 'lg:col-span-2' : ''} bg-slate-800/60 backdrop-blur-sm rounded-2xl border border-white/10 p-6 shadow-xl`}>
                        <div className="flex items-center gap-2 mb-4 pb-4 border-b border-white/10">
                            <div className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
                            <span className="font-semibold text-white">NAS.AI Antwort</span>
                            {confidence && hasResult && (
                                <span className={`ml-auto text-xs px-2 py-1 rounded-full ${confidence === 'HOCH' ? 'bg-emerald-500/20 text-emerald-400' :
                                        confidence === 'MITTEL' ? 'bg-amber-500/20 text-amber-400' :
                                            'bg-slate-500/20 text-slate-400'
                                    }`}>
                                    {confidence}
                                </span>
                            )}
                        </div>

                        {isLoading ? (
                            <div className="flex items-center gap-3 text-slate-400">
                                <Loader2 size={20} className="animate-spin text-blue-400" />
                                <span>Analysiere Dokumente...</span>
                            </div>
                        ) : (
                            <div className="text-slate-200 leading-relaxed text-lg whitespace-pre-wrap">
                                {formatAnswer(answer)}
                            </div>
                        )}
                    </div>

                    {/* Right: Files Panel */}
                    {sources.length > 0 && (
                        <div className="bg-slate-800/60 backdrop-blur-sm rounded-2xl border border-white/10 p-5 shadow-xl animate-slide-in-right">
                            <div className="flex items-center gap-2 mb-4 pb-4 border-b border-white/10">
                                <Zap size={18} className="text-cyan-400" />
                                <span className="font-semibold text-white">Quellen</span>
                                <span className="ml-auto text-xs text-slate-400">{sources.length} Dokumente</span>
                            </div>

                            <div className="flex flex-col gap-3">
                                {sources.map((source, index) => {
                                    const fileName = source.file_path?.split('/').pop() || source.file_id;
                                    const similarity = source.similarity || 0;
                                    const percent = Math.round(similarity * 100);

                                    return (
                                        <button
                                            key={index}
                                            className="group flex items-center gap-3 p-4 bg-slate-700/50 hover:bg-blue-600/30 rounded-xl border border-white/5 hover:border-blue-500/40 transition-all duration-200 text-left animate-slide-in-right"
                                            style={{ animationDelay: `${index * 100}ms` }}
                                        >
                                            <div className="p-2 bg-blue-500/20 rounded-lg group-hover:bg-blue-500/30 transition-colors">
                                                <FileText size={18} className="text-blue-400" />
                                            </div>
                                            <div className="flex-1 min-w-0">
                                                <p className="text-white font-medium truncate">{fileName}</p>
                                                <p className="text-xs text-slate-400 truncate">{source.file_path}</p>
                                            </div>
                                            <div className="flex flex-col items-end gap-1">
                                                <span className={`text-xs font-semibold px-2 py-0.5 rounded-full ${percent >= 70 ? 'bg-emerald-500/20 text-emerald-400' :
                                                        percent >= 50 ? 'bg-blue-500/20 text-blue-400' :
                                                            'bg-amber-500/20 text-amber-400'
                                                    }`}>
                                                    {percent}%
                                                </span>
                                                <ArrowRight size={14} className="text-slate-500 group-hover:text-white transition-colors" />
                                            </div>
                                        </button>
                                    );
                                })}
                            </div>
                        </div>
                    )}
                </div>
            )}

            {/* Empty State */}
            {!isLoading && !hasResult && (
                <div className="flex-1 flex flex-col items-center justify-center text-center py-20">
                    <div className="p-6 bg-gradient-to-br from-blue-500/10 to-violet-500/10 rounded-full mb-6 border border-white/5">
                        <Brain size={48} className="text-blue-400" />
                    </div>
                    <h2 className="text-xl font-semibold text-white mb-2">Frag mich etwas!</h2>
                    <p className="text-slate-400 max-w-md mb-8">
                        Ich durchsuche deine Dokumente mit KI und finde relevante Informationen.
                    </p>
                    <div className="flex flex-wrap gap-2 justify-center">
                        {[
                            'Was sind meine Krypto-Gewinne?',
                            'Welche Noten hat Finn?',
                            'Zeige mir Rechnungen',
                            'Wann ist meine nächste Reise?'
                        ].map((suggestion, i) => (
                            <button
                                key={i}
                                onClick={() => setQuery(suggestion)}
                                className="px-4 py-2 bg-slate-800/60 hover:bg-slate-700/60 text-slate-300 hover:text-white rounded-full text-sm border border-white/10 hover:border-white/20 transition-all"
                            >
                                {suggestion}
                            </button>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
