package vo

type Config struct {
	SecretKey    string `env:"secretKey"`                               // 密钥
	FileKey      string `env:"fileKey"`                                 // 下载文件的密钥
	DeveloperKey string `env:"developerKey"`                            // 开发者密钥
	DbPath       string `env:"dbPath" envDefault:"./cz_db"`             // 数据库路径
	V4File       string `env:"v4File" envDefault:"cz88_public_v4.czdb"` // 数据库文件 v4
	V6File       string `env:"v6File" envDefault:"cz88_public_v6.czdb"` // 数据库文件 v6
}
