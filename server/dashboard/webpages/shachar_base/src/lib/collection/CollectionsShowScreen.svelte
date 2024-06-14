<script>
    export let collection;
    // collection = {
    //     ID: 1,
    //     Name: "User",
    //     QueryRulesDirectoryPath: "access_rules/rules",
    //     Fields: [
    //         {
    //             ID: 1,
    //             CollectionID: 0,
    //             FieldName: "Username",
    //             FieldType: "text",
    //             FieldOptions: { String: "NOT NULL", Valid: true },
    //         },
    //         {
    //             ID: 2,
    //             CollectionID: 0,
    //             FieldName: "Password",
    //             FieldType: "BLOB",
    //             FieldOptions: { String: "NOT NULL", Valid: true },
    //         },
    //     ],
    // };
    let data = [
        {
            ID: 1,
            Username: "admin",
            Password: "admin",
        },
        {
            ID: 2,
            Username: "user",
            Password: "user",
        },
        {
            ID: 3,
            Username: "guest",
            Password: "guest",
        },
    ];

    console.log(data[0]["ID"]);

    let column;

    function changeColumnStart(dataIndex, dataName) {
        console.log(dataIndex, dataName);
        column = data[dataIndex][dataName];
        console.log(column);
    }

    function changeColumnEnd(dataIndex, dataName) {
        console.log(dataIndex, dataName);
        column = data[dataIndex][dataName];
        console.log(column);
    }

    function addRow() {}

    function deleteRow() {}
</script>

<div class="collection-show-screen">
    <h2>Data</h2>
    <hr />
    <table>
        <tr>
            <td class="fields">ID</td>
            {#if collection.Fields != null}
                {#each collection.Fields as _, i}
                    <td class="fields">{collection.Fields[i].FieldName}</td>
                {/each}
            {/if}
        </tr>
        {#if data != null}
            {#each data as _, i}
                <tr>
                    <td>{data[i]["ID"]}</td>
                    {#if collection.Fields != null}
                        {#each collection.Fields as _, j}
                            <td
                                contenteditable="true"
                                on:focus={() =>
                                    changeColumnStart(
                                        i,
                                        collection.Fields[j].FieldName
                                    )}
                                on:blur={() =>
                                    changeColumnEnd(
                                        i,
                                        collection.Fields[j].FieldName
                                    )}
                                bind:textContent={data[i][
                                    collection.Fields[j].FieldName
                                ]}>{data[i][collection.Fields[j].FieldName]}</td
                            >
                        {/each}
                    {/if}

                    <!-- remove button -->
                    <button id="remove-button" on:click={deleteRow}>
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="white"
                        >
                            <path
                                d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
                            ></path>
                        </svg>
                    </button>
                </tr>
            {/each}
        {/if}

        <button id="plus-button" on:click={addRow}>
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="white"
            >
                <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"></path>
            </svg>
        </button>
    </table>
</div>

<style>
    table {
        width: 100%;
        /* border-collapse: collapse; */
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
        /* color: black; */
        padding: 10px;
    }

    #remove-button {
        background-color: transparent;
        border: none;
        border-radius: 0;
        position: absolute;
        right: 2rem;
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
        right: 2rem;
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
