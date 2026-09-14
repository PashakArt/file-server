SELECT id
FROM users
WHERE login = ANY($1);