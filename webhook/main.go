package webhook

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var CmdWebhook = &cobra.Command{
	Use:   "run",
	Short: "Start webhook server",
	Long:  "Start webhook server",
	Args:  cobra.ExactArgs(0),
	Run:   main,
}

func main(cmd *cobra.Command, args []string) {
	http.HandleFunc("/parse", handler)
	http.HandleFunc("/", base)
	log.Println("start webhook server")

	log.Fatal(http.ListenAndServe(":8080", nil))

	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func base(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This endpoint has no function")
}
