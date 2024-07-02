const { fontFamily } = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["../views/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['MerriweatherSans', ...fontFamily.sans]
      }
    },
  },
  plugins: [
    require("daisyui")
  ],
  daisyui: {
    themes: ["dark"]
  }
}

