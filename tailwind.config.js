/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/handler/views/*.templ",
    "node_modules/preline/dist/*.js",
  ],
  // enable dark mode via class strategy
  // darkMode: "class",
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("@tailwindcss/forms")],
};
