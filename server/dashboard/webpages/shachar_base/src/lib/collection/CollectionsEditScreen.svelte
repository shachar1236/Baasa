<script>
    export let collection;
</script>

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
            {#each collection.table as row, rowIndex}
                <tr>
                    {#each row as cell, colIndex}
                        <td contenteditable={true}>{cell}</td>
                    {/each}
                </tr>
            {/each}
        </table>
    </div>
</div>

<style>
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
</style>
