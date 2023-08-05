create database paymaster;
create user paymaster with encrypted password 'paymaster';
grant all privileges on database paymaster to paymaster;

INSERT INTO users (address, created_at, updated_at) VALUES ('0x0000000000000000000000000000000000000000', now(), now());
INSERT INTO api_keys (user_id, key, enable, description, created_at, updated_at) VALUES
    (1, '1234567890', true, 'test api key', now(), now());