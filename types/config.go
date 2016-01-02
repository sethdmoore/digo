package types

type Config struct {
	Email     string `envconfig:"user" json:"email" toml:"email"`
	Password  string `envconfig:"pass" json:"password" toml:"password"`
	Trigger   string `envconfig:"trigger" json:"trigger" toml:"trigger"`
	Interface string `envconfig:"interface" json:"interface" toml:"interface"`
	Guild     string `envconfig:"server_id" json:"server_id" toml:"server_id"`
	//Token    string `json:"token"`
}
