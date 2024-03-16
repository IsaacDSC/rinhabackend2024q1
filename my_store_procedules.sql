CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE result_transaction AS (
                                      balance VARCHAR, "limit" VARCHAR
                                  );


CREATE OR REPLACE FUNCTION insert_transactions(
    in_transaction_client_id INTEGER,
    in_transaction_value VARCHAR(255),
    in_transaction_type VARCHAR(1),
    in_transaction_description VARCHAR(255)
)
    RETURNS result_transaction
    LANGUAGE plpgsql
AS $$
DECLARE
    current_balance VARCHAR;
    account_limit VARCHAR;
    result result_transaction;
BEGIN
    SELECT balance, "limit" INTO current_balance, account_limit FROM clients where id = in_transaction_client_id FOR UPDATE;

    IF in_transaction_type = 'd' THEN
        IF  in_transaction_value::int > account_limit::int or  in_transaction_value::int > current_balance::int  THEN
            RAISE EXCEPTION  'Debit must be positive';
        ELSE
            result.balance := current_balance::int - in_transaction_value::int;
        END IF;
    ELSIF in_transaction_type = 'c' THEN
        result.balance := current_balance::int + in_transaction_value::int;
    ELSE
        RAISE EXCEPTION 'Invalid transaction type:: %', in_transaction_type;
    END IF;

    result.limit := account_limit;

    UPDATE clients SET balance = result.balance WHERE id = in_transaction_client_id;

    INSERT INTO transactions (id, client_id, value, type, description)
    VALUES (gen_random_uuid(), in_transaction_client_id, in_transaction_value, in_transaction_type, in_transaction_description);

    RETURN result;
END;
$$;


select * from clients where id = 1;
select * from transactions where client_id=1;
select gen_random_uuid();

SELECT "balance", "limit" FROM insert_transactions(1,'10000', 'c', 'desc');


SELECT transactions.id, transactions.client_id, transactions.type, transactions.description, transactions.value, transactions.created_at, clients.balance, clients."limit" FROM transactions
                                                                                                                                                                                    JOIN clients on clients.id = transactions.client_id
WHERE transactions.client_id = 1 ORDER BY  transactions.created_at DESC;