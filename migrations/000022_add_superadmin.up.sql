INSERT INTO
  roles (name, value, description)
VALUES
  (
    'Superadmin',
    'superadmin',
    'Acceso global a todos los negocios'
  )
ON CONFLICT (value) DO NOTHING;

INSERT INTO
  businesses (
    slug,
    tax_id,
    company_name,
    trade_name,
    description,
    street,
    city,
    province,
    country,
    zip_code,
    email,
    phone_number
  )
VALUES
  (
    'system',
    '00000000000',
    'System',
    'System',
    'Internal system tenant',
    '-',
    '-',
    '-',
    '-',
    '0000',
    'alanmatiasloffler@gmail.com',
    '0000000000'
  )
ON CONFLICT (slug) DO NOTHING;
