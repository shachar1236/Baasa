<script>
    import { onMount } from "svelte";
    import { EditorState } from "@codemirror/state";
    import { EditorView, basicSetup } from "codemirror";
    import { autocompletion, closeBrackets } from "@codemirror/autocomplete";
    import { keymap } from "@codemirror/view";
    import { indentWithTab } from "@codemirror/commands";
    import { StreamLanguage } from "@codemirror/language";
    import { lua } from "@codemirror/legacy-modes/mode/lua";
    import { standardSQL } from "@codemirror/legacy-modes/mode/sql";

    import { generateLuaAutocomplete } from "./luaAutocomplete.js";
    import { createSqlAutocomplete } from "./sqlAutocomplete.js";
    import { dracula } from "thememirror";

    export let code = "";
    export let onChange = (value) => {};
    export let collections = [];
    export let lang = "lua";
    export let CustomKeywords = {};

    let my_extentions = [];
    if (lang == "lua") {
        my_extentions.push(StreamLanguage.define(lua));
        let luaAutocomplete = generateLuaAutocomplete(CustomKeywords);
        my_extentions.push(
            autocompletion({
                override: [luaAutocomplete],
            })
        );
    } else if (lang == "sql") {
        my_extentions.push(StreamLanguage.define(standardSQL));
        let sqlAutocomplete = createSqlAutocomplete(collections);
        my_extentions.push(autocompletion({ override: [sqlAutocomplete] }));
    }

    let editorContainer;

    onMount(() => {
        // if (CustomKeywords.length > 0) {
        //     console.log(CustomKeywords);
        //     my_extentions.push(
        //         autocompletion({
        //             override: [createCustomWordsAutocomplete()],
        //         })
        //     );
        // }

        const startState = EditorState.create({
            doc: code,
            extensions: [
                basicSetup,
                keymap.of([indentWithTab]),
                EditorView.updateListener.of((update) => {
                    if (update.changes) {
                        onChange(update.state.doc.toString());
                    }
                }),
                EditorView.lineWrapping,
                closeBrackets(),
                dracula,
                ...my_extentions,
            ],
        });

        const view = new EditorView({
            state: startState,
            parent: editorContainer,
        });

        // Ensure the editor fits initial content
        view.requestMeasure({
            read: () => {},
            write: () => {
                view.dom.style.height = `100%`;
            },
        });

        return () => {
            view.destroy();
        };
    });
</script>

<div class="editor-container" bind:this={editorContainer}></div>

<style>
    .editor-container {
        width: 100%;
        height: 100%;
    }
</style>
