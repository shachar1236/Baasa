<script>
    import CollectionCard from "./lib/CollectionCard.svelte";

    let collections = [
        {
            id: 1,
            table: [
                ["username", "string", ""],
                ["password", "string", ""],
            ],
        },
        {
            id: 2,
            table: [
                ["post", "string", ""],
                ["username", "key", ""],
                ["content", "string", ""],
                ["date", "date", ""],
            ],
        },
        {
            id: 3,
            table: [
                ["username", "string", ""],
                ["password", "string", ""],
            ],
        },
        {
            id: 4,
            table: [
                ["username", "string", ""],
                ["password", "string", ""],
            ],
        },
    ];

    let currentCollections = null;

    function openCollection(collection) {
        window.location.href = `/collection.html?id=${collection.id}`;
    }
    // <hr style="margin: 0.5rem 0;" />;
</script>

<main>
    <div class="collections-screen">
        <h1 class="home-heading">Shachar Base</h1>
        {#if currentCollections}
            <hr style="margin: 0 0.5rem;" />
        {:else}
            <hr />
        {/if}
        <div class="cards-div">
            {#each collections as collection}
                <CollectionCard
                    onClick={() => openCollection(collection)}
                    cellData={collection.table}
                />
            {/each}
        </div>
    </div>
</main>

<style>
    main {
        margin: 0;
        display: flex;
    }

    .collections-screen {
        /* make container all the screen */
        width: 100vw;
        height: 100vh;
    }

    .home-heading {
        color: #646cff;
    }

    .cards-div {
        display: flex;
        flex-wrap: wrap;
        margin-top: 3rem;
    }
</style>
