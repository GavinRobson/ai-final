/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.go",               // Go files where you build HTML
    "./templates/**/*.html",   // if you have templates
    "./static/**/*.html",      // any static HTML
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

