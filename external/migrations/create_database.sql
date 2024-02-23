CREATE TABLE IF NOT EXISTS "clients" (
    "id" SERIAL PRIMARY KEY,
    "balance" VARCHAR(255),
    "limit" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "transactions" (
    "id" VARCHAR(255) PRIMARY KEY,
    "client_id" SERIAL NOT NULL,
    "value" VARCHAR(255) NOT NULL,
    "type" VARCHAR(255) NOT NULL,
    "description" VARCHAR(255),
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
	 CONSTRAINT fk_client FOREIGN KEY(client_id) REFERENCES clients(id)
);