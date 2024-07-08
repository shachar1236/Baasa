<script>
    import CodeMirrorEditor from "../CodeMirrorEditor.svelte";
    import TextEditor from "../TextEditor.svelte";

    export let query;

    let rules = "";

    let last_updated = Date.now();
    let updated = false;

    let CustomKeywords = {
        Query: {},
        Count: {},

        Request: {
            Method: {},
            Headers: {
                Add: {},
                Set: {},
                Get: {},
                Values: {},
                Del: {},
                Write: {},
                Clone: {},
            },

            Auth: {
                ID: {},
                Username: {},
                PasswordHash: {},
                Session: {},
            },
        },

        Filters: {},
        Accept: {},
    };

    fetch("/GetQueryRules?query_id=" + query.ID).then((response) => {
        response.json().then((data) => {
            rules = data;
            console.log("rules: ", data);
        });
    });

    function changeRules() {
        fetch("SetQueryRules", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                ID: query.ID,
                new_rules: rules,
            }),
        });
    }

    function onChange(value) {
        rules = value;
        last_updated = Date.now();
        updated = false;
    }

    setInterval(() => {
        if (Date.now() - last_updated > 1000 && !updated) {
            updated = true;
            changeRules();
        }
    }, 1000);
</script>

<!-- <TextEditor value={rules} onExit={changeRules}></TextEditor> -->
<br />
<CodeMirrorEditor lang="lua" code={rules} {onChange} {CustomKeywords}
></CodeMirrorEditor>

<style></style>
