SELECT du.doc_id, u.login
FROM documents_users du
JOIN users u ON du.user_id = u.id
WHERE du.doc_id = ANY($1)