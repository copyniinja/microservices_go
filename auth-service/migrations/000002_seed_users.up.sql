-- +migrate Up
INSERT INTO users (
    email,
    first_name,
    last_name,
    password,
    active,
    created_at,
    updated_at
)
VALUES
('ajib@g.com', 'Faiyaz', 'Ahmed', '$2a$12$T52BAkwkk4o0ckg.fMTBB.5Kgx1JxvF5899J6L3Clt7uPQKCl5kcm', TRUE, NOW(), NOW()),
('a2@g.com', 'Rahim', 'Uddin', '$2a$12$HGYPAZFTUOYJmpSJsdd5i.SSce20J3lU4RvnuuG.KdNJrO4AnqyoC', TRUE, NOW(), NOW()),
('a3@g.com', 'Karim', 'Hasan', '$2a$12$HGYPAZFTUOYJmpSJsdd5i.SSce20J3lU4RvnuuG.KdNJrO4AnqyoC', TRUE, NOW(), NOW());
