package httptest

// // BeforeTestMainGin
// func BeforeTestMainGin(party func(wi *WebServer), seed func(wi *WebServer, mc *MigrationCmd)) (string, *WebServer) {
// 	dbType := admin.TestDbType
// 	if dbType != "" {
// 		CONFIG.System.DbType = dbType
// 	}
// 	if dbType == "redis" {
// 		if err := cache.Recover(); err != nil {
// 			log.Printf("cache recover fail:%s\n", err.Error())
// 		}
// 	}
// 	if err := Recover(); err != nil {
// 		log.Printf("web recover fail:%s\n", err.Error())
// 	}

// 	node, _ := snowflake.NewNode(1)
// 	uuid := str.Join("gin", "_", node.Generate().String())

// 	CONFIG.DbName = uuid
// 	if user := admin.TestMysqlName; user != "" {
// 		CONFIG.Username = user
// 	}
// 	if pwd := admin.TestMysqlPwd; pwd != "" {
// 		CONFIG.Password = pwd
// 	}
// 	if addr := admin.TestMysqlAddr; addr != "" {
// 		CONFIG.Path = addr
// 	}
// 	CONFIG.LogMode = true
// 	if err := Recover(); err != nil {
// 		log.Printf("databse recover fail:%s\n", err.Error())
// 	}

// 	if Instance() == nil {
// 		fmt.Println("database instance is nil")
// 		return uuid, nil
// 	}

// 	wi := Init()
// 	party(wi)
// 	StartTest(wi)

// 	mc := New()
// 	seed(wi, mc)
// 	err := mc.Migrate()
// 	if err != nil {
// 		fmt.Printf("migrate fail: [%s]", err.Error())
// 		return uuid, nil
// 	}
// 	err = mc.Seed()
// 	if err != nil {
// 		fmt.Printf("seed fail: [%s]", err.Error())
// 		return uuid, nil
// 	}

// 	return uuid, wi
// }

// func AfterTestMain(uuid string, isDelDb bool) {
// 	if isDelDb {
// 		dsn := CONFIG.BaseDsn()
// 		if err := DorpDB(dsn, "mysql", uuid); err != nil {
// 			log.Printf("drop table(%s) on dsn(%s) fail %s\n", uuid, dsn, err.Error())
// 		}
// 	}

// 	if db, _ := Instance().DB(); db != nil {
// 		db.Close()
// 	}

// 	// defer operation.Remove()
// 	// defer casbin.Remove()
// 	// defer Remove()
// 	// defer Remove()
// }
