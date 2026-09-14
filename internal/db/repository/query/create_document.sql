INSERT INTO documents (
    id, 
    owner_id, 
    name, 
    mime_type, 
    has_file, 
    is_public, 
    json_data, 
    file_path
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);