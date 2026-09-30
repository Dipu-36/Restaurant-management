import { Outlet } from 'react-router-dom';
import { Toaster } from 'sonner';
import { Header } from './Header';
import { LocationBar } from './LocationBar';

export function RootLayout() {
  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Header />
      <LocationBar />

      <main className="flex-1">
        <Outlet />
      </main>

      <Toaster
        position="top-center"
        toastOptions={{
          style: {
            background: '#ffffff',
            border: '1px solid #e8e3dc',
            color: '#1a1a1a',
          },
        }}
      />
    </div>
  );
}