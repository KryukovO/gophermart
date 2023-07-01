--
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables 
        WHERE table_name = '__users' AND table_schema = 'public'
    )
    THEN
        ALTER TABLE "__users" RENAME TO "users";
    ELSE
        CREATE TABLE IF NOT EXISTS "users" (
            id INTEGER GENERATED ALWAYS AS IDENTITY,
            login TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            salt TEXT NOT NULL,
            PRIMARY KEY(id)
        );
    END IF;
    --
    --
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables 
        WHERE table_name = '__user_balance' AND table_schema = 'public'
    )
    THEN
        ALTER TABLE "__user_balance" RENAME TO "user_balance";
    ELSE
        CREATE TABLE IF NOT EXISTS "user_balance" (
            id INTEGER GENERATED ALWAYS AS IDENTITY,
            user_id INTEGER NOT NULL UNIQUE,
            balance INTEGER NOT NULL,
            PRIMARY KEY(id),
            FOREIGN KEY(user_id) REFERENCES user(id)
        );
    END IF;
    --
    --
    IF EXISTS (
        SELECT 1
        FROM information_schema.user_defined_types
        WHERE user_defined_type_name = 'balance_operation' AND user_defined_type_schema = 'public';
    )
    THEN
        ALTER TYPE "balance_operation" ADD VALUE IF NOT EXISTS 'withdrawal';
        ALTER TYPE "balance_operation" ADD VALUE IF NOT EXISTS 'refill';
    ELSE
        CREATE TYPE "balance_operation" AS ENUM ('withdrawal', 'refill');
    END IF;
    --
    --
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables 
        WHERE table_name = '__user_balance_log' AND table_schema = 'public'
    )
    THEN
        ALTER TABLE "__user_balance_log" RENAME TO "user_balance_log";
    ELSE
        CREATE TABLE IF NOT EXISTS "user_balance_log" (
            id INTEGER GENERATED ALWAYS AS IDENTITY,
            login TEXT NOT NULL UNIQUE,
            user_id INTEGER NOT NULL,
            created TIMESTAMP WITH TIME ZONE NOT NULL,
            operation "balance_operation" NOT NULL,
            sum INTEGER NOT NULL,
            PRIMARY KEY(id)
        );
    END IF;
    --
    --
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables 
        WHERE table_name = '__orders' AND table_schema = 'public'
    )
    THEN
        ALTER TABLE "__orders" RENAME TO "orders";
    ELSE
        CREATE TABLE IF NOT EXISTS "orders" (
            id INTEGER GENERATED ALWAYS AS IDENTITY,
            user_id INTEGER NOT NULL,
            order_num INTEGER NOT NULL,
            PRIMARY KEY(id),
            FOREIGN KEY(user_id) REFERENCES user(id)
        );
    END IF;
END $$;
