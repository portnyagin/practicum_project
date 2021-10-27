package query

const CreateAccount = "INSERT INTO practicum_ut.accounts (id, user_id) VALUES(nextval('seq_account'), $1);"

const updateAccountForUser = "UPDATE practicum_ut.accounts \n" +
	"SET balance=$2, debit=$3, credit=$4 \n" +
	"WHERE  user_id=$1;"

const getAccountBalanceForUser = "select balance, debit from accounts where user_id = $1"

/*
UPDATE practicum_ut.accounts
SET balance=0, debit=0, credit=0
WHERE  user_id=0;

select balance, debit from accounts where user_id = 0*/
