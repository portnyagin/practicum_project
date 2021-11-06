package query

const CreateAccount = "INSERT INTO accounts (id, user_id) VALUES(nextval('seq_account'), $1);"

const UpdateAccountForUser = "UPDATE accounts \n" +
	"SET balance=$2, debit=$3, credit=$4 \n" +
	"WHERE  user_id=$1;"

const GetAccountForUpdate = "select id, user_id, balance, debit, credit from accounts where user_id = $1 for update"

const GetAccount = "select id, user_id, balance, debit, credit from accounts where user_id = $1"
