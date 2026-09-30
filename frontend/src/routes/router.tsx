import { createBrowserRouter } from 'react-router-dom';

// Layouts
import { RootLayout } from '../components/layout/RootLayout';

// Pages - we'll create placeholder files now
import { HomePage } from '../components/pages/Home/HomePage';
import { MenuPage } from '../components/pages/Menu/MenuPage';
import { NotFoundPage } from '../components/pages/NotFound/NotFoundPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: 'menu',
        element: <MenuPage />,
      },
      // These will be added in later phases:
      // { path: 'cart', element: <CartPage /> },
      // { path: 'checkout', element: <CheckoutPage /> },
      // { path: 'login', element: <LoginPage /> },
      // { path: 'register', element: <RegisterPage /> },
      // { path: 'profile', element: <ProfilePage /> },
      // { path: 'orders', element: <OrdersPage /> },
      // { path: 'admin', element: <AdminDashboard /> },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
]);