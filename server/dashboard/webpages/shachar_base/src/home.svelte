<script>
    import Collection from "./collection.svelte";
    import CollectionCard from "./lib/CollectionCard.svelte";

    let collections = [
        {
            ID: 1,
            Name: "User",
            QueryRulesDirectoryPath: "access_rules/rules",
            Fields: [
                {
                    ID: 1,
                    CollectionID: 0,
                    FieldName: "Username",
                    FieldType: "text",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
                {
                    ID: 2,
                    CollectionID: 0,
                    FieldName: "Password",
                    FieldType: "BLOB",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
            ],
        },
        {
            ID: 2,
            Name: "Animal",
            QueryRulesDirectoryPath: "access_rules/rules",
            Fields: [
                {
                    ID: 1,
                    CollectionID: 0,
                    FieldName: "Name",
                    FieldType: "text",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
                {
                    ID: 2,
                    CollectionID: 0,
                    FieldName: "Species",
                    FieldType: "text",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
            ],
        },
        {
            ID: 3,
            Name: "Book",
            QueryRulesDirectoryPath: "access_rules/rules",
            Fields: [
                {
                    ID: 1,
                    CollectionID: 0,
                    FieldName: "Title",
                    FieldType: "text",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
                {
                    ID: 2,
                    CollectionID: 0,
                    FieldName: "Author",
                    FieldType: "text",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
            ],
        },
        {
            ID: 4,
            Name: "Course",
            QueryRulesDirectoryPath: "access_rules/rules",
            Fields: [
                {
                    ID: 1,
                    CollectionID: 0,
                    FieldName: "Name",
                    FieldType: "text",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
                {
                    ID: 2,
                    CollectionID: 0,
                    FieldName: "Duration",
                    FieldType: "INTEGER",
                    FieldOptions: { String: "NOT NULL", Valid: true },
                },
            ],
        },
    ];

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

<main>
    <div class="collections-screen">
        <h1 class="home-heading">Shachar Base</h1>
        <hr />
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
