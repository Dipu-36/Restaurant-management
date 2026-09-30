import { create } from 'zustand';
import { apiClient } from '../api/client';
import type { Restaurant, Category, Dish } from '../types/api.types';

interface RestaurantState {
  info: Restaurant | null;
  categories: Category[];
  dishes: Dish[];
  featuredDishes: Dish[];
  isLoading: boolean;
  isInitialized: boolean;
  error: string | null;

  // Actions
  fetchRestaurantInfo: () => Promise<void>;
  fetchCategories: () => Promise<void>;
  fetchDishes: (categoryId?: number) => Promise<void>;
  fetchFeaturedDishes: () => Promise<void>;
  initialize: () => Promise<void>;
}

export const useRestaurantStore = create<RestaurantState>((set) => ({
  info: null,
  categories: [],
  dishes: [],
  featuredDishes: [],
  isLoading: false,
  isInitialized: false,
  error: null,

  fetchRestaurantInfo: async () => {
    try {
      const response = await apiClient.get<{ restaurant: Restaurant }>(
        '/restaurant'
      );
      set({ info: response.data.restaurant });
    } catch (error) {
      // Restaurant info might be owner-only; don't crash
      console.error('Failed to fetch restaurant info:', error);
    }
  },

  fetchCategories: async () => {
    try {
      // Note: Backend currently requires owner role for this.
      // We'll adjust when we make the public endpoint.
      const response = await apiClient.get<{ categories: Category[] }>(
        '/categories'
      );
      set({ categories: response.data.categories });
    } catch (error) {
      console.error('Failed to fetch categories:', error);
    }
  },

  fetchDishes: async (categoryId) => {
    set({ isLoading: true });
    try {
      const params: Record<string, string> = { available_only: 'true' };
      if (categoryId) params.category_id = String(categoryId);

      const response = await apiClient.get<{ dishes: Dish[] }>('/dishes', {
        params,
      });
      set({ dishes: response.data.dishes, isLoading: false });
    } catch (error) {
      set({ isLoading: false });
      console.error('Failed to fetch dishes:', error);
    }
  },

  fetchFeaturedDishes: async () => {
    try {
      const response = await apiClient.get<{ dishes: Dish[] }>('/dishes', {
        params: { available_only: 'true' },
      });
      const featured = response.data.dishes.filter((d) => d.is_featured);
      set({ featuredDishes: featured });
    } catch (error) {
      console.error('Failed to fetch featured dishes:', error);
    }
  },

  initialize: async () => {
    // Fetch everything in parallel
    await Promise.all([
      useRestaurantStore.getState().fetchRestaurantInfo(),
      useRestaurantStore.getState().fetchCategories(),
      useRestaurantStore.getState().fetchDishes(),
    ]);
    set({ isInitialized: true });
  },
}));