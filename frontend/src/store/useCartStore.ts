import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Dish } from '../types/api.types';

export interface CartItem {
  dish: Dish;
  quantity: number;
}

interface CartState {
  items: CartItem[];

  // Actions
  addItem: (dish: Dish, quantity?: number) => void;
  removeItem: (dishId: number) => void;
  updateQuantity: (dishId: number, quantity: number) => void;
  incrementQuantity: (dishId: number) => void;
  decrementQuantity: (dishId: number) => void;
  clearCart: () => void;

  // Computed getters (as functions since Zustand doesn't support computed directly)
  getSubtotal: () => number;
  getItemCount: () => number;
  getItemQuantity: (dishId: number) => number;
  isEmpty: () => boolean;
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],

      addItem: (dish, quantity = 1) => {
        set((state) => {
          const existing = state.items.find((i) => i.dish.id === dish.id);
          if (existing) {
            return {
              items: state.items.map((i) =>
                i.dish.id === dish.id
                  ? { ...i, quantity: i.quantity + quantity }
                  : i
              ),
            };
          }
          return { items: [...state.items, { dish, quantity }] };
        });
      },

      removeItem: (dishId) => {
        set((state) => ({
          items: state.items.filter((i) => i.dish.id !== dishId),
        }));
      },

      updateQuantity: (dishId, quantity) => {
        if (quantity <= 0) {
          get().removeItem(dishId);
          return;
        }
        set((state) => ({
          items: state.items.map((i) =>
            i.dish.id === dishId ? { ...i, quantity } : i
          ),
        }));
      },

      incrementQuantity: (dishId) => {
        const item = get().items.find((i) => i.dish.id === dishId);
        if (item) {
          get().updateQuantity(dishId, item.quantity + 1);
        }
      },

      decrementQuantity: (dishId) => {
        const item = get().items.find((i) => i.dish.id === dishId);
        if (item) {
          get().updateQuantity(dishId, item.quantity - 1);
        }
      },

      clearCart: () => set({ items: [] }),

      getSubtotal: () => {
        return get().items.reduce(
          (sum, item) => sum + item.dish.price * item.quantity,
          0
        );
      },

      getItemCount: () => {
        return get().items.reduce((sum, item) => sum + item.quantity, 0);
      },

      getItemQuantity: (dishId) => {
        const item = get().items.find((i) => i.dish.id === dishId);
        return item?.quantity ?? 0;
      },

      isEmpty: () => get().items.length === 0,
    }),
    {
      name: 'foodhub_cart',
    }
  )
);