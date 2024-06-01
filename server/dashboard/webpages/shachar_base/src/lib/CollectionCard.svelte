<!-- CollectionCard.svelte -->
<script>
    export let cellData = [[]];
    let rows;
    let columns;
    $: rows = cellData.length;
    $: columns = cellData[0].length;
    export let onClick = () => {};
</script>

<button class="card" on:click={onClick}>
    <table>
        {#each Array(rows) as _, rowIndex}
            <tr>
                {#each Array(columns) as _, colIndex}
                    <td>{cellData[rowIndex][colIndex]}</td>
                {/each}
            </tr>
        {/each}
    </table>
</button>

<style>
    .card {
        margin: 0 0 1rem 1rem;
        border-radius: 20px;
        padding: 20px;
        background-color: #f0f0f0;
        transition: all 0.2s ease-in-out;
    }

    .card:hover {
        transform: scale(1.05);
        background-color: #e0e0e0;
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    td {
        border: 1px solid #ccc;
        color: black;
        padding: 10px;
    }
</style>
