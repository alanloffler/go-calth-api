CREATE TYPE event_status AS ENUM(
  'absent',
  'attended',
  'cancelled',
  'in_progress',
  'pending'
);

CREATE TABLE events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  business_id UUID NOT NULL,
  professional_id UUID NOT NULL,
  user_id UUID NOT NULL,
  status event_status NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT fk_events_professional FOREIGN KEY (professional_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_events_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT chk_events_date_order CHECK (end_date > start_date)
);

CREATE INDEX idx_events_business_start ON events (business_id, start_date);

CREATE INDEX idx_events_business_professional_start ON events (business_id, professional_id, start_date);

CREATE INDEX idx_events_business_user_start ON events (business_id, user_id, start_date);
