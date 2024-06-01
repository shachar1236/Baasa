<script>
    import InputAndLabel from "./lib/InputAndLabel.svelte";
    let username = "";
    let password = "";
    let repeat_password = "";
</script>

<main>
    <form action="register" method="post" class="register-form">
        <h1 class="register-heading">Register</h1>
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
        <InputAndLabel
            label="Repeat password"
            input_type="password"
            bind:value={repeat_password}
        />
        {#if password !== repeat_password}
            <p class="error" color="red">Passwords are not equal</p>
        {/if}
        <button
            type="submit"
            class="register-button"
            disabled={password !== repeat_password &&
                username.length > 0 &&
                password.length > 0}>Register</button
        >
    </form>
</main>

<style>
    .register-form {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin: 50px;
    }

    .register-heading {
        margin-bottom: 20px;
        font-size: 3em;
        color: #646cff;
    }

    .register-button {
        background-color: #646cff;
        color: white;
        padding: 1em;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }
</style>
