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

-- name: GetEventsByBusinessProfessionalPatient :many
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
  AND e.professional_id = $2
  AND e.user_id = $3
  AND e.deleted_at IS NULL
ORDER BY
  e.start_date::date DESC,
  e.end_date::time DESC;

-- name: GetEventsFiltered :many
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
    e.created_at,
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
      'role',
      jsonb_build_object('name', r.name, 'value', r.value)
    ),
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
    )
  ) AS event
FROM
  events e
  LEFT JOIN users u ON e.user_id = u.id
  LEFT JOIN roles r ON u.role_id = r.id
  LEFT JOIN users p ON e.professional_id = p.id
  LEFT JOIN professional_profile pp ON p.id = pp.user_id
WHERE
  e.business_id = sqlc.arg (business_id)
  AND (
    sqlc.narg (start_of_day)::timestamp IS NULL
    OR e.start_date >= sqlc.narg (start_of_day)
  )
  AND (
    sqlc.narg (end_of_day)::timestamp IS NULL
    OR e.start_date <= sqlc.narg (end_of_day)
  )
  AND (
    sqlc.narg (patient_id)::uuid IS NULL
    OR u.id = sqlc.narg (patient_id)
  )
  AND (
    sqlc.narg (professional_id)::uuid IS NULL
    OR p.id = sqlc.narg (professional_id)
  )
  AND (
    sqlc.narg (status)::text IS NULL
    OR e.status::text = sqlc.narg (status)
  )
ORDER BY
  CASE
    WHEN sqlc.narg (sort_by) = 'start_date'
    AND sqlc.narg (sort_order) = 'asc' THEN e.start_date
  END ASC,
  CASE
    WHEN sqlc.narg (sort_by) = 'start_date'
    AND sqlc.narg (sort_order) = 'desc' THEN e.start_date
  END DESC,
  CASE
    WHEN sqlc.narg (sort_by) = 'status'
    AND sqlc.narg (sort_order) = 'asc' THEN e.status
  END ASC,
  CASE
    WHEN sqlc.narg (sort_by) = 'status'
    AND sqlc.narg (sort_order) = 'desc' THEN e.status
  END DESC,
  CASE
    WHEN sqlc.narg (sort_by) = 'title'
    AND sqlc.narg (sort_order) = 'asc' THEN e.title
  END ASC,
  CASE
    WHEN sqlc.narg (sort_by) = 'title'
    AND sqlc.narg (sort_order) = 'desc' THEN e.title
  END DESC,
  CASE
    WHEN sqlc.narg (sort_by) = 'professional.firstName'
    AND sqlc.narg (sort_order) = 'asc' THEN p.first_name
  END ASC,
  CASE
    WHEN sqlc.narg (sort_by) = 'professional.firstName'
    AND sqlc.narg (sort_order) = 'desc' THEN p.first_name
  END DESC,
  CASE
    WHEN sqlc.narg (sort_by) = 'user.firstName'
    AND sqlc.narg (sort_order) = 'asc' THEN u.first_name
  END ASC,
  CASE
    WHEN sqlc.narg (sort_by) = 'user.firstName'
    AND sqlc.narg (sort_order) = 'desc' THEN u.first_name
  END DESC,
  e.start_date::date DESC,
  e.start_date::time ASC
LIMIT
  sqlc.arg (query_limit)
OFFSET
  sqlc.arg (query_offset);

-- name: GetEventsFilteredCount :one
SELECT
  COUNT(e.id)::int AS total
FROM
  events e
  LEFT JOIN users u ON e.user_id = u.id
  LEFT JOIN users p ON e.professional_id = p.id
WHERE
  e.business_id = sqlc.arg (business_id)
  AND (
    sqlc.narg (start_of_day)::timestamp IS NULL
    OR e.start_date >= sqlc.narg (start_of_day)
  )
  AND (
    sqlc.narg (end_of_day)::timestamp IS NULL
    OR e.start_date <= sqlc.narg (end_of_day)
  )
  AND (
    sqlc.narg (patient_id)::uuid IS NULL
    OR u.id = sqlc.narg (patient_id)
  )
  AND (
    sqlc.narg (professional_id)::uuid IS NULL
    OR p.id = sqlc.narg (professional_id)
  )
  AND (
    sqlc.narg (status)::text IS NULL
    OR e.status::text = sqlc.narg (status)
  );

-- name: GetEventByID :one
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
    'deletedAt',
    e.deleted_at,
    'professional',
    jsonb_build_object(
      'id',
      p.id,
      'ic',
      p.ic,
      'userName',
      p.user_name,
      'firstName',
      p.first_name,
      'lastName',
      p.last_name,
      'email',
      p.email,
      'password',
      p.password,
      'phoneNumber',
      p.phone_number,
      'roleId',
      p.role_id,
      'businessId',
      p.business_id,
      'refreshToken',
      p.refresh_token,
      'createdAt',
      p.created_at,
      'updatedAt',
      p.updated_at,
      'deletedAt',
      p.deleted_at,
      'professionalProfile',
      jsonb_build_object(
        'id',
        pp.id,
        'businessId',
        pp.business_id,
        'userId',
        pp.user_id,
        'licenseId',
        pp.license_id,
        'professionalPrefix',
        pp.professional_prefix,
        'specialty',
        pp.specialty,
        'workingDays',
        pp.working_days,
        'startHour',
        pp.start_hour,
        'endHour',
        pp.end_hour,
        'slotDuration',
        pp.slot_duration,
        'dailyExceptionStart',
        pp.daily_exception_start,
        'dailyExceptionEnd',
        pp.daily_exception_end,
        'createdAt',
        pp.created_at,
        'updatedAt',
        pp.updated_at,
        'deletedAt',
        pp.deleted_at
      )
    ),
    'user',
    jsonb_build_object(
      'id',
      u.id,
      'ic',
      u.ic,
      'userName',
      u.user_name,
      'firstName',
      u.first_name,
      'lastName',
      u.last_name,
      'email',
      u.email,
      'password',
      u.password,
      'phoneNumber',
      u.phone_number,
      'roleId',
      u.role_id,
      'businessId',
      u.business_id,
      'refreshToken',
      u.refresh_token,
      'createdAt',
      u.created_at,
      'updatedAt',
      u.updated_at,
      'deletedAt',
      u.deleted_at,
      'role',
      jsonb_build_object(
        'id',
        ur.id,
        'name',
        ur.name,
        'value',
        ur.value,
        'description',
        ur.description,
        'createdAt',
        ur.created_at,
        'updatedAt',
        ur.updated_at,
        'deletedAt',
        ur.deleted_at
      )
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
  AND e.id = $2
  AND e.deleted_at IS NULL;

-- name: UpdateEventStatus :one
UPDATE events
SET
  status = $3,
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
