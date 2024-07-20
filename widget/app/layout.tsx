import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";

import { CustomerProvider } from "@/components/providers";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Zyg Widget",
  description: "Zyg Widget Powered By Zyg AI",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <CustomerProvider>{children}</CustomerProvider>
      </body>
    </html>
  );
}
