

CREATE TYPE strategy_status AS ENUM ('preparing', 'active', 'archived', 'deleted');
CREATE TYPE offer_status AS ENUM ('active', 'archived', 'deleted');
CREATE TYPE subscription_status AS ENUM ('preparing', 'active', 'archived', 'deleted', 'suspended');
CREATE TYPE trade_direction AS ENUM ('buy', 'sell');
CREATE TYPE commission_type AS ENUM ('performance', 'management', 'registration');
CREATE TYPE import_job_type AS ENUM ('trades', 'accounts', 'statistics');
CREATE TYPE import_job_status AS ENUM ('pending', 'running', 'success', 'failed');
CREATE TYPE audit_operation AS ENUM ('insert', 'update', 'delete');
