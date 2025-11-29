import { useState } from "react";
import { Search as SearchIcon, FileText, Loader2 } from "lucide-react";
import { searchFiles } from "../lib/api";

export default function Search() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSearch = async (e) => {
    e.preventDefault();

    if (!query.trim()) {
      setError("Bitte geben Sie einen Suchbegriff ein");
      return;
    }

    setLoading(true);
    setError("");
    setResults(null);

    try {
      const data = await searchFiles(query);
      setResults(data);
    } catch (err) {
      setError(err.message || "Fehler bei der Suche");
      console.error("Search error:", err);
    } finally {
      setLoading(false);
    }
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
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800">
      <div className="max-w-4xl mx-auto px-4 py-12">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold text-slate-900 dark:text-white mb-3">
            Dokumentensuche
          </h1>
          <p className="text-slate-600 dark:text-slate-400">
            Semantische Suche mit KI-gestützter Vektoranalyse
          </p>
        </div>

        {/* Search Box */}
        <form onSubmit={handleSearch} className="mb-8">
          <div className="relative">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Suchen Sie nach Rechnungen, E-Mails, Logs..."
              className="w-full px-6 py-4 pr-14 rounded-2xl border-2 border-slate-200 dark:border-slate-700
                bg-white dark:bg-slate-800 text-slate-900 dark:text-white
                focus:border-blue-500 dark:focus:border-blue-400 focus:outline-none focus:ring-4 focus:ring-blue-500/20
                text-lg transition-all shadow-lg"
              disabled={loading}
            />
            <button
              type="submit"
              disabled={loading}
              className="absolute right-2 top-1/2 -translate-y-1/2 p-3 rounded-xl
                bg-blue-600 hover:bg-blue-700 text-white
                disabled:bg-slate-400 disabled:cursor-not-allowed
                transition-colors shadow-md"
            >
              {loading ? (
                <Loader2 size={24} className="animate-spin" />
              ) : (
                <SearchIcon size={24} />
              )}
            </button>
          </div>

          {error && (
            <div className="mt-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl">
              <p className="text-red-800 dark:text-red-200">{error}</p>
            </div>
          )}
        </form>

        {/* Results */}
        {results && (
          <div className="space-y-4">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-semibold text-slate-900 dark:text-white">
                Ergebnisse für "{results.query}"
              </h2>
              <span className="text-sm text-slate-600 dark:text-slate-400">
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
                      className="group p-6 bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700
                        hover:shadow-xl hover:border-blue-300 dark:hover:border-blue-600 transition-all cursor-pointer"
                    >
                      {/* File Header */}
                      <div className="flex items-start justify-between mb-3">
                        <div className="flex items-center gap-3">
                          <div className="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                            <FileText size={20} className="text-blue-600 dark:text-blue-400" />
                          </div>
                          <div>
                            <h3 className="font-semibold text-slate-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                              {result.file_path.replace('/mnt/data/test_corpus/', '')}
                            </h3>
                            <p className="text-xs text-slate-500 dark:text-slate-400">
                              {result.file_path}
                            </p>
                          </div>
                        </div>

                        {/* Similarity Score */}
                        <div className="flex flex-col items-end gap-1">
                          <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
                            {similarity}% Übereinstimmung
                          </span>
                          <div className="w-24 h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
                            <div
                              className="h-full bg-gradient-to-r from-green-500 to-blue-500 transition-all"
                              style={{ width: `${similarity}%` }}
                            />
                          </div>
                        </div>
                      </div>

                      {/* Content Snippet */}
                      <div
                        className="text-sm text-slate-600 dark:text-slate-300 leading-relaxed whitespace-pre-wrap"
                        dangerouslySetInnerHTML={{
                          __html: highlightText(snippet, query)
                        }}
                      />
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="text-center py-12 bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700">
                <p className="text-slate-600 dark:text-slate-400 text-lg">
                  Keine Ergebnisse gefunden
                </p>
                <p className="text-slate-500 dark:text-slate-500 text-sm mt-2">
                  Versuchen Sie es mit anderen Suchbegriffen
                </p>
              </div>
            )}
          </div>
        )}

        {/* Initial State */}
        {!results && !loading && (
          <div className="text-center py-16">
            <div className="inline-flex p-4 bg-slate-200 dark:bg-slate-800 rounded-full mb-4">
              <SearchIcon size={48} className="text-slate-400" />
            </div>
            <p className="text-slate-600 dark:text-slate-400">
              Geben Sie einen Suchbegriff ein, um zu beginnen
            </p>
            <div className="mt-6 flex flex-wrap justify-center gap-2">
              {["Rechnung Müller", "Server Kosten", "API Fehler", "E-Mail Support"].map(suggestion => (
                <button
                  key={suggestion}
                  onClick={() => {
                    setQuery(suggestion);
                    setTimeout(() => {
                      document.querySelector('form').dispatchEvent(
                        new Event('submit', { cancelable: true, bubbles: true })
                      );
                    }, 100);
                  }}
                  className="px-4 py-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700
                    rounded-lg text-sm text-slate-700 dark:text-slate-300
                    hover:border-blue-500 hover:text-blue-600 dark:hover:text-blue-400
                    transition-colors"
                >
                  {suggestion}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
