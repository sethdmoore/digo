package types

type Config struct {
	Email     string `envconfig:"user" json:"email"`
	Password  string `envconfig:"pass" json:"password"`
	Trigger   string `envconfig:"trigger" json:"trigger"`
	Interface string `envconfig:"interface" json:"interface"`
	Guild     string `envconfig:"server_id" json:"server_id"`
	//Token    string `json:"token"`
}
