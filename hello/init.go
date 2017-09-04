package hello

import (
	"expvar"
	"net/http"

	"fmt"
)

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type HelloWorldModule struct {
	cfg       *Config
	something string
	stats     *expvar.Int
}

func NewHelloWorldModule() *HelloWorldModule {

	var cfg Config

	//ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
	//if !ok {
	// when the app is run with -e switch, this message will automatically be redirected to the log file specified
	//	log.Fatalln("failed to read config")
	//}

	// this message only shows up if app is run with -debug option, so its great for debugging
	//logging.Debug.Println("hello init called", cfg.Server.Name)

	return &HelloWorldModule{
		cfg:       &cfg,
		something: "John Doe",
		stats:     expvar.NewInt("rpsStats"),
	}

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
	hlm.stats.Add(1)
	r.ParseForm()
	for key, val := range r.Form {
		fmt.Println(key, val)
	}

	fmt.Println(r.Form)
	name := r.FormValue("name")
	w.Write([]byte("Hello " + name))
}

func (hlm *HelloWorldModule) HelloWebService(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	w.Write([]byte("Halo nama saya" + name))
}
