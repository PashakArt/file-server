INSERT INTO documents (
    id, 
    owner_id, 
    name, 
    mime, 
    has_file, 
    is_public, 
    json_data, 
    file_path
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);