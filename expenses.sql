CREATE SCHEMA IF NOT EXISTS expenses;

CREATE TABLE IF NOT EXISTS expenses.currencies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(2) NOT NULL,
    is_main BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS expenses.accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    currency_id INT REFERENCES expenses.currencies(id)
);

CREATE TABLE IF NOT EXISTS expenses.exchange_rate (
    id SERIAL PRIMARY KEY,
    currency_id INT REFERENCES expenses.currencies(id),
    rate_date DATE NOT NULL,
    rate DECIMAL(12,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS expenses.categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS expenses.labels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS expenses.transactions (
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES expenses.accounts(id),
    category_id INT REFERENCES expenses.categories(id),
    currency_id INT REFERENCES expenses.currencies(id),
    amount DECIMAL(12,2) NOT NULL,
    transaction_date DATE NOT NULL,
    note TEXT,
    is_transfer BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS expenses.transaction_labels (
    transaction_id INT REFERENCES expenses.transactions(id),
    label_id INT REFERENCES expenses.labels(id),
    PRIMARY KEY (transaction_id, label_id)
);
