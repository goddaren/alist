package conf
import ("os";"strconv")
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
	p1,err:=strconv.Atoi(os.Getenv("PORT"))
	p2,err:=strconv.Atoi(os.Getenv("PGPORT"))
	if err==nil{}
	return &Config{
		Address: "0.0.0.0",
		Port:    p1,
		Database: Database{
			Type:        "postgres",
			Port:        p2,
			TablePrefix: "x_",
			Host: os.Getenv("PGHOST"),
			Name: os.Getenv("PGDATABASE"),
			User:os.Getenv("PGUSER"),
			Host:os.Getenv("PGHOST"),
			Password:os.Getenv("PGPASSWORD")
			
		},
	}
}
