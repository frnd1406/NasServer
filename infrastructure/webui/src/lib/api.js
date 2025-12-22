const envBaseUrl = (import.meta.env.VITE_API_BASE_URL || "").trim();

function deriveDefaultBaseUrl() {
  if (typeof window === "undefined") {
    return "http://localhost:8080";
  }
  const { protocol, hostname } = window.location;

  // In production (via nginx), use relative URLs so requests go through the proxy
  // Only use explicit port for localhost development
  if (hostname === "localhost" || hostname === "127.0.0.1") {
    const defaultPort = protocol === "https:" ? "8443" : "8080";
    return `${protocol}//${hostname}:${defaultPort}`;
  }

  // Production: use empty base URL (relative paths work via nginx proxy)
  return "";
}

function normalizeBaseUrl(url) {
  const base = url || deriveDefaultBaseUrl();
  return base.replace(/\/+$/, "");
}

const API_BASE_URL = normalizeBaseUrl(envBaseUrl);
const LOGOUT_COUNTDOWN_SECONDS = 4;

function buildUrl(path = "") {
  if (!path.startsWith("/")) {
    return `${API_BASE_URL}/${path}`;
  }
  return `${API_BASE_URL}${path}`;
}

export function getApiBaseUrl() {
  return API_BASE_URL;
}

let logoutOverlay;
let logoutCountdownInterval;
let logoutRedirectScheduled = false;

function ensureLogoutStyles() {
  if (typeof document === "undefined") return;
  if (document.getElementById("session-warning-styles")) return;

  const style = document.createElement("style");
  style.id = "session-warning-styles";
  style.textContent = `
    .session-warning-overlay {
      position: fixed;
      inset: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(0,0,0,0.35);
      backdrop-filter: blur(12px);
      z-index: 9999;
      opacity: 0;
      pointer-events: none;
      transition: opacity 180ms ease;
    }
    .session-warning-overlay.is-visible {
      opacity: 1;
      pointer-events: auto;
    }
    .session-warning-card {
      max-width: 420px;
      width: 90%;
      padding: 24px 26px;
      border-radius: 18px;
      background: linear-gradient(135deg, rgba(239,68,68,0.25), rgba(127,29,29,0.35));
      border: 1px solid rgba(248,113,113,0.55);
      box-shadow: 0 20px 70px rgba(239,68,68,0.35);
      color: #fff;
      backdrop-filter: blur(18px);
      font-family: 'Inter', system-ui, -apple-system, sans-serif;
    }
    .session-warning-pill {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 6px 12px;
      border-radius: 999px;
      background: rgba(248,113,113,0.28);
      border: 1px solid rgba(248,113,113,0.55);
      text-transform: uppercase;
      letter-spacing: 0.08em;
      font-size: 12px;
      font-weight: 700;
    }
    .session-warning-title {
      margin: 14px 0 6px 0;
      font-size: 22px;
      font-weight: 800;
      letter-spacing: 0.01em;
    }
    .session-warning-body {
      margin: 0 0 12px 0;
      color: rgba(255,255,255,0.9);
      line-height: 1.4;
      font-size: 15px;
    }
    .session-warning-timer {
      font-weight: 800;
      font-size: 17px;
      color: #fecdd3;
    }
    .session-warning-subtle {
      margin: 0;
      color: rgba(255,255,255,0.78);
      font-size: 13px;
    }
  `;

  document.head.appendChild(style);
}

function showLogoutOverlay(seconds) {
  if (typeof document === "undefined") return;
  ensureLogoutStyles();

  let overlay = document.getElementById("session-warning-overlay");

  if (!overlay) {
    overlay = document.createElement("div");
    overlay.id = "session-warning-overlay";
    overlay.className = "session-warning-overlay";
    overlay.innerHTML = `
      <div class="session-warning-card">
        <div class="session-warning-pill">Warnung · Session läuft ab</div>
        <div class="session-warning-title">Gleich wirst du abgemeldet</div>
        <p class="session-warning-body">
          Wir konnten deinen Token nicht erneuern. Du wirst in
          <span class="session-warning-timer" data-session-timer></span>
          abgemeldet.
        </p>
        <p class="session-warning-subtle">Bitte melde dich erneut an, um weiterzuarbeiten.</p>
      </div>
    `;
    document.body.appendChild(overlay);
  }

  // Make sure to show it
  requestAnimationFrame(() => {
    overlay.classList.add("is-visible");
  });


  const timerEl = overlay.querySelector("[data-session-timer]");
  if (!timerEl) return;

  let remaining = seconds;
  timerEl.textContent = `${remaining}s`;

  if (logoutCountdownInterval) {
    clearInterval(logoutCountdownInterval);
  }

  logoutCountdownInterval = setInterval(() => {
    remaining -= 1;
    if (remaining <= 0) {
      clearInterval(logoutCountdownInterval);
      return;
    }
    timerEl.textContent = `${remaining}s`;
  }, 1000);
}

function clearAuth() {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("csrfToken");
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("csrf_token");
}

async function refreshAccessToken() {
  const refreshToken = localStorage.getItem("refreshToken");
  if (!refreshToken) return null;

  try {
    const res = await fetch(buildUrl("/auth/refresh"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!res.ok) {
      return null;
    }

    const data = await res.json().catch(() => null);
    const newAccessToken = data?.access_token || data?.token;

    if (!newAccessToken) {
      return null;
    }

    localStorage.setItem("accessToken", newAccessToken);
    if (data?.refresh_token) {
      localStorage.setItem("refreshToken", data.refresh_token);
    }
    if (data?.csrf_token) {
      localStorage.setItem("csrfToken", data.csrf_token);
    }

    return newAccessToken;
  } catch (err) {
    return null;
  }
}

function redirectToLogin() {
  if (typeof window === "undefined") {
    clearAuth();
    return;
  }

  if (logoutRedirectScheduled) return;
  logoutRedirectScheduled = true;

  clearAuth();
  showLogoutOverlay(LOGOUT_COUNTDOWN_SECONDS);

  setTimeout(() => {
    window.location.href = "/login";
  }, LOGOUT_COUNTDOWN_SECONDS * 1000);
}

function buildHeaders(accessToken, headersOverride = {}) {
  const csrfToken = localStorage.getItem("csrfToken") || localStorage.getItem("csrf_token") || "";
  const headers = {
    "Content-Type": "application/json",
    ...headersOverride,
  };

  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }
  if (csrfToken && !headers["X-CSRF-Token"]) {
    headers["X-CSRF-Token"] = csrfToken;
  }

  return headers;
}

function extractErrorMessage(res, data) {
  return data?.error?.message || data?.error || res.statusText || "Request failed";
}

async function performRequest(path, options, tokenOverride) {
  const accessToken = tokenOverride || localStorage.getItem("accessToken");
  let res;
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), options.timeout || 30000);

  try {
    res = await fetch(buildUrl(path), {
      ...options,
      signal: controller.signal,
      credentials: 'include', // <--- WICHTIG: HINZUFÜGEN [Fix für BUG-JS-002]
      headers: buildHeaders(accessToken, options.headers || {}),
    });
  } catch (err) {
    if (err.name === 'AbortError') {
      throw new Error(`Request timed out after ${(options.timeout || 30000) / 1000} seconds`);
    }
    throw new Error(`Cannot reach API at ${API_BASE_URL} (${err.message})`);
  } finally {
    clearTimeout(timeoutId);
  }

  const isJson = res.headers.get("content-type")?.includes("application/json");
  let data = null;

  if (isJson) {
    try {
      data = await res.json();
    } catch (err) {
      data = null;
    }
  }

  return { res, data };
}

export async function apiRequest(path, options = {}) {
  const firstAttempt = await performRequest(path, options);

  if (firstAttempt.res.ok) {
    return firstAttempt.data;
  }

  if (firstAttempt.res.status === 401) {
    const newAccessToken = await refreshAccessToken();

    if (newAccessToken) {
      const retry = await performRequest(path, options, newAccessToken);

      if (retry.res.ok) {
        return retry.data;
      }

      if (retry.res.status === 401) {
        redirectToLogin();
      }

      const retryMessage = extractErrorMessage(retry.res, retry.data);
      const retryError = new Error(retryMessage);
      retryError.status = retry.res.status;
      throw retryError;
    }

    redirectToLogin();
    const refreshError = new Error("Session expired. Please log in again.");
    refreshError.status = 401;
    throw refreshError;
  }

  // GLOBAL VAULT INTERCEPTOR: Handle 423 Locked (Vault is locked)
  // Redirect to vault unlock page immediately
  if (firstAttempt.res.status === 423) {
    if (typeof window !== "undefined" && !window.location.pathname.includes('/vault/unlock')) {
      window.location.href = "/vault/unlock";
    }
    const vaultError = new Error("Vault is locked. Please unlock to access encrypted files.");
    vaultError.status = 423;
    throw vaultError;
  }

  const message = extractErrorMessage(firstAttempt.res, firstAttempt.data);
  const error = new Error(message);
  error.status = firstAttempt.res.status;
  throw error;
}

/**
 * Search for documents using semantic search
 * @param {string} query - Search query
 * @returns {Promise<{query: string, results: Array}>}
 */
export async function searchFiles(query) {
  if (!query || !query.trim()) {
    throw new Error("Search query is required");
  }

  const encodedQuery = encodeURIComponent(query.trim());
  const response = await apiRequest(`/api/v1/search?q=${encodedQuery}`, {
    method: "GET",
  });

  return response;
}

/**
 * Unified AI Query - AI decides whether to search or answer
 * Supports both async (default) and sync (?sync=true) modes
 * 
 * @param {string} query - User query (question or search)
 * @param {Object} options - Optional configuration
 * @param {function} options.onProgress - Progress callback (status, poll count)
 * @param {boolean} options.sync - Force sync mode (default: false)
 * @returns {Promise<{
 *   mode: "search" | "answer",
 *   intent: { type: string, count_hint: string, limit: number },
 *   files?: Array<{ file_id: string, file_path: string, content: string, similarity: number }>,
 *   answer?: string,
 *   sources?: Array<{ file_id: string, file_path: string, similarity: number }>,
 *   confidence?: string,
 *   query: string
 * }>}
 */
export async function queryAI(query, options = {}) {
  if (!query?.trim()) {
    throw new Error("Query is required");
  }

  const { onProgress, sync = false } = options;
  const endpoint = sync ? "/api/v1/query?sync=true" : "/api/v1/query";

  // Step 1: Submit query
  const submitResponse = await apiRequest(endpoint, {
    method: "POST",
    body: JSON.stringify({ query: query.trim() }),
    timeout: sync ? 120000 : 10000, // Longer timeout for sync mode
  });

  // Check if this is a direct response (sync mode or backward compat)
  if (!submitResponse.job_id) {
    // Direct response - return as-is
    return submitResponse;
  }

  // Async mode - poll for result
  const jobId = submitResponse.job_id;
  const maxPolls = 80; // 80 * 1.5s = 120s max (matches JOB_TIMEOUT)
  const pollInterval = 1500; // 1.5 seconds

  if (onProgress) {
    onProgress({ status: "pending", poll: 0, jobId });
  }

  for (let i = 0; i < maxPolls; i++) {
    await new Promise(resolve => setTimeout(resolve, pollInterval));

    try {
      const statusResponse = await apiRequest(`/api/v1/jobs/${jobId}`, {
        timeout: 5000,
      });

      if (onProgress) {
        onProgress({
          status: statusResponse.status,
          poll: i + 1,
          jobId,
          elapsed: Math.round((i + 1) * pollInterval / 1000)
        });
      }

      if (statusResponse.status === "completed") {
        return statusResponse;
      }

      if (statusResponse.status === "failed") {
        throw new Error(statusResponse.error || "AI processing failed");
      }

      // Still pending or processing - continue polling
    } catch (err) {
      // If job not found yet, might be a race condition - keep trying
      if (err.status === 404 && i < 3) {
        continue;
      }
      throw err;
    }
  }

  throw new Error("AI processing timed out - please try again");
}


/**
 * Download multiple files/folders as a ZIP
 * @param {string[]} paths - Array of file/folder paths to download
 * @returns {Promise<Blob>} ZIP file blob
 */
export async function batchDownload(paths) {
  if (!paths || paths.length === 0) {
    throw new Error("No files selected for download");
  }

  const accessToken = localStorage.getItem("accessToken");
  const csrfToken = localStorage.getItem("csrfToken") || localStorage.getItem("csrf_token") || "";

  const res = await fetch(buildUrl("/api/v1/storage/batch-download"), {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${accessToken}`,
      "X-CSRF-Token": csrfToken,
    },
    body: JSON.stringify({ paths }),
  });

  if (!res.ok) {
    throw new Error(`Batch download failed: ${res.status}`);
  }

  return await res.blob();
}

/**
 * Download a folder as a ZIP file
 * @param {string} path - Path to the folder
 * @returns {Promise<Blob>} ZIP file blob
 */
export async function downloadFolderAsZip(path) {
  if (!path) {
    throw new Error("Folder path is required");
  }

  const accessToken = localStorage.getItem("accessToken");
  const csrfToken = localStorage.getItem("csrfToken") || localStorage.getItem("csrf_token") || "";

  const res = await fetch(buildUrl(`/api/v1/storage/download-zip?path=${encodeURIComponent(path)}`), {
    method: "GET",
    credentials: "include",
    headers: {
      Authorization: `Bearer ${accessToken}`,
      "X-CSRF-Token": csrfToken,
    },
  });

  if (!res.ok) {
    throw new Error(`ZIP download failed: ${res.status}`);
  }

  return await res.blob();
}
