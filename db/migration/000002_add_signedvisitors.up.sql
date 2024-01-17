CREATE TABLE "signedvisitors" (
  "id" bigserial PRIMARY KEY,
  "visitor_id" bigint NOT NULL,
  "signed_in" timestamptz NOT NULL DEFAULT (now()),
  "signed_out" timestamptz,
  "visit_id" bigint,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

ALTER TABLE "signedvisitors" ADD FOREIGN KEY ("visit_id") REFERENCES "visits" ("id");

CREATE TABLE "admin" (
  "id" bigserial PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "email" varchar NOT NULL,
  "role" varchar, -- admin, superadmin
  "logged_in" timestamptz,
  "logged_out" timestamptz,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "adminlogs" (
  "id" bigserial PRIMARY KEY,
  "admin_id" bigint NOT NULL,
  "logged_in" timestamptz NOT NULL DEFAULT (now()),
  "logged_out" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "adminlogs" ADD FOREIGN KEY ("admin_id") REFERENCES "admin" ("id");