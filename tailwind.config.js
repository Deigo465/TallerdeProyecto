/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./pkg/web/views/**/*.tmpl",
    "./pkg/web/views/**/*.vue",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#0C065C',
        secondary: '#7FA3FF',
      },
    },
  },
  plugins: [],
}

