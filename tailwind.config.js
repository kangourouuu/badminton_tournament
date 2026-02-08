/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Outfit', 'sans-serif'], // Ép dùng Outfit
      },
      borderRadius: {
        DEFAULT: '2px', // Chuẩn Flat/Sharp
        'sm': '1px',
        'md': '3px',
        'lg': '4px',
        'full': '9999px',
      },
      boxShadow: {
        // Xóa sạch shadow mặc định để ép dùng Border
        'sm': 'none',
        DEFAULT: 'none',
        'md': 'none',
        'lg': 'none',
        'xl': 'none',
        '2xl': 'none',
        'inner': 'none',
      },
      colors: {
        // Pale Purple Palette
        brand: {
          50: '#FDFBFF',  // Background chính
          100: '#F3E8FF', // Background phụ
          200: '#E9D5FF', // Border chính
          500: '#A855F7', // Primary Action
          600: '#9333EA', // Hover
          900: '#581C87', // Text đậm
        }
      }
    },
  },
  plugins: [],
}