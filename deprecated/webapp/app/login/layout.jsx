export const metadata = {
  title: "Login Or Sign Up | Zyg AI",
};

export default function Layout({ children }) {
  return (
    <div className="flex min-h-screen flex-col justify-center p-4">
      {children}
    </div>
  );
}
