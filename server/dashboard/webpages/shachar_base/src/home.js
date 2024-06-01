import "./app_v2.css";
import Home from "./home.svelte";

const app = new Home({
    target: document.getElementById("app"),
});

export default app;
