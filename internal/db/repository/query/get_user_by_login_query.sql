SELECT id, login, password_hash
FROM users
where login = $1