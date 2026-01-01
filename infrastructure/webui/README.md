# NAS AI Web Interface

A modern, responsive web client built with **React 18**, **Vite**, and **TailwindCSS**. It adheres to a "Glassmorphism" design aesthetic for a premium user experience.

## ✨ Key Features

*   **Virtual Scroll**: Efficient rendering of large file lists using `react-window`.
*   **Smart Uploads**:
    *   Drag & Drop support.
    *   Client-side encryption preparation.
    *   Chunked uploads for large files.
*   **Security**:
    *   Vault locking/unlocking UI.
    *   Visual indicators for encrypted files.
*   **Responsive Design**: Mobile-friendly layout with glass-effect components.

## 📂 Directory Structure

*   `src/components/`: Reusable UI components (Buttons, Modals, File Cards).
*   `src/hooks/`: Custom React hooks (`useFileStorage`, `useFileSelection`) for logic encapsulation.
*   `src/pages/`: Main application views (Files, Dashboard, Settings).
*   `src/context/`: Global state management (`VaultContext`, `AuthContext`).
*   `src/lib/`: Utility libraries and API clients.

## 🚀 Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```
