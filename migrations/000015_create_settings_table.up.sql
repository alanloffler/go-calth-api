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
