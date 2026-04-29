CREATE TABLE business_role_permissions (
  business_id UUID NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
  role_id UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
  effect VARCHAR(10) NOT NULL CHECK (effect IN ('grant', 'deny')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (business_id, role_id, permission_id)
);

CREATE INDEX idx_brp_business_role ON business_role_permissions (business_id, role_id);

CREATE INDEX idx_brp_permission ON business_role_permissions (permission_id);
