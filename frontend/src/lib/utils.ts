import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Merge Tailwind classes intelligently.
 * Handles conditional classes and resolves conflicts.
 * 
 * Example:
 *   cn('px-2 py-1', 'px-4') // => 'py-1 px-4'
 *   cn('text-red-500', isActive && 'text-blue-500') // => conditional
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}