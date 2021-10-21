package ddl

const clrUsers = "drop table if exists users cascade;\n"
const clrOrders = "drop table if exists orders cascade;\n"
const clrAccounts = "drop table if exists accounts cascade;\n"
const clrOperations = "drop table if exists operations cascade;\n"

const ClearDatabaseStructure = clrUsers + clrOrders + clrAccounts + clrOperations
