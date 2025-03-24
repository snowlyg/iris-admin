/**

======== for redis driver ==============
options := &redis.UniversalOptions{
		Addrs:       []string{"127.0.0.1:6379"},
		Password:    "",
		PoolSize:    10,
		IdleTimeout: 300 * time.Second,
		// Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 	conn, err := net.Dial(network, addr)
		// 	if err == nil {
		// 		go func() {
		// 			time.Sleep(5 * time.Second)
		// 			conn.Close()
		// 		}()
		// 	}
		// 	return conn, err
		// },
	}

	err := multi.InitDriver(&multi.Config{
		DriverType:      "redis",
		TokenMaxCount: 10,
		UniversalClient: redis.NewUniversalClient(options)})
	if err != nil {
		panic(err)
	}

======== for local driver ==============
	err := multi.InitDriver(&multi.Config{
		DriverType:      "local"
	if err != nil {
		panic(err)
	}


======== for jwt driver ==============
	err := multi.InitDriver(&multi.Config{
		DriverType:      "jwt",
		HmacSecret:      "your hmac_secret",
	if err != nil {
		panic(err)
	}
*/

package v2
