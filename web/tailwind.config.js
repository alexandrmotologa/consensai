/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#0B0F17',
        surface: '#111827',
        'surface-elevated': '#1F2937',
        border: '#374151',
        primary: {
          50: '#ecfeff',
          500: '#06b6d4',
          600: '#0891b2',
        },
        advocate: '#10B981', // emerald
        skeptic: '#F43F5E',  // rose
        auditor: '#F59E0B',  // amber
      },
      fontFamily: {
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace'],
        sans: ['Inter', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
