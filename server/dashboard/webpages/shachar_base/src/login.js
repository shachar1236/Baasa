import "./app.css";
import Login from "./login.svelte";

const app = new Login({
    target: document.getElementById("app"),
});

export default app;
