<script>
    export let value = "";

    $: newValue = value;

    export let onExit = (data) => {};

    function changeValue() {
        onExit(newValue);
    }
</script>

<textarea
    rows="10"
    cols="50"
    placeholder="Enter Rules Here"
    bind:value={newValue}
    on:blur={() => changeValue()}
></textarea>

<style>
    textarea {
        width: 100%;
        height: 100%;
        resize: none;
        border: 1px solid #ccc;
        border-radius: 4px;
        padding: 10px;
        margin: 1rem;
        background-color: #1a1a1a;
    }
</style>
