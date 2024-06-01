<script>
    export let input_type = "text";
    export let name = "";
    export let label;
    export let value = "";

    function typeAction(node) {
        node.type = input_type;
    }
</script>

<div class="input-and-label">
    <label for="my_input">{label}</label>
    <input id="my_input" {name} use:typeAction placeholder={label} bind:value />
</div>

<style>
    .input-and-label {
        display: flex;
        flex-direction: column;
        align-items: stretch;
        margin: 1em 0;
    }

    .input-and-label label {
        margin-bottom: 0.5em;
        text-align: left;
    }

    .input-and-label input {
        padding: 0.5em;
        border: 1px solid #ccc;
        border-radius: 4px;
    }
</style>
