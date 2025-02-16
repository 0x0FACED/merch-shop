CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$ 
DECLARE 
    i INT := 0;
    user_id INT;
BEGIN
    WHILE i < 100000 LOOP
        INSERT INTO shop.users (username, password_hash) 
        VALUES (
            'user' || i, 
            crypt('pass' || i, gen_salt('bf', 4))
        )
        RETURNING id INTO user_id;

        INSERT INTO shop.wallets (user_id, balance) 
        VALUES (user_id, 1000);

        i := i + 1;
    END LOOP;
END $$;
