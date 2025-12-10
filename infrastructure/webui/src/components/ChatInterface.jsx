import React, { useState, useRef, useEffect } from 'react';
import { Send, Bot, User, Loader2, FileText, AlertCircle, CheckCircle2 } from 'lucide-react';
import { apiRequest } from '../lib/api';

const ChatInterface = () => {
    const [messages, setMessages] = useState([
        {
            role: 'assistant',
            content: 'Hallo! Ich bin dein AI-Assistent. Frag mich etwas über deine Dokumente, z.B. "Wie hoch waren die Serverkosten?"',
            timestamp: new Date()
        }
    ]);
    const [input, setInput] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const messagesEndRef = useRef(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!input.trim() || isLoading) return;

        const userMessage = {
            role: 'user',
            content: input,
            timestamp: new Date()
        };

        setMessages(prev => [...prev, userMessage]);
        setInput('');
        setIsLoading(true);

        // Helper function with retry logic for 502/504 errors
        const fetchWithRetry = async (url, options, retryCount = 0) => {
            try {
                return await apiRequest(url, options);
            } catch (error) {
                // Auto-retry once on gateway errors (502/504)
                const status = error.status || error.response?.status;
                const isGatewayError = status === 502 || status === 504;

                if (isGatewayError && retryCount < 1) {
                    console.warn(`Gateway error ${status}, retrying in 1s...`);
                    await new Promise(resolve => setTimeout(resolve, 1000));
                    return fetchWithRetry(url, options, retryCount + 1);
                }
                throw error; // Re-throw for non-retryable errors
            }
        };

        try {
            const data = await fetchWithRetry(
                `/api/v1/ask?q=${encodeURIComponent(userMessage.content)}`,
                { method: 'GET' }
            );

            const aiMessage = {
                role: 'assistant',
                content: data.answer,
                sources: data.cited_sources || [],
                confidence: data.confidence,
                timestamp: new Date()
            };

            setMessages(prev => [...prev, aiMessage]);
        } catch (error) {
            console.error('Chat error:', error);
            setMessages(prev => [...prev, {
                role: 'assistant',
                content: 'Entschuldigung, es gab einen Fehler bei der Verarbeitung deiner Anfrage. Bitte versuche es später erneut.',
                isError: true,
                timestamp: new Date()
            }]);
        } finally {
            setIsLoading(false);
        }
    };

    const ConfidenceBadge = ({ level }) => {
        if (!level) return null;

        const colors = {
            'HOCH': 'bg-green-500/20 text-green-400 border-green-500/30',
            'MITTEL': 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30',
            'NIEDRIG': 'bg-red-500/20 text-red-400 border-red-500/30'
        };

        return (
            <span className={`text-xs px-2 py-0.5 rounded border ${colors[level] || colors['MITTEL']} ml-2`}>
                {level}
            </span>
        );
    };

    return (
        <div className="flex flex-col h-[600px] bg-[#0f172a] rounded-xl border border-white/10 overflow-hidden shadow-2xl">
            {/* Header */}
            <div className="p-4 border-b border-white/10 bg-[#1e293b]/50 backdrop-blur flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-indigo-500/20 flex items-center justify-center text-indigo-400">
                    <Bot size={20} />
                </div>
                <div>
                    <h3 className="font-medium text-white">AI Knowledge Assistant</h3>
                    <p className="text-xs text-slate-400">Powered by Llama 3.2 & RAG</p>
                </div>
            </div>

            {/* Messages Area */}
            <div className="flex-1 overflow-y-auto p-4 space-y-6 scrollbar-thin scrollbar-thumb-white/10 scrollbar-track-transparent">
                {messages.map((msg, idx) => {
                    // Ghost Message Filter: Don't render empty bubbles
                    if (!msg.content && !msg.isLoading && !msg.isError) return null;

                    return (
                        <div key={idx} className={`flex gap-4 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}>
                            {/* Avatar */}
                            <div className={`w-8 h-8 rounded-lg flex-shrink-0 flex items-center justify-center ${msg.role === 'user'
                                ? 'bg-indigo-500 text-white'
                                : 'bg-slate-700 text-indigo-400'
                                }`}>
                                {msg.role === 'user' ? <User size={18} /> : <Bot size={18} />}
                            </div>

                            {/* Message Bubble */}
                            <div className={`flex flex-col max-w-[80%] ${msg.role === 'user' ? 'items-end' : 'items-start'}`}>
                                <div className={`rounded-2xl px-4 py-3 ${msg.role === 'user'
                                    ? 'bg-indigo-600 text-white rounded-tr-none'
                                    : 'bg-slate-800/80 text-slate-200 rounded-tl-none border border-white/5'
                                    }`}>
                                    <p className="whitespace-pre-wrap leading-relaxed">{msg.content}</p>
                                </div>

                                {/* Metadata (Sources & Confidence) for AI messages */}
                                {msg.role === 'assistant' && !msg.isError && (
                                    <div className="mt-2 space-y-2 w-full">
                                        {/* Confidence Badge */}
                                        {msg.confidence && (
                                            <div className="flex items-center gap-2 text-xs text-slate-400">
                                                <span>Konfidenz:</span>
                                                <ConfidenceBadge level={msg.confidence} />
                                            </div>
                                        )}

                                        {/* Sources List */}
                                        {msg.sources && msg.sources.length > 0 && (
                                            <div className="bg-slate-900/50 rounded-lg p-3 border border-white/5 text-sm">
                                                <div className="flex items-center gap-2 text-slate-400 mb-2 text-xs uppercase tracking-wider font-medium">
                                                    <FileText size={12} />
                                                    Quellen
                                                </div>
                                                <div className="space-y-2">
                                                    {msg.sources.map((source, sIdx) => (
                                                        <div key={sIdx} className="flex items-start gap-2 text-slate-300 bg-white/5 p-2 rounded hover:bg-white/10 transition-colors cursor-pointer group">
                                                            <div className="mt-0.5 text-indigo-400 group-hover:text-indigo-300">
                                                                <CheckCircle2 size={14} />
                                                            </div>
                                                            <div className="flex-1 min-w-0">
                                                                <p className="truncate font-medium text-xs text-indigo-300">{source.file_id}</p>
                                                                <p className="text-xs text-slate-500 mt-0.5">Match: {Math.round(source.similarity * 100)}%</p>
                                                            </div>
                                                        </div>
                                                    ))}
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                )}

                                <span className="text-[10px] text-slate-500 mt-1 px-1">
                                    {msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                </span>
                            </div>
                        </div>
                    );
                })}

                {isLoading && (
                    <div className="flex gap-4">
                        <div className="w-8 h-8 rounded-lg bg-slate-700 text-indigo-400 flex-shrink-0 flex items-center justify-center">
                            <Bot size={18} />
                        </div>
                        <div className="bg-slate-800/80 rounded-2xl rounded-tl-none px-4 py-3 border border-white/5 flex items-center gap-2">
                            <Loader2 size={16} className="animate-spin text-indigo-400" />
                            <span className="text-slate-400 text-sm">Analysiere Dokumente...</span>
                        </div>
                    </div>
                )}
                <div ref={messagesEndRef} />
            </div>

            {/* Input Area */}
            <form onSubmit={handleSubmit} className="p-4 bg-[#1e293b]/50 border-t border-white/10">
                <div className="relative flex gap-2">
                    <input
                        type="text"
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        placeholder="Stelle eine Frage an deine Dokumente..."
                        className="flex-1 bg-slate-900/50 border border-white/10 rounded-xl px-4 py-3 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500/50 transition-all"
                        disabled={isLoading}
                    />
                    <button
                        type="submit"
                        disabled={!input.trim() || isLoading}
                        className="bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed text-white p-3 rounded-xl transition-all shadow-lg shadow-indigo-500/20 flex items-center justify-center min-w-[50px]"
                    >
                        {isLoading ? <Loader2 size={20} className="animate-spin" /> : <Send size={20} />}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default ChatInterface;
