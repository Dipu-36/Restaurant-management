import { create } from 'zustand';
import { apiClient } from '../api/client';
import type {
  User,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  ActivateRequest,
  ActivateResponse,
} from '../types/api.types';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isOwner: boolean;
  isCustomer: boolean;
  isLoading: boolean;

  // Actions
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<User>;
  activate: (data: ActivateRequest) => Promise<void>;
  logout: () => void;
  initializeAuth: () => void;
}

const TOKEN_KEY = 'foodhub_token';
const USER_KEY = 'foodhub_user';

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isAuthenticated: false,
  isOwner: false,
  isCustomer: false,
  isLoading: false,

  /**
   * Called once on app mount to restore session from localStorage
   */
  initializeAuth: () => {
    const token = localStorage.getItem(TOKEN_KEY);
    const userJson = localStorage.getItem(USER_KEY);

    if (token && userJson) {
      try {
        const user = JSON.parse(userJson) as User;
        set({
          user,
          token,
          isAuthenticated: true,
          isOwner: user.role === 'owner',
          isCustomer: user.role === 'customer',
        });
      } catch {
        // Corrupted storage, clear it
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
      }
    }
  },

  /**
   * Login: POST /v1/tokens/authentication
   * Note: The backend only returns the token, not the user.
   * We'll decode the JWT or make a second call to get user info.
   */
  login: async (credentials) => {
    set({ isLoading: true });
    try {
      const response = await apiClient.post<LoginResponse>(
        '/tokens/authentication',
        credentials
      );

      const { authentication_token } = response.data;

      // Decode JWT to extract user info
      const payload = decodeJwt(authentication_token);

      const user: User = {
        id: payload.user_id,
        role: payload.role,
        // These fields aren't in the JWT, we'll fetch them later or use minimal data
        created_at: '',
        name: '',
        email: credentials.email,
        phone: '',
        activated: true,
        version: 0,
      };

      localStorage.setItem(TOKEN_KEY, authentication_token);
      localStorage.setItem(USER_KEY, JSON.stringify(user));

      set({
        user,
        token: authentication_token,
        isAuthenticated: true,
        isOwner: user.role === 'owner',
        isCustomer: user.role === 'customer',
        isLoading: false,
      });
    } catch (error) {
      set({ isLoading: false });
      throw error;
    }
  },

  /**
   * Register: POST /v1/users
   * Returns the created user (not activated yet).
   */
  register: async (data) => {
    set({ isLoading: true });
    try {
      const response = await apiClient.post<RegisterResponse>('/users', data);
      set({ isLoading: false });
      return response.data.user;
    } catch (error) {
      set({ isLoading: false });
      throw error;
    }
  },

  /**
   * Activate: PUT /v1/users/activated
   * After activation, the user can log in.
   */
  activate: async (data) => {
    set({ isLoading: true });
    try {
      await apiClient.put<ActivateResponse>('/users/activated', data);
      set({ isLoading: false });
    } catch (error) {
      set({ isLoading: false });
      throw error;
    }
  },

  /**
   * Logout: clear state and storage
   */
  logout: () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    set({
      user: null,
      token: null,
      isAuthenticated: false,
      isOwner: false,
      isCustomer: false,
    });
  },
}));

// ---------------------------------------------------------------------------
// JWT decoder (no external library needed)
// ---------------------------------------------------------------------------
interface JwtPayload {
  user_id: number;
  role: 'customer' | 'owner';
  iss: string;
  sub: string;
  exp: number;
  iat: number;
  nbf: number;
}

function decodeJwt(token: string): JwtPayload {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch {
    throw new Error('Invalid token');
  }
}