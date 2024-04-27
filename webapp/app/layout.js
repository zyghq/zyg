import "./globals.css";
import { fontSans, fontMono } from "@/lib/fonts";
import {
  ThemeProvider,
  ReactQueryClientProvider,
} from "@/components/providers";
import { Toaster } from "@/components/ui/toaster";
import { cn } from "@/lib/utils";

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
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            enableSystem
            disableTransitionOnChange
          >
            {children}
            <Toaster />
          </ThemeProvider>
        </body>
      </html>
    </ReactQueryClientProvider>
  );
}
