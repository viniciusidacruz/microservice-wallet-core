INSERT INTO clients (id, name, email, created_at, updated_at) VALUES
  ('11111111-1111-1111-1111-111111111111', 'John Doe', 'john.doe@example.com', NOW(), NOW()),
  ('22222222-2222-2222-2222-222222222222', 'Jane Doe', 'jane.doe@example.com', NOW(), NOW());

INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 1000, NOW(), NOW()),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 500, NOW(), NOW());
