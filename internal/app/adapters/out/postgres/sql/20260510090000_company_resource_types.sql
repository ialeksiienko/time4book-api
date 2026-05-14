-- +goose Up
ALTER TABLE resources ALTER COLUMN type TYPE VARCHAR(32);

CREATE TABLE company_resource_types (
    id UUID PRIMARY KEY NOT NULL,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    icon_key VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX company_resource_types_company_name_lower_idx
ON company_resource_types (company_id, lower(name));

ALTER TABLE resources
ADD COLUMN company_resource_type_id UUID REFERENCES company_resource_types(id) ON DELETE SET NULL;

CREATE INDEX idx_resources_company_resource_type_id ON resources(company_resource_type_id);

-- +goose Down
DROP INDEX IF EXISTS idx_resources_company_resource_type_id;
ALTER TABLE resources DROP COLUMN IF EXISTS company_resource_type_id;
DROP TABLE IF EXISTS company_resource_types;
ALTER TABLE resources ALTER COLUMN type TYPE VARCHAR(20);
