package query

const CreateOrder = "INSERT INTO practicum_ut.orders \n" +
	"(id, user_id, num, status, upload_at, updated_at) \n" +
	"VALUES(nextval('seq_order'),  $1, $2, $3, $4,$5);"

const UpdateOrderStatus = "UPDATE practicum_ut.orders \n" +
	"SET  status=$3, updated_at=CURRENT_TIMESTAMP \n" +
	"where user_id=$1 and num =$2 and status!=$3;"

const FindOrdersByUser = "select id, num,user_id, status, upload_at, updated_at from orders where user_id = $1;"

const GetOrderByID = "select id, user_id, num, status, upload_at, updated_at from orders where id = $1;"
const GetOrderByNum = "select id, user_id, num, status, upload_at, updated_at from orders where num = $1;"
