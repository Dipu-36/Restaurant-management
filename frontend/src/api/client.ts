import axios, { type AxiosError, type AxiosInstance } from 'axios';
import { toast } from 'sonner';
import type { ApiError } from '../types/api.types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:4000/v1';

export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 15000,
});

// ---------------------------------------------------------------------------
// Request Interceptor: Attach JWT token
// ---------------------------------------------------------------------------
apiClient.interceptors.request.use(
  (config) => {
    // We'll read the token from localStorage to avoid circular deps with store
    const token = localStorage.getItem('foodhub_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// ---------------------------------------------------------------------------
// Response Interceptor: Handle errors globally
// ---------------------------------------------------------------------------
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiError>) => {
    // Network error (no response from server)
    if (!error.response) {
      toast.error('Network error. Please check your connection.');
      return Promise.reject(error);
    }

    const { status, data } = error.response;

    // Unauthorized — clear token and redirect to login
    if (status === 401) {
      localStorage.removeItem('foodhub_token');
      localStorage.removeItem('foodhub_user');

      // Only redirect if not already on auth pages
      const currentPath = window.location.pathname;
      if (
        !currentPath.startsWith('/login') &&
        !currentPath.startsWith('/register') &&
        !currentPath.startsWith('/activate')
      ) {
        window.location.href = '/login';
      }
      return Promise.reject(error);
    }

    // Validation errors (422) — let callers handle these
    if (status === 422) {
      return Promise.reject(error);
    }

    // Forbidden (403)
    if (status === 403) {
      toast.error('You do not have permission to perform this action.');
      return Promise.reject(error);
    }

    // Not found (404)
    if (status === 404) {
      // Let callers handle 404s silently
      return Promise.reject(error);
    }

    // Generic server error (500+)
    if (status >= 500) {
      toast.error('Server error. Please try again later.');
      return Promise.reject(error);
    }

    return Promise.reject(error);
  }
);

/**
 * Helper to extract error message from backend error response.
 * The backend can return:
 *   { "error": "some string" }
 *   { "error": { "field": "message", ... } }  // validation errors
 */
export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError<ApiError>(error)) {
    const data = error.response?.data;
    if (!data) return 'Something went wrong.';

    if (typeof data.error === 'string') {
      return data.error;
    }

    if (typeof data.error === 'object' && data.error !== null) {
      // Validation errors — join first few
      const messages = Object.values(data.error);
      return messages.slice(0, 2).join(', ');
    }
  }
  return 'Something went wrong.';
}

/**
 * Helper to extract field-level validation errors.
 * Returns an object like { email: "must be a valid email address" }
 */
export function getFieldErrors(error: unknown): Record<string, string> {
  if (axios.isAxiosError<ApiError>(error)) {
    const data = error.response?.data;
    if (data && typeof data.error === 'object' && data.error !== null) {
      return data.error;
    }
  }
  return {};
}