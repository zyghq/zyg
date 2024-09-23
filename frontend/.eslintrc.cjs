module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react-hooks/recommended",
    "plugin:perfectionist/recommended-natural-legacy",
  ],
  ignorePatterns: ["dist", ".eslintrc.cjs", "**/components/ui/**"],
  parser: "@typescript-eslint/parser",
  plugins: ["react-refresh", "perfectionist"],
  rules: {
    // "react-refresh/only-export-components": [
    //   "warn",
    //   { allowConstantExport: true },
    // ],
    "@typescript-eslint/no-explicit-any": "off",
  },
  // ignores: ["**/components/ui/**"],
};
