import type { Metadata } from 'next'
import { Gabarito } from 'next/font/google'
import './globals.css'
import { ThemeProvider, AuthProvider } from '@/components/provider'
import { Toaster } from '@/components/ui/toaster'

const gabarito = Gabarito({
  weight: '400',
  subsets: ['latin'],
  display: 'swap',
  adjustFontFallback: false,
})

export const metadata: Metadata = {
  title: 'Stargazer',
  description: 'A GitHub repository stargazer search engine',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang='en' suppressHydrationWarning>
      <body className={gabarito.className}>
        <ThemeProvider attribute='class' defaultTheme='system' enableSystem disableTransitionOnChange>
          <AuthProvider>{children}</AuthProvider>
        </ThemeProvider>
        <Toaster />
      </body>
    </html>
  )
}
