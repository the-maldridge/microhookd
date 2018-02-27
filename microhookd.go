package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

var (
	addr = flag.String("bind", "", "Bind address")
	port = flag.Int("port", 8080, "Bind port")
	cmd  = flag.String("cmd", "", "Command to run on poke")

	hooksSinceLastRestart = 0
	errorsSinceLastRestart = 0
)

func runHook(w http.ResponseWriter, r *http.Request) {
	hooksSinceLastRestart++
	log.Println("Running the hook!")
	if err := exec.Command(*cmd).Run(); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "ERROR")
		errorsSinceLastRestart++
		return
	}
	fmt.Fprintf(w, "OK")
}

func help(w http.ResponseWriter, r *http.Request) {
	helpmsg := `Welcome to the hook server!  To fire the hook, you
should call /hook on this server.  Ideally this server should
be firewall protected so that only authorized users can call
it, since its *very* simple and will interpret any request at
all as a reason to invoke the hook!`

	fmt.Fprintf(w, helpmsg)
}

func vars(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "calls-since-last-restart: %d\n", hooksSinceLastRestart)
	fmt.Fprintf(w, "errors-since-last-restart: %d", errorsSinceLastRestart)
}

func main() {
	flag.Parse()
	log.Println("microhookd is starting")
	log.Printf("Will bind on %s:%d", *addr, *port)
	log.Printf("Command will be %s", *cmd)
	http.HandleFunc("/hook", runHook)
	http.HandleFunc("/vars", vars)
	http.HandleFunc("/", help)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *addr, *port), nil))
}
