export default function Loading() {
  return (
    <div className="relative h-1 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
      <div className="absolute inset-0 bg-gradient-to-r from-gray-300 to-gray-500 dark:from-gray-200 dark:to-gray-400 animate-indeterminate" />
    </div>
  );
}
