package repository

const CreateOrder = "INSERT INTO orders \n" +
	"(id, user_id, num, status, upload_at, updated_at) \n" +
	"VALUES(nextval('seq_order'),  $1, $2, $3, $4,$5);"

const UpdateOrderStatus = "UPDATE orders \n" +
	"SET  status=$2, updated_at=$3 \n" +
	"where id=$1  and status!=$2;"

const FindOrdersByUser = "select id, num,user_id, status, upload_at, updated_at from orders where user_id = $1 order by upload_at"

const GetOrderByID = "select id, user_id, num, status, upload_at, updated_at from orders where id = $1;"
const GetOrderByNum = "select id, user_id, num, status, upload_at, updated_at from orders where num = $1;"

const GetOrderByNumForUpdate = "select id, user_id, num, status, upload_at, updated_at from orders where num = $1 for update"

const FindOrderByStatuses = "select id, user_id, num, status, upload_at, updated_at from orders where status in ($1, $2, $3, $4, $5) limit $6"
