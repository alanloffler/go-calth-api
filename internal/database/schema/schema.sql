CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  value VARCHAR(100) NOT NULL,
  description VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_role_value ON roles (value);

CREATE TABLE businesses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  slug VARCHAR(50) NOT NULL,
  tax_id VARCHAR(11) NOT NULL,
  company_name VARCHAR(100) NOT NULL,
  trade_name VARCHAR(100) NOT NULL,
  description VARCHAR(100) NOT NULL,
  street VARCHAR(50) NOT NULL,
  city VARCHAR(50) NOT NULL,
  province VARCHAR(50) NOT NULL,
  country VARCHAR(50) NOT NULL,
  zip_code VARCHAR(10) NOT NULL,
  email VARCHAR(100) NOT NULL,
  phone_number VARCHAR(10) NOT NULL,
  whatsapp_number VARCHAR(10),
  website VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ic VARCHAR(8) NOT NULL,
  user_name VARCHAR(100) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL,
  password VARCHAR(100) NOT NULL,
  phone_number VARCHAR(10) NOT NULL,
  role_id UUID NOT NULL REFERENCES roles (id),
  business_id UUID NOT NULL REFERENCES businesses (id),
  refresh_token TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  UNIQUE (business_id, email),
  UNIQUE (business_id, ic),
  UNIQUE (business_id, user_name)
);

CREATE INDEX IF NOT EXISTS idx_users_business ON users (business_id);

CREATE INDEX IF NOT EXISTS idx_users_business_id_id ON users (business_id, id);

CREATE INDEX IF NOT EXISTS idx_users_business_email ON users (business_id, email);

CREATE INDEX IF NOT EXISTS idx_users_business_role ON users (business_id, role_id);

CREATE TABLE patient_profile (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id UUID NOT NULL,
  user_id UUID NOT NULL,
  gender VARCHAR(20) NOT NULL,
  birth_day DATE NOT NULL,
  blood_type VARCHAR(20) NOT NULL,
  weight NUMERIC NOT NULL,
  height NUMERIC NOT NULL,
  emergency_contact_name VARCHAR(50) NOT NULL,
  emergency_contact_phone VARCHAR(10) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT fk_patient_profile_business FOREIGN KEY (business_id) REFERENCES businesses (id) ON DELETE CASCADE,
  CONSTRAINT fk_patient_profile_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_patient_profile_business ON patient_profile (business_id);

CREATE INDEX idx_patient_profile_business_user ON patient_profile (business_id, user_id);

CREATE UNIQUE INDEX uq_patient_profile_business_user ON patient_profile (business_id, user_id);

CREATE TABLE professional_profile (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id UUID NOT NULL,
  user_id UUID NOT NULL,
  license_id VARCHAR NOT NULL,
  professional_prefix VARCHAR NOT NULL,
  specialty VARCHAR NOT NULL,
  working_days VARCHAR NOT NULL,
  start_hour VARCHAR NOT NULL DEFAULT '07:00',
  end_hour VARCHAR NOT NULL DEFAULT '20:00',
  slot_duration VARCHAR NOT NULL DEFAULT '60',
  daily_exception_start VARCHAR,
  daily_exception_end VARCHAR,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT uq_prof_profile_business_license UNIQUE (business_id, license_id),
  CONSTRAINT fk_professional_profile_business FOREIGN KEY (business_id) REFERENCES businesses (id) ON DELETE CASCADE,
  CONSTRAINT fk_professional_profile_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_prof_profile_business ON professional_profile (business_id);

CREATE INDEX idx_prof_profile_business_user ON professional_profile (business_id, user_id);

CREATE TABLE permissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  category VARCHAR(100) NOT NULL,
  action_key VARCHAR(100) NOT NULL,
  description VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_permission_action_key ON permissions (action_key);

CREATE TABLE role_permissions (
  role_id UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (role_id, permission_id)
);

-- // Events //
CREATE TYPE event_status AS ENUM(
  'absent',
  'cancelled',
  'in_progress',
  'pending',
  'present'
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

CREATE INDEX idx_events_business_status_start ON events (business_id, status, start_date);

CREATE INDEX idx_events_business_title ON events (business_id, title);

CREATE INDEX idx_events_not_deleted ON events (business_id, start_date)
WHERE
  deleted_at IS NULL;

CREATE INDEX idx_users_first_name ON users (first_name);

-- // Medical Histories //
CREATE TABLE medical_histories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id UUID NOT NULL,
  user_id UUID NOT NULL,
  professional_id UUID NOT NULL,
  event_id UUID NULL,
  date TIMESTAMPTZ NOT NULL,
  reason VARCHAR NOT NULL,
  recipe BOOLEAN NOT NULL DEFAULT FALSE,
  comments VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ NULL,
  CONSTRAINT fk_mh_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
  CONSTRAINT fk_mh_professional FOREIGN KEY (professional_id) REFERENCES users (id) ON DELETE SET NULL,
  CONSTRAINT fk_mh_event FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE SET NULL
);

CREATE INDEX idx_mh_business_user ON medical_histories (business_id, user_id);

CREATE INDEX idx_mh_business_event ON medical_histories (business_id, event_id);

CREATE INDEX idx_mh_business_user_created ON medical_histories (business_id, user_id, created_at);

-- // Settings //
CREATE TABLE settings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module VARCHAR(50) NOT NULL,
  submodule VARCHAR(50),
  key VARCHAR(50) NOT NULL,
  value VARCHAR(255) NOT NULL,
  title VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
