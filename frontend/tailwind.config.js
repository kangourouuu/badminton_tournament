/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Outfit", "sans-serif"],
      },
      colors: {
        background: "#FDFBFF",
        surface: "#FFFFFF",
        primary: "#7C3AED",
        border: "#E9D5FF",
      },
      borderRadius: {
        DEFAULT: "2px",
        sm: "2px",
        md: "4px",
        lg: "8px",
      },
      boxShadow: {
        none: "none",
      },
    },
  },
  plugins: [],
};
