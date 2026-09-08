SELECT id, login, password_hash
FROM users
where id = $1