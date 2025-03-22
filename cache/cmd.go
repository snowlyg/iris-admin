package cache

// // InitConfig initialize redis's config file
// func InitConfig() error {
// 	var cover string
// 	if IsExist() {
// 		fmt.Println("Your redis config is initialized , reinitialized redis will cover your redis config.")
// 		fmt.Println("Did you want to do it ?  [Y/N]")
// 		fmt.Scanln(&cover)
// 		switch strings.ToUpper(cover) {
// 		case "Y":
// 		case "N":
// 			return nil
// 		default:
// 		}
// 	} else {
// 		fmt.Println("Redis config file is not exist!")
// 	}

// 	err := Remove()
// 	if err != nil {
// 		return err
// 	}

// 	err = initConfig()
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("redis initialized finished!")
// 	return nil
// }

// func initConfig() error {
// 	var addr, dbPwd string
// 	var db, poolSize int
// 	fmt.Println("Please input your redis addr: ")
// 	fmt.Printf("Redis addr default is '%s'\n", CONFIG.Addr)
// 	fmt.Scanln(&addr)
// 	if addr != "" {
// 		CONFIG.Addr = addr
// 	}

// 	fmt.Println("Please input your redis db: ")
// 	fmt.Printf("Redis db default is '%d'\n", CONFIG.DB)
// 	fmt.Scanln(&db)
// 	if db > 0 {
// 		CONFIG.DB = db
// 	}

// 	fmt.Println("Please input your redis password: ")
// 	fmt.Printf("Redis password default is '%s'\n", CONFIG.Password)
// 	fmt.Scanln(&dbPwd)
// 	if dbPwd != "" {
// 		CONFIG.Password = dbPwd
// 	}

// 	fmt.Println("Please input your redis pool size: ")
// 	fmt.Scanln(&poolSize)
// 	if poolSize > 0 {
// 		CONFIG.PoolSize = poolSize
// 	}
// 	// admin.Init(getViperConfig())
// 	if Instance() == nil {
// 		return ErrRedisInit
// 	}
// 	return nil
// }
