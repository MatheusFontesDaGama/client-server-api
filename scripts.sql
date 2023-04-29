-- sqlite3 dollar_quote.db
-- .database

CREATE TABLE dollar_quotes(
    id VARCHAR(50) NOT NULL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    code_in VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    high DECIMAL(9, 4) NOT NULL,
    low DECIMAL(9, 4) NOT NULL,
    var_bid DECIMAL (9, 4) NOT NULL,
    pct_change DECIMAL (9, 4) NOT NULL,
    bid DECIMAL (9, 4) NOT NULL,
    ask DECIMAL (9, 4) NOT NULL,
    timestamp VARCHAR NOT NULL,
    create_date DATETIME NOT NULL
);