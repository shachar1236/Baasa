import "./app_v2.css";
import Collection from "./collection.svelte";

const app = new Collection({
    target: document.getElementById("app"),
});

export default app;
