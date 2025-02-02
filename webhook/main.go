package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog/v2"
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

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func main(cmd *cobra.Command, args []string) {
	http.HandleFunc("/allow", serveAlwaysAllow)

	log.Println("start webhook server")

	log.Fatal(http.ListenAndServeTLS(":443", "/app/certs/tls.crt", "/app/certs/tls.key", nil))

	//log.Fatal(http.ListenAndServe(":8080", nil))

	return
}

type admitv1Func func(v1.AdmissionReview) *v1.AdmissionResponse

type admitHandler struct {
	v1 admitv1Func
}

func serveAlwaysAllow(w http.ResponseWriter, r *http.Request) {
	serve(w, r, admitHandler{v1: alwaysAllowV1})
}

func init() {
	v1.AddToScheme(scheme)
}

func serve(w http.ResponseWriter, r *http.Request, admit admitHandler) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))
	deserializer := codecs.UniversalDeserializer()
	obj, gvk, err := deserializer.Decode(body, nil, nil)
	if err != nil {
		msg := fmt.Sprintf("Request could not be decoded: %v", err)
		klog.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var responseObj runtime.Object

	switch *gvk {
	case v1.SchemeGroupVersion.WithKind("AdmissionReview"):
		requestedAdmissionReview, ok := obj.(*v1.AdmissionReview)
		if !ok {
			klog.Errorf("Expected v1.AdmissionReview but got %#v", obj)
			return
		}
		responseAdmissionReview := &v1.AdmissionReview{}
		responseAdmissionReview.SetGroupVersionKind(*gvk)
		responseAdmissionReview.Response = admit.v1(*requestedAdmissionReview)
		responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		responseObj = responseAdmissionReview
	default:
		klog.Errorf("unknown kind: %v", gvk)
		http.Error(w, fmt.Sprintf("unknown kind: %v", gvk), http.StatusBadRequest)
		return
	}

	responseBytes, err := json.Marshal(responseObj)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseBytes); err != nil {
		klog.Error(err)
	}
}

func alwaysAllowV1(ar v1.AdmissionReview) *v1.AdmissionResponse {
	log.Println("Allow request")
	return &v1.AdmissionResponse{
		Allowed: true,
		Result:  &metav1.Status{Message: "This admission controller allows all requests"},
	}
}
