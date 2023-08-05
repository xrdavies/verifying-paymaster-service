# create database paymaster;
# create user paymaster with encrypted password 'paymaster';
# grant all privileges on database paymaster to paymaster;

CREATE TABLE "users" (
    "id" BIGSERIAL,
    "created_at" TIMESTAMPTZ,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "address" VARCHAR(42) DEFAULT NULL,
    PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "users_address_key" ON "users"("address");
