--
BEGIN TRANSACTION;
--
ALTER TABLE "users" RENAME TO "__users";
--
ALTER TABLE "user_balance" RENAME TO "__user_balance";
--
ALTER TABLE "user_balance_log" RENAME TO "__user_balance_log";
--
ALTER TABLE "orders" RENAME TO "__orders";
--
COMMIT TRANSACTION;
