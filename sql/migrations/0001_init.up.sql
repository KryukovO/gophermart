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
            id BIGINT GENERATED ALWAYS AS IDENTITY,
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
            id BIGINT GENERATED ALWAYS AS IDENTITY,
            user_id BIGINT NOT NULL UNIQUE,
            balance DOUBLE PRECISION NOT NULL CHECK (balance >= 0),
            PRIMARY KEY(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        );
    END IF;
    --
    --
    IF EXISTS (
        SELECT 1
        FROM pg_type t 
        JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace 
        WHERE n.nspname = 'public' AND t.typname = 'balance_operation'
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
            id BIGINT GENERATED ALWAYS AS IDENTITY,
            user_id BIGINT NOT NULL,
            processed TIMESTAMP WITH TIME ZONE NOT NULL,
            operation "balance_operation" NOT NULL,
            order_num TEXT NOT NULL,
            sum DOUBLE PRECISION NOT NULL,
            PRIMARY KEY(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        );
    END IF;
    --
    --
    IF EXISTS (
        SELECT 1
        FROM pg_type t 
        JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace 
        WHERE n.nspname = 'public' AND t.typname = 'order_status'
    )
    THEN
        ALTER TYPE "order_status" ADD VALUE IF NOT EXISTS 'NEW';
        ALTER TYPE "order_status" ADD VALUE IF NOT EXISTS 'PROCESSING';
        ALTER TYPE "order_status" ADD VALUE IF NOT EXISTS 'INVALID';
        ALTER TYPE "order_status" ADD VALUE IF NOT EXISTS 'PROCESSED';
    ELSE
        CREATE TYPE "order_status" AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
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
            id BIGINT GENERATED ALWAYS AS IDENTITY,
            user_id BIGINT NOT NULL,
            order_num TEXT NOT NULL UNIQUE,
            status order_status NOT NULL,
            accrual DOUBLE PRECISION,
            uploaded TIMESTAMP WITH TIME ZONE NOT NULL,
            PRIMARY KEY(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        );
    END IF;
END $$;
