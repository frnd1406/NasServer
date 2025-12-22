export function getAuth() {
  const accessToken = localStorage.getItem("accessToken") || localStorage.getItem("access_token") || "";
  const refreshToken = localStorage.getItem("refreshToken") || localStorage.getItem("refresh_token") || "";
  const csrfToken = localStorage.getItem("csrfToken") || localStorage.getItem("csrf_token") || "";
  return {
    accessToken,
    refreshToken,
    csrfToken,
  };
}

export function authHeaders() {
  const { accessToken, csrfToken } = getAuth();
  const headers = {};
  if (accessToken) headers["Authorization"] = `Bearer ${accessToken}`;
  if (csrfToken) headers["X-CSRF-Token"] = csrfToken;
  return headers;
}

export function setAuth({ accessToken = "", refreshToken = "", csrfToken = "" }) {
  if (accessToken) {
    localStorage.setItem("accessToken", accessToken);
    localStorage.setItem("access_token", accessToken);
  }
  if (refreshToken) {
    localStorage.setItem("refreshToken", refreshToken);
    localStorage.setItem("refresh_token", refreshToken);
  }
  if (csrfToken) {
    localStorage.setItem("csrfToken", csrfToken);
    localStorage.setItem("csrf_token", csrfToken);
  }
}

export function clearAuth() {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("csrfToken");
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("csrf_token");
}

export function isAuthenticated() {
  const { accessToken } = getAuth();
  return Boolean(accessToken);
}
