<script>
    import Modal from "./../Modal.svelte";

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
    let new_row = {};
    collection.Fields.forEach((field) => {
        new_row[field.FieldName] = "";
    });

    let showModal = false;

    fetch("GetData", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            collection_name: collection.Name,
            from: 0,
            to: 10,
        }),
    }).then((response) => {
        response.json().then((new_data) => {
            console.log(new_data);
            data = new_data;
        });
    });

    console.log(data[0]["ID"]);

    let column;

    function changeColumnStart(dataIndex, dataName) {
        console.log(dataIndex, dataName);
        column = data[dataIndex][dataName];
        console.log(column);
    }

    function changeColumnEnd(dataIndex, dataName) {
        fetch("SetById", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                id: data[dataIndex]["id"],
                collection_name: collection.Name,
                column_name: dataName,
                to: data[dataIndex][dataName],
            }),
        }).then((response) => {
            // read message and see if its error
            console.log(response);
        });
    }

    function addRow() {
        fetch("AddWithArgs", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                collection_name: collection.Name,
                args: new_row,
            }),
        }).then((response) => {
            console.log(response);
            // convert from json
            response.json().then((new_data) => {
                console.log(new_data.id);
                let id = parseInt(new_data.id);
                // if its a number and not string
                new_row["id"] = id;
                // @ts-ignore
                data.push({ ...new_row });
                data = data;

                collection.Fields.forEach((field) => {
                    new_row[field.FieldName] = "";
                });
            });
        });
    }

    function deleteRow(dataIndex) {
        console.log(dataIndex);
        fetch("DeleteById", {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                id: data[dataIndex]["id"],
                collection_name: collection.Name,
            }),
        }).then((response) => {
            // read message and see if its error
            console.log(response);
            // remove from data
            data = data.filter((_, i) => i != dataIndex);
        });
    }
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
                    <td>{data[i]["id"]}</td>
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
                    <button id="remove-button" on:click={() => deleteRow(i)}>
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

        <button id="plus-button" on:click={() => (showModal = true)}>
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

<Modal bind:showModal>
    <h2 slot="header">Create row</h2>
    <table>
        <tr>
            {#if collection.Fields != null}
                {#each collection.Fields as _, i}
                    <td class="fields">{collection.Fields[i].FieldName}</td>
                {/each}
            {/if}
        </tr>
        <tr>
            {#if collection.Fields != null}
                {#each collection.Fields as _, i}
                    <td
                        contenteditable="true"
                        bind:textContent={new_row[
                            collection.Fields[i].FieldName
                        ]}>{new_row[collection.Fields[i].FieldName]}</td
                    >
                {/each}
            {/if}
        </tr>
    </table>
    <!-- save and exit button -->
    <button
        id="save-button"
        on:click={() => {
            let can_save = true;
            collection.Fields.forEach((field) => {
                let contains_not_null =
                    field.FieldOptions.String.toLowerCase().includes(
                        "not null"
                    );
                let contains_default =
                    field.FieldOptions.String.toLowerCase().includes("default");
                if (
                    contains_not_null &&
                    !contains_default &&
                    !new_row[field.FieldName]
                ) {
                    alert(field.FieldName + " must have a value");
                    can_save = false;
                }
            });

            if (can_save) {
                addRow();
                showModal = false;
            }
        }}
    >
        save
    </button>
    <button
        id="exit-button"
        on:click={() => {
            showModal = false;
        }}
    >
        exit
    </button>
</Modal>

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
