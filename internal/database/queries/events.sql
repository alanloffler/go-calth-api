-- name: CreateEvent :one
INSERT INTO
  events (
    title,
    start_date,
    end_date,
    business_id,
    professional_id,
    user_id
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  *;

-- name: GetEventsByProfessionalID :many
SELECT
  jsonb_build_object(
    'id',
    e.id,
    'title',
    e.title,
    'startDate',
    e.start_date,
    'endDate',
    e.end_date,
    'businessId',
    e.business_id,
    'professionalId',
    e.professional_id,
    'userId',
    e.user_id,
    'status',
    e.status,
    'createdAt',
    e.created_at,
    'updatedAt',
    e.updated_at,
    'professional',
    jsonb_build_object(
      'id',
      p.id,
      'firstName',
      p.first_name,
      'lastName',
      p.last_name,
      'ic',
      p.ic,
      'role',
      jsonb_build_object('name', pr.name, 'value', pr.value),
      'professionalProfile',
      jsonb_build_object('professionalPrefix', pp.professional_prefix)
    ),
    'user',
    jsonb_build_object(
      'id',
      u.id,
      'firstName',
      u.first_name,
      'lastName',
      u.last_name,
      'email',
      u.email,
      'phoneNumber',
      u.phone_number,
      'ic',
      u.ic,
      'role',
      jsonb_build_object('name', ur.name, 'value', ur.value)
    )
  ) AS event
FROM
  events e
  LEFT JOIN users u ON u.id = e.user_id
  LEFT JOIN roles ur ON ur.id = u.role_id
  LEFT JOIN users p ON p.id = e.professional_id
  LEFT JOIN roles pr ON pr.id = p.role_id
  LEFT JOIN professional_profile pp ON pp.user_id = p.id
WHERE
  e.business_id = $1
  AND e.professional_id = $2
  AND e.deleted_at IS NULL
ORDER BY
  e.start_date;

-- name: GetEventsByBusinessID :many
SELECT
  jsonb_build_object(
    'id',
    e.id,
    'title',
    e.title,
    'startDate',
    e.start_date,
    'endDate',
    e.end_date,
    'businessId',
    e.business_id,
    'userId',
    e.user_id,
    'status',
    e.status,
    'createdAt',
    e.created_at,
    'updatedAt',
    e.updated_at,
    'professional',
    jsonb_build_object(
      'id',
      p.id,
      'firstName',
      p.first_name,
      'lastName',
      p.last_name,
      'ic',
      p.ic,
      'professionalProfile',
      jsonb_build_object('professionalPrefix', pp.professional_prefix)
    ),
    'user',
    jsonb_build_object(
      'id',
      u.id,
      'firstName',
      u.first_name,
      'lastName',
      u.last_name,
      'email',
      u.email,
      'phoneNumber',
      u.phone_number,
      'ic',
      u.ic,
      'role',
      jsonb_build_object('name', ur.name, 'value', ur.value)
    )
  ) AS event
FROM
  events e
  LEFT JOIN users u ON u.id = e.user_id
  LEFT JOIN roles ur ON ur.id = u.role_id
  LEFT JOIN users p ON p.id = e.professional_id
  LEFT JOIN professional_profile pp ON pp.user_id = p.id
WHERE
  e.business_id = $1
  AND e.deleted_at IS NULL
ORDER BY
  e.start_date::date DESC,
  e.end_date::time DESC
LIMIT
  $2;

-- name: GetProfessionalEventsByDay :many
SELECT
  jsonb_build_object(
    'id',
    e.id,
    'title',
    e.title,
    'startDate',
    e.start_date,
    'endDate',
    e.end_date,
    'businessId',
    e.business_id,
    'professionalId',
    e.professional_id,
    'userId',
    e.user_id,
    'status',
    e.status,
    'createdAt',
    e.created_at,
    'updatedAt',
    e.updated_at,
    'professional',
    jsonb_build_object(
      'id',
      p.id,
      'firstName',
      p.first_name,
      'lastName',
      p.last_name,
      'ic',
      p.ic,
      'role',
      jsonb_build_object('name', pr.name, 'value', pr.value),
      'professionalProfile',
      jsonb_build_object('professionalPrefix', pp.professional_prefix)
    ),
    'user',
    jsonb_build_object(
      'id',
      u.id,
      'firstName',
      u.first_name,
      'lastName',
      u.last_name,
      'email',
      u.email,
      'phoneNumber',
      u.phone_number,
      'ic',
      u.ic,
      'role',
      jsonb_build_object('name', ur.name, 'value', ur.value)
    )
  ) AS event
FROM
  events e
  LEFT JOIN users u ON u.id = e.user_id
  LEFT JOIN roles ur ON ur.id = u.role_id
  LEFT JOIN users p ON p.id = e.professional_id
  LEFT JOIN roles pr ON pr.id = p.role_id
  LEFT JOIN professional_profile pp ON pp.user_id = p.id
WHERE
  e.business_id = $1
  AND e.professional_id = $2
  AND e.start_date >= $3
  AND e.start_date <= $4
  AND e.deleted_at IS NULL
ORDER BY
  e.start_date;

-- name: GetProfessionalEventsByDayArray :many
SELECT
  e.start_date
FROM
  events e
WHERE
  e.business_id = $1
  AND e.professional_id = $2
  AND e.start_date >= $3
  AND e.start_date <= $4
  AND e.deleted_at IS NULL
ORDER BY
  e.start_date;

-- name: GetEventByID :one
SELECT
  *
FROM
  events
WHERE
  id = $1;

-- name: UpdateEventStatus :one
UPDATE events
SET
  status = $3,
  updated_at = now()
WHERE
  id = $1
  AND business_id = $2
  AND deleted_at IS NULL
RETURNING
  id,
  title,
  start_date,
  end_date,
  business_id,
  professional_id,
  user_id,
  status,
  created_at,
  updated_at,
  deleted_at;

-- name: UpdateEvent :one
UPDATE events
SET
  title = COALESCE(sqlc.narg ('title'), title),
  start_date = COALESCE(sqlc.narg ('start_date'), start_date),
  end_date = COALESCE(sqlc.narg ('end_date'), end_date),
  professional_id = COALESCE(sqlc.narg ('professional_id'), professional_id),
  user_id = COALESCE(sqlc.narg ('user_id'), user_id),
  status = COALESCE(sqlc.narg ('status'), status),
  updated_at = now()
WHERE
  business_id = $1
  AND id = $2
  AND deleted_at IS NULL
RETURNING
  id,
  title,
  start_date,
  end_date,
  business_id,
  professional_id,
  user_id,
  status,
  created_at,
  updated_at,
  deleted_at;
