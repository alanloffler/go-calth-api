DELETE FROM businesses
WHERE
  slug = 'system';

DELETE FROM roles
WHERE
  value = 'superadmin';
