BEGIN;
INSERT INTO tasks (user_id, description, start_time, end_time, created_at, updated_at)
VALUES

(1, 'Разработать конкурента chatgpt', NOW(), NOW(), NOW(), NOW()),
(1, 'Разработать конкурента google', NOW(), NOW(), NOW(), NOW()),
(2, 'Разработать конкурента youtube', NOW(), NOW(), NOW(), NOW()),
(3, 'Разработать конкурента apple', NOW(), NOW(), NOW(), NOW()),
(4, 'Разработать конкурента android', NOW(), NOW(), NOW(), NOW()),
(5, 'Разработать конкурента amazon', NOW(), NOW(), NOW(), NOW()),
(6, 'Разработать конкурента windows', NOW(), NOW(), NOW(), NOW());

COMMIT;
