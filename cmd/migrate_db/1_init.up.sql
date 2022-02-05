CREATE TABLE auth_tokens
(
    id                    uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    client_id             varchar(255) NOT NULL,
    user_id               varchar(255),
    redirect_uri          varchar(255),
    scope                 varchar(255),

    code                  varchar(255) UNIQUE,
    code_challenge        varchar(255),
    code_challenge_method varchar(255),
    code_create_at        timestamp,
    code_expires_at       timestamp,

    access                varchar(255) UNIQUE,
    access_create_at      timestamp,
    access_expires_at     timestamp,

    refresh               varchar(255) UNIQUE,
    refresh_create_at     timestamp,
    refresh_expires_at    timestamp
);

CREATE TABLE auth_clients
(
    id      varchar(255) NOT NULL PRIMARY KEY,
    secret  varchar(255) NOT NULL,
    domain  varchar(255) NOT NULL,
    user_id varchar(255) NOT NULL
);

INSERT INTO auth_clients(id, secret, domain, user_id)
VALUES ('e4212aad-79ea-49cd-bf92-cd102806b68f', 'changeme', 'http://localhost', 'tinyfluffs');
