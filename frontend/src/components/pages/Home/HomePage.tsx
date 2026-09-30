import { ShoppingBag, Mail, Search, Trash2, Plus } from 'lucide-react';
import { Container } from '@/components/layout/Container';
import {
  Button,
  IconButton,
  Input,
  Textarea,
  Card,
  Badge,
  Spinner,
  Skeleton,
  EmptyState,
} from '@/components/ui';
import { toast } from 'sonner';

export function HomePage() {
  return (
    <Container>
      <div className="py-12 space-y-16 max-w-4xl">
        <header>
          <h1 className="text-3xl font-bold">Component Gallery</h1>
          <p className="mt-2 text-text-secondary">
            Temporary page for Phase 2A review. Will be replaced in Phase 3.
          </p>
        </header>

        {/* ---------------- BUTTONS ---------------- */}
        <Section title="Buttons">
          <div className="flex flex-wrap gap-3 items-center">
            <Button>Primary</Button>
            <Button variant="outline">Outline</Button>
            <Button variant="ghost">Ghost</Button>
            <Button variant="danger">Danger</Button>
          </div>

          <div className="flex flex-wrap gap-3 items-center">
            <Button size="sm">Small</Button>
            <Button size="md">Medium</Button>
            <Button size="lg">Large</Button>
          </div>

          <div className="flex flex-wrap gap-3 items-center">
            <Button leftIcon={<Plus size={16} />}>Add item</Button>
            <Button rightIcon={<ShoppingBag size={16} />} variant="outline">
              View cart
            </Button>
            <Button isLoading>Loading</Button>
            <Button disabled>Disabled</Button>
          </div>

          <div className="max-w-xs">
            <Button fullWidth>Full width</Button>
          </div>
        </Section>

        {/* ---------------- ICON BUTTONS ---------------- */}
        <Section title="Icon Buttons">
          <div className="flex flex-wrap gap-3 items-center">
            <IconButton aria-label="Cart" size="sm">
              <ShoppingBag size={16} />
            </IconButton>
            <IconButton aria-label="Cart" size="md">
              <ShoppingBag size={20} />
            </IconButton>
            <IconButton aria-label="Cart" size="lg">
              <ShoppingBag size={22} />
            </IconButton>
          </div>

          <div className="flex flex-wrap gap-3 items-center">
            <IconButton aria-label="Delete" variant="ghost">
              <Trash2 size={20} />
            </IconButton>
            <IconButton aria-label="Cart" variant="primary">
              <ShoppingBag size={20} />
            </IconButton>
          </div>
        </Section>

        {/* ---------------- INPUTS ---------------- */}
        <Section title="Inputs">
          <div className="grid gap-4 max-w-md">
            <Input label="Email" placeholder="you@example.com" />
            <Input
              label="Search"
              placeholder="Find a dish..."
              leftIcon={<Search size={16} />}
            />
            <Input
              label="Email"
              placeholder="you@example.com"
              leftIcon={<Mail size={16} />}
              error="This email is already in use"
            />
            <Input
              label="Phone"
              placeholder="+91 98765 43210"
              helperText="We'll only use this for delivery updates"
            />
            <Input label="Disabled" placeholder="Can't type here" disabled />
            <Textarea
              label="Special instructions"
              placeholder="Any allergies or preferences?"
              helperText="Optional"
            />
          </div>
        </Section>

        {/* ---------------- CARDS ---------------- */}
        <Section title="Cards">
          <div className="grid md:grid-cols-2 gap-4">
            <Card className="p-6">
              <h4 className="font-semibold">Default card</h4>
              <p className="text-sm text-text-secondary mt-1">
                With border and subtle shadow.
              </p>
            </Card>
            <Card variant="interactive" className="p-6" onClick={() => toast.success('Card clicked!')}>
              <h4 className="font-semibold">Interactive card</h4>
              <p className="text-sm text-text-secondary mt-1">
                Hover me — I lift up and cast a bigger shadow.
              </p>
            </Card>
          </div>
        </Section>

        {/* ---------------- BADGES ---------------- */}
        <Section title="Badges">
          <div className="flex flex-wrap gap-2">
            <Badge>Neutral</Badge>
            <Badge variant="primary">New</Badge>
            <Badge variant="success">Delivered</Badge>
            <Badge variant="warning">Pending</Badge>
            <Badge variant="danger">Cancelled</Badge>
            <Badge variant="info">Info</Badge>
          </div>
          <div className="flex flex-wrap gap-2 items-center">
            <Badge size="md" variant="success">
              Delivered
            </Badge>
            <Badge size="md" variant="warning">
              Preparing
            </Badge>
          </div>
        </Section>

        {/* ---------------- LOADERS ---------------- */}
        <Section title="Loaders">
          <div className="flex items-center gap-6">
            <Spinner size="sm" />
            <Spinner size="md" />
            <Spinner size="lg" />
          </div>
          <div className="space-y-3 max-w-md">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
            <Skeleton className="h-32 w-full" />
          </div>
        </Section>

        {/* ---------------- EMPTY STATE ---------------- */}
        <Section title="Empty State">
          <Card>
            <EmptyState
              icon={<ShoppingBag size={28} />}
              title="Your cart is empty"
              description="Add some delicious dishes to get started."
              action={<Button>Browse menu</Button>}
            />
          </Card>
        </Section>

        {/* ---------------- TOASTS ---------------- */}
        <Section title="Toasts">
          <div className="flex flex-wrap gap-3">
            <Button
              variant="outline"
              onClick={() => toast.success('Order placed successfully!')}
            >
              Success toast
            </Button>
            <Button
              variant="outline"
              onClick={() => toast.error('Something went wrong')}
            >
              Error toast
            </Button>
            <Button
              variant="outline"
              onClick={() => toast.info('Pickup time updated')}
            >
              Info toast
            </Button>
          </div>
        </Section>
      </div>
    </Container>
  );
}

// ---------------------------------------------------------------------------
// Helper component for gallery sections
// ---------------------------------------------------------------------------
function Section({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section>
      <h2 className="text-xl font-semibold mb-4 pb-2 border-b border-border">
        {title}
      </h2>
      <div className="space-y-4">{children}</div>
    </section>
  );
}