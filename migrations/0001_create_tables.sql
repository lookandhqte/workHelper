CREATE TABLE Accounts (
    account_id TEXT PRIMARY KEY,
    access_token TEXT,
    refresh_token TEXT,
    expires TIMESTAMP
);

CREATE TABLE Account_Integration (
    account_id TEXT,
    secret_key TEXT,
    client_id TEXT,
    redirect_url TEXT,
    auth_code TEXT,
    FOREIGN KEY (account_id) REFERENCES Accounts(account_id)
);