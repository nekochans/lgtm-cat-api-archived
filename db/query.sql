-- name: ListLgtmImageIds :many
SELECT id FROM lgtm_images;

-- name: ListLgtmImages :many
SELECT id, filename, path FROM lgtm_images
WHERE id = ? OR id = ? OR id = ? OR id = ? OR id = ? OR id = ? OR id = ? OR id = ? OR id = ?;
