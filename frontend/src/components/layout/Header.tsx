import { Link, NavLink } from 'react-router-dom';
import { ShoppingBag, UserCircle, Menu as MenuIcon } from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCartStore } from '../../store/useCartStore';
import { useAuthStore } from '../../store/useAuthStore';
import { APP_NAME } from '../../lib/constants';

export function Header() {
  const itemCount = useCartStore((s) => s.getItemCount());
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return (
    <header className="sticky top-0 z-40 bg-background/90 backdrop-blur-md border-b border-border">
      <div className="container">
        <div className="grid h-16 grid-cols-[1fr_auto_1fr] items-center gap-4">
          {/* Left: Navigation (Desktop) */}
          <nav className="hidden md:flex items-center gap-8">
            <HeaderNavLink to="/">Menu</HeaderNavLink>
            <HeaderNavLink to="/orders">Orders</HeaderNavLink>
            <HeaderNavLink to="/profile">Account</HeaderNavLink>
          </nav>

          {/* Mobile Menu Button */}
          <button
            type="button"
            aria-label="Open menu"
            className="md:hidden flex h-10 w-10 items-center justify-center rounded-full hover:bg-muted transition-colors"
          >
            <MenuIcon size={22} />
          </button>

          {/* Center: Logo */}
          <Link
            to="/"
            className="justify-self-center text-2xl font-bold tracking-tight text-primary hover:opacity-80 transition-opacity"
          >
            {APP_NAME}
          </Link>

          {/* Right: Actions */}
          <div className="flex items-center justify-end gap-2">
            <Link
              to="/menu"
              className={cn(
                'hidden md:inline-flex items-center justify-center',
                'h-10 px-5 rounded-full',
                'text-sm font-medium',
                'border border-border bg-transparent',
                'transition-colors hover:bg-muted',
                'cursor-pointer'
              )}
            >
              Order
            </Link>

            {/* Cart */}
            <Link
              to="/cart"
              aria-label="Cart"
              className="relative flex h-10 w-10 items-center justify-center rounded-full hover:bg-muted transition-colors"
            >
              <ShoppingBag size={22} />
              {itemCount > 0 && (
                <span className="absolute -top-0.5 -right-0.5 flex h-5 w-5 items-center justify-center rounded-full bg-primary text-white text-xs font-semibold">
                  {itemCount > 9 ? '9+' : itemCount}
                </span>
              )}
            </Link>

            {/* Profile / Login */}
            <Link
              to={isAuthenticated ? '/profile' : '/login'}
              aria-label={isAuthenticated ? 'Profile' : 'Login'}
              className="flex h-10 w-10 items-center justify-center rounded-full hover:bg-muted transition-colors"
            >
              <UserCircle size={24} />
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
}

// ---------------------------------------------------------------------------
// Sub-component: NavLink with underline animation on hover
// ---------------------------------------------------------------------------
interface HeaderNavLinkProps {
  to: string;
  children: React.ReactNode;
}

function HeaderNavLink({ to, children }: HeaderNavLinkProps) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        cn(
          'relative text-sm font-medium text-foreground transition-colors',
          'cursor-pointer',
          'hover:text-primary',
          'py-1',
          // Underline animation
          'after:absolute after:left-0 after:bottom-0 after:h-[2px] after:w-0 after:bg-primary',
          'after:transition-all after:duration-300',
          'hover:after:w-full',
          isActive && 'text-primary after:w-full'
        )
      }
    >
      {children}
    </NavLink>
  );
}