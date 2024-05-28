import { fontMono, fontSans } from "@/lib/fonts";
import { cn } from "@/lib/utils";

import { Toaster } from "@/components/ui/toaster";

import {
  JotaiProvider,
  ReactQueryClientProvider,
  ThemeProvider,
} from "@/components/providers";

import "./globals.css";

export const metadata = {
  title: "Zyg - Developer first customer support platform.",
  description:
    "Transform your customer support experience, with a developer-first and GPT powered platform.",
};

export default function RootLayout({ children }) {
  return (
    <ReactQueryClientProvider>
      <html
        lang="en"
        className={`${fontSans.variable} ${fontMono.variable}`}
        suppressHydrationWarning
      >
        <head />
        <body className={cn("bg-background antialiased")}>
          <JotaiProvider>
            <ThemeProvider
              attribute="class"
              defaultTheme="system"
              enableSystem
              disableTransitionOnChange
            >
              {children}
              <Toaster />
            </ThemeProvider>
          </JotaiProvider>
        </body>
      </html>
    </ReactQueryClientProvider>
  );
}
