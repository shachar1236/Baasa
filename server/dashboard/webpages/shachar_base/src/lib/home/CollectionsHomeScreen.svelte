<script>
    import Collection from "../../collection.svelte";
    import CollectionCard from "../CollectionCard.svelte";

    export let collections;

    fetch("/GetCollections").then((response) => {
        response.json().then((data) => {
            collections = data;
        });
    });

    function openCollection(collection) {
        window.location.href = `/collection?id=${collection.ID}`;
    }

    function addCollection() {
        console.log("add collection");
        fetch("/AddCollection", { method: "POST" }).then((response) => {
            response.json().then((data) => {
                window.location.href = window.location.href;
            });
        });
    }

    function deleteCollection(name) {
        console.log("delete collection");
        fetch("/DeleteCollection?name=" + name, { method: "DELETE" }).then(
            (response) => {
                window.location.href = window.location.href;
            }
        );
    }
    // <hr style="margin: 0.5rem 0;" />;
</script>

<div class="cards-div">
    {#each collections as collection}
        <CollectionCard
            onClick={() => openCollection(collection)}
            deleteFunc={deleteCollection}
            {collection}
        />
    {/each}
    <CollectionCard isPlus={true} onClick={() => addCollection()} />
</div>

<style>
    .cards-div {
        display: flex;
        flex-wrap: wrap;
        margin-top: 3rem;
    }
</style>
