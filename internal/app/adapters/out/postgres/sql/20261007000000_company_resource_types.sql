-- +goose Up
CREATE TABLE company_resource_types (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    icon_key VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (id, company_id)
);

CREATE UNIQUE INDEX company_resource_types_company_name_lower_idx
    ON company_resource_types (company_id, lower(name));

ALTER TABLE resources
    ADD COLUMN company_resource_type_id UUID;

ALTER TABLE resources
    ADD CONSTRAINT resources_company_resource_type_fkey
    FOREIGN KEY (company_resource_type_id, company_id)
    REFERENCES company_resource_types (id, company_id);

ALTER TABLE resources
    ADD CONSTRAINT resources_company_resource_type_check
    CHECK ((type = 'custom') = (company_resource_type_id IS NOT NULL));

CREATE INDEX resources_company_resource_type_idx
    ON resources (company_resource_type_id);

-- +goose Down
DROP INDEX resources_company_resource_type_idx;
ALTER TABLE resources DROP CONSTRAINT resources_company_resource_type_check;
ALTER TABLE resources DROP CONSTRAINT resources_company_resource_type_fkey;
ALTER TABLE resources DROP COLUMN company_resource_type_id;
DROP TABLE company_resource_types;
