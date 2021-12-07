package conf
import ("os")
type Database struct {
	Type        string `json:"type"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Name        string `json:"name"`
	TablePrefix string `json:"table_prefix"`
	DBFile      string `json:"db_file"`
}
type Config struct {
	Address  string   `json:"address"`
	Port     int      `json:"port"`
	Database Database `json:"database"`
	Https    bool     `json:"https"`
	CertFile string   `json:"cert_file"`
	KeyFile  string   `json:"key_file"`
}

func DefaultConfig() *Config {
	return &Config{
		Address: "0.0.0.0",
		Port:    os.Getenv("POTR"),
		Database: Database{
			Type:        "mongodb",
			Port:        os.Getenv("MONGOPORT"),
			TablePrefix: "x_",
			DBFile:      "data/data.db",
			User:	     os.Getenv("MONGOUSER"),
			Password:  os.Getenv("MONGOPASSWORD"),
			Host: os.Getenv("MONGOHOST"),
		},
	}
}
