package ddl

const createUsers = "create table if not exists users (\n" +
	"id numeric primary key,\n" +
	"login varchar not null,\n" +
	"pass varchar  not null\n" +
	");\n" +
	"" +
	"create sequence if not exists seq_user increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by users.id;\n" +
	"create index if not exists user_login_idx on users (login);"

const createOrders = "create table if not exists orders (\n" +
	"id numeric primary key,\n" +
	"user_id numeric  not null,\n" +
	"num  varchar  not null,\n" +
	"status varchar  not null,\n" +
	"upload_at timestamp with time zone  not null,\n" +
	"updated_at timestamp  with time zone\n" +
	");\n" +
	"" +
	"create sequence if not exists seq_order increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by orders.id;\n" +
	"create index if not exists order_user_id_idx on orders (user_id,status );\n" +
	"create index if not exists order_num_idx on orders (num);\n"

const createAccounts = "create table if not exists accounts ( \n" +
	"id numeric primary key,\n" +
	"user_id numeric  not null,\n" +
	"balance numeric  not null default 0,\n" +
	"debit numeric not null default 0,\n" +
	"credit numeric not null default 0\n" +
	");\n" +
	"create sequence if not exists seq_account increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by accounts.id;\n" +
	"create unique index if not exists account_user_id_idx on accounts (user_id );\n"

const createOperations = "create table if not exists operations (\n" +
	"id numeric primary key,\n" +
	"account_id numeric not null,\n" +
	"order_id numeric not null,\n" +
	"operation_type varchar not null ,\n" +
	"amount numeric not null,\n" +
	"processed_at numeric not null\n" +
	");\n" +
	"" +
	"create sequence if not exists seq_operation increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by operations.id;\n" +
	"create index if not exists operation_account_id_idx on operations (account_id );\n" +
	"create index if not exists operation_order_id_idx on operations (order_id );\n"

const CreateDatabaseStructure = createUsers + createAccounts + createOrders + createOperations
