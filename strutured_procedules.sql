select * from clients;



CREATE OR REPLACE FUNCTION efetuar_transacao(
    clienteIdParam int,
    tipoParam varchar,
    valorParam int,
    descricaoParam varchar
)
    RETURNS TABLE (saldoRetorno int, limiteRetorno int) AS $$
DECLARE
client clients%rowtype;
    novoSaldo int;
    numeroLinhasAfetadas int;
BEGIN
    PERFORM * FROM client where id = clienteIdParam FOR UPDATE;

IF tipoParam = 'd' THEN
        novoSaldo := valorParam * -1;
ELSE
        novoSaldo := valorParam;
END IF;

UPDATE clients
SET balance = balance + novoSaldo
WHERE id = clienteIdParam
  AND (novoSaldo > 0 OR "limit" * -1 <= balance + novoSaldo) RETURNING *
INTO client;

GET DIAGNOSTICS numeroLinhasAfetadas = ROW_COUNT;

IF numeroLinhasAfetadas = 0 THEN
        RAISE EXCEPTION 'Cliente nao possui limite';
END IF;

INSERT INTO transactions (client_id, value, type, description, created_at)
VALUES (clienteIdParam, valorParam, tipoParam, descricaoParam, current_timestamp);


RETURN QUERY SELECT client.balance, client.limit;

END
$$
LANGUAGE plpgsql;



call efetuar_transacao(1, 'd', 2000, 'descricao');


CREATE TYPE my_type AS (
    balance VARCHAR, "limit" VARCHAR
    );

CREATE OR REPLACE PROCEDURE insert_transaction(
    in_transaction_client_id INTEGER,
    in_transaction_value VARCHAR(255),
    in_transaction_type VARCHAR(1),
    in_transaction_description VARCHAR(255)
)RETURNS my_type AS $$
DECLARE
client  clients%rowtype;
    new_transaction my_type;
BEGIN
    PERFORM * FROM clients where id = in_transaction_client_id;

INSERT INTO transactions (id,client_id, "value", "type", description)
VALUES ('123',in_transaction_client_id,in_transaction_value, in_transaction_type, in_transaction_description);

new_transaction.balance := client.balance;
    new_transaction.limit := client."limit";

RETURN new_transaction;
END;
$$  LANGUAGE plpgsql;

-- Step 3: Call the stored procedure with your values
CALL insert_transaction(1,'10000', 'credit', 'desc');

select * from transactions;