-- CREATE TYPE visit_status AS ENUM (
--   'approved',
--   'denied'
-- );

CREATE TABLE "employees" (
  "id" bigserial PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "mobile" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "token" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "visitors" (
  "id" bigserial PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "email" varchar NOT NULL,
  "mobile" varchar NOT NULL,
  "company_name" varchar NOT NULL,
  "picture" varchar NOT NULL,
  "sign_in" timestamptz NOT NULL DEFAULT (now()),
  "sign_out" timestamptz,
  "employee_id" bigint DEFAULT null,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "visits" (
  "id" bigserial PRIMARY KEY,
  -- "status" visit_status NOT NULL,
  "status" varchar NOT NULL,
  "reason" varchar NOT NULL,
  "employee_id" bigint DEFAULT null,
  "visitor_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE INDEX ON "employees" ("fullname");

CREATE INDEX ON "employees" ("email");

CREATE INDEX ON "visitors" ("fullname");

CREATE INDEX ON "visitors" ("email");

CREATE INDEX ON "visits" ("status");

CREATE INDEX ON "visits" ("employee_id");

CREATE INDEX ON "visits" ("visitor_id");

CREATE INDEX ON "visits" ("reason");

ALTER TABLE "visitors" ADD FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON DELETE SET DEFAULT;

ALTER TABLE "visits" ADD FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON DELETE SET DEFAULT;

ALTER TABLE "visits" ADD FOREIGN KEY ("visitor_id") REFERENCES "visitors" ("id");
