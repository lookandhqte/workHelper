CREATE TABLE Accounts (
    account_id TEXT PRIMARY KEY NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires TIMESTAMP NOT NULL
);

CREATE TABLE Account_Integration (
    account_id TEXT NOT NULL,
    secret_key TEXT NOT NULL,
    client_id TEXT NOT NULL,
    redirect_url TEXT NOT NULL,
    auth_code TEXT NOT NULL,
    FOREIGN KEY (account_id) REFERENCES Accounts(account_id)
);