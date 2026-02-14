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
