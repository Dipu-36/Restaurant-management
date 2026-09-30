import { Link } from 'react-router-dom';
import { Container } from '../../layout/Container';

export function NotFoundPage() {
  return (
    <Container>
      <div className="py-24 text-center">
        <h1 className="text-6xl font-bold text-primary">404</h1>
        <h2 className="mt-4 text-2xl font-semibold">Page not found</h2>
        <p className="mt-2 text-text-secondary">
          The page you're looking for doesn't exist.
        </p>
        <Link
          to="/"
          className="mt-6 inline-flex items-center justify-center h-11 px-6 rounded-full bg-primary text-white font-medium hover:opacity-90 transition-opacity"
        >
          Back to Home
        </Link>
      </div>
    </Container>
  );
}