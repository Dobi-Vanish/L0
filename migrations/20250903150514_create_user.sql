-- +goose Up
CREATE USER app_user_2 WITH LOGIN PASSWORD 'secure_password_123' NOSUPERUSER NOCREATEDB NOCREATEROLE;

GRANT CONNECT ON DATABASE orders TO app_user_2;
GRANT USAGE ON SCHEMA public TO app_user_2;
GRANT SELECT, INSERT, UPDATE ON orders, deliveries, payments, items TO app_user_2;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO app_user_2;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user_2;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON SEQUENCES TO app_user_2;
REVOKE CREATE ON SCHEMA public FROM app_user_2;


-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
REVOKE ALL ON TABLE orders, deliveries, payments, items FROM app_user_2;
REVOKE USAGE ON SCHEMA public FROM app_user_2;
REVOKE CONNECT ON DATABASE orders FROM app_user_2;
REVOKE USAGE ON ALL SEQUENCES IN SCHEMA public FROM app_user_2;

ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM app_user_2;
ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON SEQUENCES FROM app_user_2;

DROP USER IF EXISTS app_user_2;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd