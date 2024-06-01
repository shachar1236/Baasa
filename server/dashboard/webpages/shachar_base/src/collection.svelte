<script>
    import CollectionsEditScreen from "./lib/collection/CollectionsEditScreen.svelte";
    import CollectionsRulesScreen from "./lib/collection/CollectionsRulesScreen.svelte";
    import { writable } from "svelte/store";

    let collection = {
        id: 1,
        table: [
            ["post", "string", ""],
            ["username", "key", ""],
            ["content", "string", ""],
            ["date", "date", ""],
        ],
    };

    let currentScreen = CollectionsEditScreen;
</script>

<main>
    <h1 class="collection-heading">Collection</h1>

    <nav>
        <ul>
            <li>
                <button
                    on:click={() => {
                        currentScreen = CollectionsEditScreen;
                    }}>Table</button
                >
            </li>
            <li>
                <button
                    on:click={() => (currentScreen = CollectionsRulesScreen)}
                    >Rules</button
                >
            </li>
        </ul>
    </nav>
    <hr />

    <svelte:component this={currentScreen} {collection} />
</main>

<style>
    nav {
        overflow: hidden;
    }

    nav ul {
        list-style-type: none;
        padding: 0;
        overflow: hidden;
    }

    nav li {
        float: left;
    }

    nav li button {
        display: block;
        color: white;
        text-align: center;
        margin-left: 1rem;
        padding: 14px 16px;
        text-decoration: none;
    }

    nav li button:focus {
        outline: none;
    }
</style>
