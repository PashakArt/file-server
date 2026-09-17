SELECT 
    d.id, d.owner_id, d.name, d.mime, d.has_file, 
    d.is_public, d.file_path, d.json_data, d.created_at
FROM documents d
LEFT JOIN documents_users du ON d.id = du.doc_id
WHERE d.id = $1 AND (d.is_public = true OR d.owner_id = $2 OR du.user_id = $2)
LIMIT 1