/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Brand colors (Sweetgreen-inspired)
        primary: {
          DEFAULT: '#0a6b4d',
          light: '#1a8a5a',
          dark: '#08503a',
        },
        // Backgrounds
        background: '#faf8f5',
        surface: '#ffffff',
        // Borders & neutrals
        border: '#e8e3dc',
        muted: '#f0ede8',
        // Text
        foreground: '#1a1a1a',
        'text-secondary': '#5a5a5a',
        'text-tertiary': '#8a8a8a',
        // Status colors (for order statuses)
        status: {
          pending: '#f59e0b',
          accepted: '#3b82f6',
          preparing: '#8b5cf6',
          ready: '#10b981',
          out_for_delivery: '#06b6d4',
          delivered: '#22c55e',
          cancelled: '#ef4444',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      container: {
        center: true,
        padding: '1rem',
        screens: {
          '2xl': '1280px',
        },
      },
      boxShadow: {
        'soft': '0 2px 8px rgba(0, 0, 0, 0.04)',
        'medium': '0 4px 16px rgba(0, 0, 0, 0.08)',
      },
      animation: {
        'fade-in': 'fadeIn 0.2s ease-in-out',
        'slide-in': 'slideIn 0.3s ease-out',
        'shimmer': 'shimmer 1.5s infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
          shimmer: {'100%': { transform: 'translateX(100%)' },},
        },
        slideIn: {
          '0%': { transform: 'translateX(100%)' },
          '100%': { transform: 'translateX(0)' },
        },
      },
    },
  },
  plugins: [],
}