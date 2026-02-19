-- name: CreateBusiness :one
INSERT INTO businesses (
   slug, tax_id, company_name, trade_name, description, street, city, province,
   country, zip_code, email, phone_number, whatsapp_number, website
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetBusinesses :many
SELECT * FROM businesses;

-- name: GetBusiness :one
SELECT * FROM businesses WHERE id = $1;

-- name: UpdateBusiness :one
UPDATE businesses SET
  slug = COALESCE(sqlc.narg('slug'), slug),
  tax_id = COALESCE(sqlc.narg('tax_id'), tax_id),
  company_name = COALESCE(sqlc.narg('company_name'), company_name),
  trade_name = COALESCE(sqlc.narg('trade_name'), trade_name),
  description = COALESCE(sqlc.narg('description'), description),
  street = COALESCE(sqlc.narg('street'), street),
  city = COALESCE(sqlc.narg('city'), city),
  province = COALESCE(sqlc.narg('province'), province),
  country = COALESCE(sqlc.narg('country'), country),
  zip_code = COALESCE(sqlc.narg('zip_code'), zip_code),
  email = COALESCE(sqlc.narg('email'), email),
  phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
  whatsapp_number = COALESCE(sqlc.narg('whatsapp_number'), whatsapp_number),
  website = COALESCE(sqlc.narg('website'), website)
WHERE id = $1
RETURNING *;

-- name: DeleteBusiness :exec
DELETE FROM businesses
WHERE id = $1;
