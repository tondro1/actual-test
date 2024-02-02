/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./public/**/*.html"],
    theme: {
        extend: {},
    },
    plugins: [require("daisyui")],
    daisyui: {
        themes: ["light", "dark", "forest", "dim"],
    },
};
