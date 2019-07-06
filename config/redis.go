package config

// Redis配置构造
type RedisConfig struct {
	Address  string //HOST:PORT
	Password string //数据库密码
	DB       int    //选择的数据库
	MinIdle  int    //最小空闲数
	MaxRetry int    //重试最大次数
	PoolSize int    //连接池大小
}
