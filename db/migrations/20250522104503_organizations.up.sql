CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(200)
);

CREATE INDEX idx_organizations_name ON organizations(name);