/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,ts,tsx,vue}'],
  safelist: [
    'text-wald-400', 'text-wald-300', 'text-wald-400/70',
    'text-ritter-400', 'text-ritter-300', 'text-ritter-400/70', 'text-ritter-500',
    'text-wald-500', 'from-wald-300', 'from-ritter-300',
    'to-wald-300', 'to-ritter-300',
  ],
  theme: {
    extend: {
      colors: {
        wald: {
          50: '#edfff5',
          100: '#d5ffea',
          200: '#aeffd6',
          300: '#00FF88',
          400: '#1BAE70',
          500: '#159458',
          600: '#0f7a48',
          700: '#0A5E35',
          800: '#0A3D1F',
          900: '#0A2A16',
          950: '#0A1F14',
        },
        ritter: {
          50: '#fffbeb',
          100: '#fff3c6',
          200: '#ffe588',
          300: '#FFD700',
          400: '#D4A017',
          500: '#B8860B',
          600: '#9A7209',
          700: '#7A5A07',
          800: '#4A3000',
          900: '#3A2500',
          950: '#2A1A00',
        },
      },
      fontFamily: {
        mono: ['"JetBrains Mono"', 'monospace'],
        display: ['"Orbitron"', 'sans-serif'],
        body: ['"Exo 2"', 'sans-serif'],
      },
      animation: {
        'glow-pulse': 'glow-pulse 3s ease-in-out infinite',
        'scan-line': 'scan-line 4s linear infinite',
        'flicker': 'flicker 0.15s infinite',
        'float': 'float 6s ease-in-out infinite',
        'wald-spin': 'wald-spin 20s linear infinite',
        'hologram-interference': 'hologram-interference 8s linear infinite',
      },
      keyframes: {
        'glow-pulse': {
          '0%, 100%': { opacity: '0.4', filter: 'blur(20px)' },
          '50%': { opacity: '0.8', filter: 'blur(30px)' },
        },
        'scan-line': {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100vh)' },
        },
        'flicker': {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.8' },
        },
        'float': {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' },
        },
        'wald-spin': {
          '0%': { transform: 'rotate(0deg)' },
          '100%': { transform: 'rotate(360deg)' },
        },
        'hologram-interference': {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' },
        },
      },
    },
  },
  plugins: [require('@tailwindcss/typography')],
};
