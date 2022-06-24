package orm

// func TestPostgreSQL(t *testing.T) {
// 	logger.Init()
//
// 	Init(PostgreSQL, database.Addrs("host=0.0.0.0 port=5432 user=postgres password=123456 dbname=postgres sslmode=disable"))
// 	fmt.Println(defaultDB.gormDB.Name())
// }
//
// func TestTransaction(t *testing.T) {
// 	logger.Init()
//
// 	Init(PostgreSQL, database.Addrs("host=0.0.0.0 port=5432 user=postgres password=123456 dbname=postgres sslmode=disable"))
//
// 	Transaction(DB(),
// 		func(tx *gorm.DB) error {
// 			tx.Table("test1").Updates(nil)
// 			return nil
// 		},
// 		func(tx *gorm.DB) error {
// 			tx.Table("test2").Updates(nil)
// 			return fmt.Errorf("test error")
// 		})
// }
