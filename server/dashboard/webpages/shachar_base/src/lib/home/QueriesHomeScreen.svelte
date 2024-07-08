<script>
    import { writable } from "svelte/store";
    import Modal from "../Modal.svelte";
    import RulesEditor from "../TextEditor.svelte";
    import QueriesScreenEditQuery from "./QueriesScreenEditQuery.svelte";
    import QueriesScreenEditRules from "./QueriesScreenEditRules.svelte";

    export let collections;

    let queiries = [];
    let selected_query = writable(null);
    let editScreen = "Query";

    let new_query = {
        Name: "",
    };
    let showModal = false;

    fetch("/GetQueries").then((response) => {
        response.json().then((data) => {
            queiries = data;
        });
    });

    function changeQuery(new_query) {
        console.log(new_query);
        selected_query.set(new_query);
    }

    function addQuery() {
        fetch("/AddQuery", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(new_query),
        }).then((response) => {
            response.json().then((data) => {
                console.log(data);
                queiries = [...queiries, data];
                new_query = {
                    Name: "",
                };
            });
        });
    }

    function deleteQuery(id) {
        fetch("/DeleteQuery?id=" + id, { method: "DELETE" }).then(
            (response) => {
                if (response.ok) {
                    queiries = queiries.filter((query) => query.ID != id);
                } else {
                    console.log(response);
                    alert("Failed to delete query: " + response.body);
                }
            }
        );
    }
</script>

<div class="rules-edit-screen">
    <div class="rules-table">
        <nav>
            <ul>
                {#if $selected_query != null}
                    <li>
                        <h3 style="margin: 0; color: green;">
                            {$selected_query.Name}
                        </h3>
                        <hr />
                    </li>
                {/if}
                {#each queiries as query}
                    <li class="query-li">
                        <button on:click={() => changeQuery(query)}
                            >{query.Name}</button
                        >
                        <button
                            class="delete-button"
                            on:click|stopPropagation={() => {
                                deleteQuery(query.ID);
                            }}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                x="0px"
                                y="0px"
                                width="1.5rem"
                                height="1.5rem"
                                viewBox="0 0 50 50"
                            >
                                <path
                                    d="M 42 5 L 32 5 L 32 3 C 32 1.347656 30.652344 0 29 0 L 21 0 C 19.347656 0 18 1.347656 18 3 L 18 5 L 8 5 C 7.449219 5 7 5.449219 7 6 C 7 6.550781 7.449219 7 8 7 L 9.085938 7 L 12.695313 47.515625 C 12.820313 48.90625 14.003906 50 15.390625 50 L 34.605469 50 C 35.992188 50 37.175781 48.90625 37.300781 47.515625 L 40.914063 7 L 42 7 C 42.554688 7 43 6.550781 43 6 C 43 5.449219 42.554688 5 42 5 Z M 20 44 C 20 44.554688 19.550781 45 19 45 C 18.449219 45 18 44.554688 18 44 L 18 11 C 18 10.449219 18.449219 10 19 10 C 19.550781 10 20 10.449219 20 11 Z M 20 3 C 20 2.449219 20.449219 2 21 2 L 29 2 C 29.550781 2 30 2.449219 30 3 L 30 5 L 20 5 Z M 26 44 C 26 44.554688 25.550781 45 25 45 C 24.449219 45 24 44.554688 24 44 L 24 11 C 24 10.449219 24.449219 10 25 10 C 25.550781 10 26 10.449219 26 11 Z M 32 44 C 32 44.554688 31.554688 45 31 45 C 30.445313 45 30 44.554688 30 44 L 30 11 C 30 10.449219 30.445313 10 31 10 C 31.554688 10 32 10.449219 32 11 Z"
                                ></path>
                            </svg>
                        </button>
                    </li>
                {/each}
                <hr />
                <li>
                    <button on:click={() => (showModal = true)}
                        >Add query</button
                    >
                </li>
            </ul>
        </nav>
        <hr style="margin: 0; height: 100%;" />

        <div class="editors">
            <h2 style="margin-left: 1rem;">{editScreen}</h2>
            <div class="editors-choise">
                <button on:click={() => (editScreen = "Query")}>Query</button>
                <button on:click={() => (editScreen = "Rules")}>Rules</button>
            </div>
            {#if $selected_query != null}
                {#if editScreen == "Query"}
                    <QueriesScreenEditQuery
                        query={$selected_query}
                        {collections}
                    ></QueriesScreenEditQuery>
                {:else if editScreen == "Rules"}
                    <QueriesScreenEditRules
                        query={$selected_query}
                        {collections}
                    ></QueriesScreenEditRules>
                {/if}
            {/if}
        </div>
    </div>
</div>

<!-- create query modal -->
<Modal bind:showModal>
    <h2 slot="header">Create query</h2>
    <table>
        <tr>
            <td>Name</td>
            <td contenteditable="true" bind:textContent={new_query.Name}
                >{new_query.Name}</td
            >
        </tr>
    </table>
    <!-- save and exit button -->
    <div>
        <button on:click={addQuery} id="save-button">save query</button>
        <button
            on:click={() => {
                showModal = false;
            }}
            id="exit-button">exit</button
        >
    </div>
</Modal>

<style>
    .rules-edit-screen {
        height: 100vh;
    }

    .rules-table {
        display: flex;
        height: 100%;
    }

    .editors {
        height: 100%;
        width: 88%;
        margin-top: 1%;
        margin-left: 2%;
    }

    .editors-choise {
        height: 5%;
        width: 100%;
        display: flex;
    }

    .editors-choise button {
        margin-left: 1%;
        width: 20%;
    }

    nav {
        margin-top: 1rem;
        width: 8rem;
        height: 100%;
        border-radius: 4px;
        padding: 10px;
    }

    ul {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    li {
        /* margin-bottom: 1rem; */
    }

    li button {
        width: 100%;
    }

    /* remove border outline after clinking on button */
    li button:focus {
        outline: none;
    }

    .query-li .delete-button {
        background-color: transparent;
        border: none;
        cursor: pointer;
        margin-top: 0.5rem;
        margin-bottom: 0;
        padding: 0;
    }

    #save-button {
        background-color: transparent;
        border: none;
        border-radius: 0;
        cursor: pointer;
    }

    #exit-button {
        background-color: transparent;
        border: none;
        border-radius: 0;
        cursor: pointer;
    }

    #exit-button:hover {
        color: red;
    }

    #save-button:hover {
        color: green;
    }
</style>
