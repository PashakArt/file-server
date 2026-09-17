DELETE FROM documents
WHERE id = $1 AND owner_id = $2
RETURNING has_file, file_path;