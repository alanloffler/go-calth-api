-- name: GetSettings :many
SELECT
  *
FROM
  settings;

-- name: GetSettingsByModule :many
SELECT
  *
FROM
  settings
WHERE
  module = $1;

-- name: UpdateSetting :execrows
UPDATE settings
SET
  module = COALESCE(sqlc.narg ('module'), module),
  submodule = COALESCE(sqlc.narg ('submodule'), submodule),
  key = COALESCE(sqlc.narg ('key'), key),
  value = COALESCE(sqlc.narg ('value'), value),
  title = COALESCE(sqlc.narg ('title'), title),
  updated_at = now()
WHERE
  id = $1;
