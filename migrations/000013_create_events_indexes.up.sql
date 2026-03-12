CREATE INDEX idx_events_business_status_start ON events (business_id, status, start_date);

CREATE INDEX idx_events_business_title ON events (business_id, title);

CREATE INDEX idx_events_not_deleted ON events (business_id, start_date)
WHERE
  deleted_at IS NULL;

CREATE INDEX idx_users_first_name ON users (first_name);
