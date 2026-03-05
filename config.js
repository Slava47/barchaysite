// API Configuration for local development.
// In Docker the entrypoint overwrites this file with the real API_BASE value.
// Leave API_BASE empty to use built-in MENU_DATA / QUIZ_DATA (offline mode).
window.__CONFIG__ = {
  API_BASE: ''
};
