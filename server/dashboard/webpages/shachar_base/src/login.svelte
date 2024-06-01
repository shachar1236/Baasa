<script>
    import InputAndLabel from "./lib/InputAndLabel.svelte";
    let username;
    let password = "";
</script>

<main>
    <form action="login" method="post" class="login-form">
        <h1 class="login-heading">login</h1>
        <InputAndLabel
            label="Username"
            name="username"
            input_type="text"
            bind:value={username}
        />
        <InputAndLabel
            label="Password"
            name="password"
            input_type="password"
            bind:value={password}
        />
        <button type="submit" class="login-button">login</button>
    </form>
</main>

<style>
    .login-form {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin: 50px;
    }

    .login-heading {
        margin-bottom: 20px;
        font-size: 3em;
        color: #646cff;
    }

    .login-button {
        background-color: #646cff;
        color: white;
        padding: 1em;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }
</style>
