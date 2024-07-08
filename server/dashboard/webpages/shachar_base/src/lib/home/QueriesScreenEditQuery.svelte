<script>
    import CodeMirrorEditor from "../CodeMirrorEditor.svelte";
    import TextEditor from "../TextEditor.svelte";
    import RulesEditor from "../TextEditor.svelte";

    export let query;
    export let collections;

    let new_query = query;

    function changeQuery() {
        console.log(new_query);
        query.Query = new_query;
        console.log(query);
        fetch("SetQuery", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(query),
        });
    }

    function onChange(value) {
        new_query = value;
        changeQuery();
    }
</script>

<br />
<CodeMirrorEditor {collections} lang="sql" code={query.Query} {onChange}
></CodeMirrorEditor>

<style></style>
