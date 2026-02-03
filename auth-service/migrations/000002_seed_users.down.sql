-- +migrate Down
DELETE FROM users
WHERE email IN ('a1@g.com', 'a2@g.com', 'a3@g.com');
