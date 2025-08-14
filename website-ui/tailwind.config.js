/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],

  theme: {
    fontFamily: {
      heading: ['"Lato"', "system-ui", "sans-serif"],
      body: ['"Montserrat"', "system-ui", "sans-serif"],
    },
    extend: {
      typography: (theme) => ({
        DEFAULT: {
          css: {
            fontSize: theme('fontSize.base'),
            fontFamily: theme('fontFamily.body'),
            h1: {
              fontFamily: theme('fontFamily.heading'),
              fontWeight: '400',
              fontSize: theme('fontSize.2xl'),
            },
            h2: {
              fontFamily: theme('fontFamily.heading'),
              fontWeight: '200',
              fontSize: theme('fontSize.xl'),
            },
            h3: {
              fontFamily: theme('fontFamily.heading'),
              fontWeight: '300',
              fontSize: theme('fontSize.lg'),
            },
            h4: {
              fontFamily: theme('fontFamily.heading'),
              fontWeight: '400',
              fontSize: theme('fontSize.base'),
            },
            h5: {
              fontFamily: theme('fontFamily.heading'),
              fontWeight: '500',
              fontSize: theme('fontSize.sm'),
            },
            h6: {
              fontFamily: theme('fontFamily.heading'),
              fontWeight: '600',
              fontSize: theme('fontSize.xs'),
            },
          },
        },
        sm: {
          css: {
            fontSize: theme('fontSize.sm'),
            h1: {
              fontSize: theme('fontSize.lg'),
            },
            h2: {
              fontSize: theme('fontSize.base'),
            },
            h3: {
              fontSize: theme('fontSize.sm'),
            },
            h4: {
              fontSize: theme('fontSize.xs'),
            },
            h5: {
              fontSize: theme('fontSize.xs'),
            },
            h6: {
              fontSize: theme('fontSize.xs'),
            },
          },
        },
        lg: {
          css: {
            fontSize: theme('fontSize.lg'),
            h1: {
              fontSize: theme('fontSize.3xl'),
            },
            h2: {
              fontSize: theme('fontSize.2xl'),
            },
            h3: {
              fontSize: theme('fontSize.xl'),
            },
            h4: {
              fontSize: theme('fontSize.lg'),
            },
            h5: {
              fontSize: theme('fontSize.base'),
            },
            h6: {
              fontSize: theme('fontSize.sm'),
            },
          },
        },
      }),
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

