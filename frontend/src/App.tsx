import { RouterProvider } from 'react-router-dom';
import { router } from './routes/router';
import { useAuthStore } from './store/useAuthStore';
import { useRestaurantStore } from './store/useRestaurantStore';
import { useEffect } from 'react';

export default function App() {
  const initializeAuth = useAuthStore((s) => s.initializeAuth);
  const initializeRestaurant = useRestaurantStore((s) => s.initialize);

  useEffect(() => {
    // Restore user session from localStorage
    initializeAuth();
    // Load restaurant data
    initializeRestaurant();
  }, [initializeAuth, initializeRestaurant]);

  return <RouterProvider router={router} />;
}