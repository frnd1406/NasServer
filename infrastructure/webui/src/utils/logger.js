/**
 * FIX [BUG-JS-010]: Production-safe logging utility
 * Only logs in development mode, silent in production
 */

const isDevelopment = import.meta.env.MODE === 'development';

export const logger = {
  error: (...args) => {
    if (isDevelopment) {
      console.error('[ERROR]', ...args);
    }
    // In production, could send to error tracking service here
  },

  warn: (...args) => {
    if (isDevelopment) {
      console.warn('[WARN]', ...args);
    }
  },

  info: (...args) => {
    if (isDevelopment) {
      console.info('[INFO]', ...args);
    }
  },

  debug: (...args) => {
    if (isDevelopment) {
      console.debug('[DEBUG]', ...args);
    }
  },
};

export default logger;
