package repository

const CreateUser = "INSERT INTO users " +
	"( id, login, pass) \n" +
	"VALUES($1, $2, $3) returning id;"

const CheckUser = "select 1 from users where login = $1 and active <> 0 and pass=$2;"

const GetUserByLogin = "select id, login, pass from users where active <> 0 and login=$1"

const GetNextUserID = "select nextval('seq_user')"
