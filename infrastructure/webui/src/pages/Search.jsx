import { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { Search as SearchIcon, FileText, Loader2, Sparkles, RefreshCw } from "lucide-react";
import { searchFiles } from "../lib/api";

export default function Search() {
  const [searchParams] = useSearchParams();
  const [query, setQuery] = useState(searchParams.get('q') || "");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [isFocused, setIsFocused] = useState(false);

  // Auto-search wenn Query aus URL kommt
  useEffect(() => {
    const urlQuery = searchParams.get('q');
    if (urlQuery && urlQuery.trim()) {
      setQuery(urlQuery);
      performSearch(urlQuery);
    }
  }, [searchParams]);

  const performSearch = async (searchQuery) => {
    if (!searchQuery.trim()) {
      setError("Bitte geben Sie einen Suchbegriff ein");
      return;
    }

    setLoading(true);
    setError("");
    setResults(null);

    try {
      const data = await searchFiles(searchQuery);
      setResults(data);
    } catch (err) {
      setError(err.message || "Fehler bei der Suche");
      console.error("Search error:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    performSearch(query);
  };

  const highlightText = (text, query) => {
    if (!query || !text) return text;

    const words = query.toLowerCase().split(/\s+/);
    let highlighted = text;

    words.forEach(word => {
      const regex = new RegExp(`(${word})`, 'gi');
      highlighted = highlighted.replace(regex, '<mark class="bg-yellow-200 dark:bg-yellow-900">$1</mark>');
    });

    return highlighted;
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="text-center mb-8">
        <h1 className="text-4xl font-bold text-white mb-3 tracking-tight">
          Semantische Suche
        </h1>
        <p className="text-slate-400">
          KI-gestützte Vektoranalyse durchsucht deine Dokumente
        </p>
      </div>

      {/* AI Search Box with Glow Effect */}
      <div className="relative w-full max-w-2xl mx-auto">
        <div className={`relative group transition-all duration-300 ${isFocused ? 'scale-[1.02]' : ''}`}>
          {/* Glow Background */}
          <div className={`absolute -inset-0.5 bg-gradient-to-r from-blue-500 to-violet-600 rounded-2xl blur opacity-30 group-hover:opacity-60 transition duration-500 ${loading ? 'animate-pulse' : ''}`} />

          <form onSubmit={handleSearch} className="relative flex items-center bg-slate-900 border border-white/10 rounded-xl shadow-2xl">
            <div className="pl-5 text-slate-400">
              {loading ? (
                <RefreshCw className="animate-spin text-blue-400" size={22} />
              ) : (
                <Sparkles className={`transition-colors ${isFocused ? 'text-violet-400' : 'text-slate-500'}`} size={22} />
              )}
            </div>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onFocus={() => setIsFocused(true)}
              onBlur={() => setIsFocused(false)}
              placeholder="Frag deine Daten... (z.B. 'Zeige mir Rechnungen über 50€')"
              className="w-full bg-transparent text-white px-4 py-5 focus:outline-none placeholder-slate-500 text-lg"
              disabled={loading}
            />
            <div className="pr-4 flex items-center gap-2">
              <kbd className="hidden sm:inline-block px-2 py-1 text-xs font-semibold text-slate-500 bg-slate-800 border border-slate-700 rounded-lg">Enter</kbd>
              <button
                type="submit"
                disabled={loading}
                className="p-2 rounded-lg bg-blue-600 hover:bg-blue-500 text-white disabled:bg-slate-700 disabled:cursor-not-allowed transition-colors"
              >
                <SearchIcon size={20} />
              </button>
            </div>
          </form>
        </div>

        {error && (
          <div className="mt-4 p-4 bg-red-500/10 border border-red-500/20 rounded-xl">
            <p className="text-red-400">{error}</p>
          </div>
        )}
      </div>

      {/* Results */}
      {results && (
        <div className="space-y-4 max-w-4xl mx-auto">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-semibold text-white">
              Ergebnisse für "{results.query}"
            </h2>
            <span className="text-sm text-slate-400">
              {results.results?.length || 0} gefunden
            </span>
          </div>

          {results.results && results.results.length > 0 ? (
            <div className="space-y-4">
              {results.results.map((result, idx) => {
                const similarity = Math.round(result.similarity * 100);
                const snippetLength = 200;
                const snippet = result.content.length > snippetLength
                  ? result.content.substring(0, snippetLength) + "..."
                  : result.content;

                return (
                  <div
                    key={idx}
                    className="group p-6 bg-slate-800/50 backdrop-blur-sm rounded-xl border border-white/10
                      hover:border-blue-500/30 hover:bg-slate-800/70 transition-all cursor-pointer"
                  >
                    {/* File Header */}
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="p-2 bg-blue-500/20 rounded-lg">
                          <FileText size={20} className="text-blue-400" />
                        </div>
                        <div>
                          <h3 className="font-semibold text-white group-hover:text-blue-400 transition-colors">
                            {result.file_path.replace('/mnt/data/test_corpus/', '')}
                          </h3>
                          <p className="text-xs text-slate-500">
                            {result.file_path}
                          </p>
                        </div>
                      </div>

                      {/* Similarity Score */}
                      <div className="flex flex-col items-end gap-1">
                        <span className="text-sm font-medium text-slate-300">
                          {similarity}% Match
                        </span>
                        <div className="w-24 h-2 bg-slate-700/50 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-gradient-to-r from-blue-500 to-violet-500 transition-all"
                            style={{ width: `${similarity}%` }}
                          />
                        </div>
                      </div>
                    </div>

                    {/* Content Snippet */}
                    <div
                      className="text-sm text-slate-300 leading-relaxed whitespace-pre-wrap"
                      dangerouslySetInnerHTML={{
                        __html: highlightText(snippet, query)
                      }}
                    />
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="text-center py-12 bg-slate-800/50 rounded-xl border border-white/10">
              <p className="text-slate-400 text-lg">
                Keine Ergebnisse gefunden
              </p>
              <p className="text-slate-500 text-sm mt-2">
                Versuchen Sie es mit anderen Suchbegriffen
              </p>
            </div>
          )}
        </div>
      )}

      {/* Initial State */}
      {!results && !loading && (
        <div className="text-center py-16">
          <div className="inline-flex p-4 bg-slate-800/50 rounded-full mb-4">
            <SearchIcon size={48} className="text-slate-500" />
          </div>
          <p className="text-slate-400">
            Geben Sie einen Suchbegriff ein, um zu beginnen
          </p>
          <div className="mt-6 flex flex-wrap justify-center gap-2">
            {["Rechnung Müller", "Server Kosten", "API Fehler", "E-Mail Support"].map(suggestion => (
              <button
                key={suggestion}
                onClick={() => {
                  setQuery(suggestion);
                  performSearch(suggestion);
                }}
                className="px-4 py-2 bg-slate-800/50 border border-white/10 rounded-lg text-sm text-slate-300 hover:border-blue-500/50 hover:text-blue-400 transition-colors"
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
