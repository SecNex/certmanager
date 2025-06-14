package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/secnex/certmanager/database"
	"github.com/secnex/certmanager/logger"
	"github.com/secnex/certmanager/manager"
	"golang.org/x/crypto/acme/autocert"
)

type ApiServer struct {
	Host      *string
	Port      *int
	Databases map[string]*database.Database
	Manager   *manager.Manager
	Router    *mux.Router
}

func NewApiServer(host *string, port *int, manager *manager.Manager) *ApiServer {
	databases := make(map[string]*database.Database)
	databases["certmanager"] = database.NewConnection("localhost", 5432, "postgres", "postgres", "cert")
	return &ApiServer{
		Host:      host,
		Port:      port,
		Databases: databases,
		Manager:   manager,
		Router:    mux.NewRouter(),
	}
}

func (s *ApiServer) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *ApiServer) Start() {
	https := true

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Println("DOMAIN is not set, using HTTP only!")
		https = false
	}

	addr := ":443"
	if !https {
		addr = ":80"
	}

	if s.Host != nil && s.Port != nil {
		addr = fmt.Sprintf("%s:%d", *s.Host, *s.Port)
	}

	log.Printf("Starting api server on %s...", addr)

	s.Router.HandleFunc("/healthz", s.Healthz).Methods(http.MethodGet)

	if https {
		s.StartHTTPS(logger.LogHTTPRequest(s.Router), addr, domain)
	} else {
		s.StartHTTP(logger.LogHTTPRequest(s.Router))
	}

}

func (s *ApiServer) StartHTTP(handler http.Handler) {
	server := &http.Server{
		Addr:    ":80",
		Handler: handler,
	}

	log.Fatal(server.ListenAndServe())
}

func (s *ApiServer) StartHTTPS(handler http.Handler, addr string, domain string) {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist(domain),
	}

	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
		Handler: handler,
	}

	go func() {
		http.ListenAndServe(":80", m.HTTPHandler(nil))
	}()

	log.Fatal(server.ListenAndServeTLS("", ""))
}
