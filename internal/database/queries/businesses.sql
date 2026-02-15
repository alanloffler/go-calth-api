-- name: CreateBusiness :one
INSERT INTO businesses (
   slug, tax_id, company_name, trade_name, description, street, city, province,
   country, zip_code, email, phone_number, whatsapp_number, website
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;
