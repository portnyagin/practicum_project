package query

const CreateUser = "INSERT INTO users " +
	"( id, login, pass) \n" +
	"VALUES(nextval('seq_user'), $1, $2) returning id;"

const CheckUser = "select 1 from users where login = $1 and active <> 0 and pass=$2;"
