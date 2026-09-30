import { MapPin, Clock } from 'lucide-react';
import { useRestaurantStore } from '../../store/useRestaurantStore';

export function LocationBar() {
  const info = useRestaurantStore((s) => s.info);

  // Fallback values while restaurant info loads
  const location = info?.street_address ?? 'Loading location...';
  const closingTime = info?.closing_time ?? '10:00pm';
  const isOpen = info?.is_open ?? true;

  return (
    <div className="border-b border-border bg-surface">
      <div className="container">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-2 py-2.5 text-sm">
          {/* Pickup Info */}
          <div className="flex items-center gap-2 text-foreground">
            <MapPin size={16} className="text-primary shrink-0" />
            <span className="font-medium">Pickup from</span>
            <span className="text-text-secondary truncate max-w-[200px] md:max-w-none">
              {location}
            </span>
            <span className="text-text-tertiary">·</span>
            <Clock size={14} className="text-text-tertiary shrink-0" />
            {isOpen ? (
              <span className="text-text-secondary">
                Open until {closingTime}
              </span>
            ) : (
              <span className="text-red-600 font-medium">Closed</span>
            )}
          </div>

          {/* Actions */}
          <div className="flex items-center gap-4 text-sm">
            <button
              type="button"
              className="text-primary font-medium hover:underline cursor-pointer"
            >
              Change location
            </button>
            <button
              type="button"
              className="text-primary font-medium hover:underline cursor-pointer"
            >
              More info
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}