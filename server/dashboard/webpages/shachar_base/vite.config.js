import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [svelte()],
    build: {
        rollupOptions: {
            input: {
                register: "register.html", // Entry point for the home page
                login: "login.html", // Entry point for the about page
                collection: "collection.html",
                home: "home.html",
            },
        },
    },
});
