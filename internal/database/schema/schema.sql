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
    role_id UUID NOT NULL REFERENCES roles(id),
    business_id UUID NOT NULL REFERENCES businesses(id),
    refresh_token TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    UNIQUE (business_id, email),
    UNIQUE (business_id, ic),
    UNIQUE (business_id, user_name)
);

CREATE INDEX idx_users_business ON users (business_id);
CREATE INDEX idx_users_business_email ON users (business_id, email);
CREATE INDEX idx_users_business_role ON users (business_id, role_id);

CREATE TABLE patient_profile (
  id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id             UUID        NOT NULL,
  user_id                 UUID        NOT NULL,
  gender                  VARCHAR(20) NOT NULL,
  birth_day               DATE        NOT NULL,
  blood_type              VARCHAR(20) NOT NULL,
  weight                  NUMERIC     NOT NULL,
  height                  NUMERIC     NOT NULL,
  emergency_contact_name  VARCHAR(50) NOT NULL,
  emergency_contact_phone VARCHAR(10) NOT NULL,
  created_at              TIMESTAMPTZ NOT NULL    DEFAULT now(),
  updated_at              TIMESTAMPTZ NOT NULL    DEFAULT now(),
  deleted_at              TIMESTAMPTZ,

  CONSTRAINT fk_patient_profile_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_patient_profile_business ON patient_profile (business_id);
CREATE INDEX idx_patient_profile_business_user ON patient_profile (business_id, user_id);
CREATE UNIQUE INDEX uq_patient_profile_business_user ON patient_profile (business_id, user_id);

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
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);
