-- Accept = Request.Auth.Username == "admin"
local user = Query("SELECT * FROM users WHERE username = :username", {
    username = "admin",
})

print(user[1].username)

print("count: " .. Count("users", "username = ? AND id = ?", {"admin", 2}))

Accept = Request.Auth.ID == user[1].id
