<script>
    import TextEditor from "../TextEditor.svelte";
    import RulesEditor from "../TextEditor.svelte";

    export let collection;
</script>

<div class="rules-edit-screen">
    <div class="rules-table">
        <nav>
            <ul>
                <li>
                    <button>list</button>
                </li>
                <li>
                    <button>view</button>
                </li>
                <li>
                    <button>create</button>
                </li>
                <li>
                    <button>update</button>
                </li>
                <li>
                    <button>delete</button>
                </li>
            </ul>
        </nav>
        <TextEditor></TextEditor>
    </div>
</div>

<style>
    .rules-edit-screen {
        width: 94vw;
        height: 100vh;
        margin-top: 2rem;
        background-color: #f0f0f0;
        border-radius: 20px;
        padding: 20px;
        transition: all 0.2s ease-in-out;
    }

    .rules-table {
        display: flex;
        height: 100%;
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
        margin-bottom: 1rem;
    }

    li button {
        width: 100%;
    }

    /* remove border outline after clinking on button */
    li button:focus {
        outline: none;
    }
</style>
