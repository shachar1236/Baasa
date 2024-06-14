<script>
    export let collection;

    function saveChanges() {
        console.log(collection);
        // post request
        fetch("/SaveCollectionChanges", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(collection),
        }).then((response) => {
            // refresh
            window.location.href = window.location.href;
        });
    }

    function addField() {
        // add field
        if (collection.Fields == null) collection.Fields = [];

        collection.Fields.push({
            ID: 0,
            CollectionID: collection.ID,
            FieldName: "",
            FieldType: "",
            FieldOptions: { String: "", Valid: true },
        });
        collection = collection;
    }
</script>

<h1
    id="collection-name"
    contenteditable="true"
    bind:textContent={collection.Name}
>
    {collection.Name}
</h1>
<div class="collections-edit-screen">
    <div class="edit-table">
        <h2>Fields</h2>
        <hr />
        <table>
            <tr>
                <td class="fields">field name</td>
                <td class="fields">field type</td>
                <td class="fields">options</td>
            </tr>
            {#if collection.Fields != null}
                {#each collection.Fields as _, i}
                    <tr>
                        <td
                            contenteditable="true"
                            bind:textContent={collection.Fields[i].FieldName}
                            >{collection.Fields[i].FieldName}</td
                        >
                        <td
                            contenteditable="true"
                            bind:textContent={collection.Fields[i].FieldType}
                            >{collection.Fields[i].FieldType}</td
                        >
                        <td
                            contenteditable="true"
                            bind:textContent={collection.Fields[i].FieldOptions
                                .String}
                            >{collection.Fields[i].FieldOptions.String}</td
                        >

                        <!-- remove button -->
                        <button
                            id="remove-button"
                            on:click={() => {
                                collection.Fields.splice(i, 1);
                                collection = collection;
                            }}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
                                ></path>
                            </svg>
                        </button>
                    </tr>
                {/each}
            {/if}

            <button id="plus-button" on:click={addField}>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                >
                    <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"></path>
                </svg>
            </button>
        </table>
        <!-- save changes button -->
        <button id="save-button" on:click={saveChanges}>save changes</button>
    </div>
</div>

<style>
    #collection-name {
        font-size: 2rem;
    }

    .collections-edit-screen {
        width: 94vw;
        /* height: 100vh; */
        margin-top: 2rem;
        background-color: #f0f0f0;
        border-radius: 20px;
        padding: 4rem 1.5rem;
        transition: all 0.2s ease-in-out;
    }

    h2 {
        color: #111;
    }

    hr {
        width: 100%;
        margin: 0.5rem 0;
        border-color: #111;
    }

    .edit-table {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        margin-top: 3rem;
        width: 100%;
        height: 60%;
        position: relative;
    }

    table {
        background-color: #e9e9e9;
        width: 100%;
        border-collapse: collapse;
        margin: 3rem 0;
    }

    .fields {
        font-weight: bold;
        color: #646cff;
        font-size: 1.3rem;
    }

    td {
        border: 1px solid #666666;
        text-align: center;
        color: black;
        padding: 10px;
    }

    #remove-button {
        background-color: transparent;
        border: none;
        border-radius: 0;
        position: absolute;
        right: 0;
        cursor: pointer;
    }

    #remove-button:focus {
        outline: none;
    }

    #plus-button {
        background-color: transparent;
        border: none;
        border-radius: 0;
        position: absolute;
        bottom: 0;
        right: 0;
        cursor: pointer;
    }

    #plus-button:focus {
        outline: none;
    }

    #save-button {
        background-color: #646cff;
        color: white;
        padding: 1em;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }
</style>
