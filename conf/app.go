/*
@Time : 2023/12/15 22:13
@Author : chiqing_85
@Software: GoLand
*/
package conf

type Gmeial struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	From string `yaml:"from"`
	Key  string `yaml:"key"`
}
type Web struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}
type Database struct {
	Name     string `yaml:"name"`
	Pas      string `yaml:"pass"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	Type     string `yaml:"type"`
	Prefix   string `yaml:"prefix"`
}
type Token struct {
	Key string `yaml:"key"`
}
type Config struct {
	Web      Web      `yaml:"web"`
	Meial    Gmeial   `yaml:"meial"`
	Database Database `yaml:"database"`
	Token    Token    `yaml:"token"`
}
