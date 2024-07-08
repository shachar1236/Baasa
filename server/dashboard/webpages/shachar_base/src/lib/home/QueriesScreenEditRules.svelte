<script>
    import CodeMirrorEditor from "../CodeMirrorEditor.svelte";
    import TextEditor from "../TextEditor.svelte";

    export let query;

    let rules = null;
    export let collections;

    let CustomKeywords = [
        [
            "Query",
            "Query(query_cmd string, args []string) - query from the database",
        ],
        [
            "Count",
            "Count(collection_name string, filters string, args []string) - count rows in the database",
        ],
        [
            "Get",
            "Get(collection_name string, filters string, args []string) - get one row from the database",
        ],

        [
            "Request",
            `
Request struct {
    Method string,
    Headers http.Header,
    Auth types.User,
}`,
        ],
        ["Request.Method", "string"],
        ["Request.Headers", "map[string][]string"],

        [
            "Request.Auth",
            `
User struct {
	ID           int64,
	Username     string,
	PasswordHash PasswordHash,
	Session      string,
}
`,
        ],
        ["Request.Auth.ID", "int64"],
        ["Request.Auth.Username", "string"],
        ["Request.Auth.PasswordHash", "PasswordHash"],
        ["Request.Auth.Session", "string"],

        [
            "Filters",
            "string - the filters you want to apply to the requested query",
        ],
        [
            "Accept",
            "bool - if you want to accept or reject the requested query",
        ],

        // user collection
        ["users", "users"],
        [
            "users.get",
            "users.get(filters string, args []string) - get one row from the database",
        ],
        [
            "users.count",
            "users.count(filters string, args []string) - count rows in the database",
        ],
    ];

    for (let collection of collections) {
        CustomKeywords.push([collection.Name, collection.Name]);
        CustomKeywords.push([
            collection.Name + ".get",
            collection.Name +
                ".get(filters string, args []string) - get one row from the database",
        ]);
        CustomKeywords.push([
            collection.Name + ".count",
            collection.Name +
                ".count(filters string, args []string) - count rows in the database",
        ]);
    }

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
        changeRules();
    }
</script>

<!-- <TextEditor value={rules} onExit={changeRules}></TextEditor> -->
<br />
{#if rules != null}
    <CodeMirrorEditor lang="lua" code={rules} {onChange} {CustomKeywords}
    ></CodeMirrorEditor>
{/if}

<style></style>
