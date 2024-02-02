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
    extend: {},
  },
  plugins: [nextui()],
}
export default config
