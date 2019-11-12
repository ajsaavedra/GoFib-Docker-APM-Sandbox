USE gofib;

CREATE TABLE IF NOT EXISTS sequences (
	idx INT(3) NOT NULL,
    fib VARCHAR(24) NOT NULL,
    elapsed VARCHAR(24) NOT NULL,
    PRIMARY KEY(idx)
);
