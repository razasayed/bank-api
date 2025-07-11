CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    document_number VARCHAR(20) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS operation_types (
    operation_type_id SERIAL PRIMARY KEY,
    description VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(account_id),
    operation_type_id INT REFERENCES operation_types(operation_type_id),
    amount NUMERIC NOT NULL,
    event_date TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO operation_types (operation_type_id, description) VALUES
(1, 'PURCHASE'),
(2, 'INSTALLMENT PURCHASE'),
(3, 'WITHDRAWAL'),
(4, 'PAYMENT')
ON CONFLICT DO NOTHING;