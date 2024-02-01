import type { Config } from 'tailwindcss'
import { nextui } from '@nextui-org/theme'
const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './node_modules/@nextui-org/theme/dist/components/input.js',
    './node_modules/@nextui-org/theme/dist/components/button.js',
  ],
  theme: {
    screens: {
      xss: '320px',
      xs: '480px',
      sm: '640px',
      md: '768px',
      lg: '1024px',
    },
    extend: {
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic': 'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
      },
    },
  },
  plugins: [nextui()],
}
export default config
