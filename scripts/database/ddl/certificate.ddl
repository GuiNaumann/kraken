CREATE TABLE certificate
(
    id              SERIAL PRIMARY KEY,
    id_user         INT                                 NOT NULL,
    name            VARCHAR(100)                        NOT NULL,
    image_url       text                                null,
    street          VARCHAR(150),
    address_number  INT,
    district        VARCHAR(70),
    zip_code        VARCHAR(8),
    city            VARCHAR(70),
    is_legal_person BOOLEAN   DEFAULT FALSE,
    cpf             VARCHAR(11),
    cnpj            VARCHAR(14),
    state           INT,
    phone           VARCHAR(15),
    email           VARCHAR(120),
    is_active       BOOLEAN   DEFAULT TRUE,
    status_code     INT       DEFAULT 0                 NOT NULL,
    last_visit_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_certificate_state FOREIGN KEY (state) REFERENCES states (id),
    CONSTRAINT fk_user_certificate FOREIGN KEY (id_user) REFERENCES "user" (id)
);