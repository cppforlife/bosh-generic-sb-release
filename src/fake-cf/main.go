package main

import (
  "fmt"
  "strings"
  "net/http"
)

type routes struct {}

func main() {
  http.ListenAndServe(":8081", routes{})
}

func (r routes) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  switch {
    case req.URL.Path == "/v2/info":
      r.serveInfo(w, req)
    case req.URL.Path == "/v2/services":
      r.serveServices(w, req)
    case strings.HasPrefix(req.URL.Path, "/v2/service_plans/"):
      r.serveServicePlan(w, req)
    case strings.HasPrefix(req.URL.Path, "/v2/service_instances/"):
      r.serveServiceInstance(w, req)
    default:
      r.notImplemented(w, req)
  }
}

func (r routes) serveInfo(w http.ResponseWriter, _ *http.Request) {
  fmt.Printf("info\n")

  w.Header().Add("Content-Type", "application/json")

  fmt.Fprintf(w, `{
    "name": "",
    "build": "",
    "support": "https://support.run.pivotal.io",
    "version": 0,
    "description": "Cloud Foundry sponsored by Pivotal",
    "authorization_endpoint": "https://login.run.pivotal.io",
    "token_endpoint": "https://uaa.run.pivotal.io",
    "min_cli_version": "6.22.0",
    "min_recommended_cli_version": "latest",
    "api_version": "2.92.0",
    "app_ssh_endpoint": "ssh.run.pivotal.io:2222",
    "app_ssh_host_key_fingerprint": "e7:13:4e:32:ee:39:62:df:54:41:d7:f7:8b:b2:a7:6b",
    "app_ssh_oauth_client": "ssh-proxy",
    "doppler_logging_endpoint": "wss://doppler.run.pivotal.io:443",
    "routing_endpoint": "https://api.run.pivotal.io/routing"
  }`)
}

func (r routes) serveServices(w http.ResponseWriter, _ *http.Request) {
  fmt.Printf("services\n")

  w.Header().Add("Content-Type", "application/json")

  fmt.Fprintf(w, `{
    "total_results": 29,
    "total_pages": 1,
    "prev_url": null,
    "next_url": null,
    "resources": []
  }`)
}

func (r routes) serveServicePlan(w http.ResponseWriter, _ *http.Request) {
  fmt.Printf("service_plan\n")

  w.Header().Add("Content-Type", "application/json")

  // todo stateful?
  fmt.Fprintf(w, `{
    "entity": {
      "unique_id": "service_plan_id_unique_id"
    }
  }`)
}

func (r routes) serveServiceInstance(w http.ResponseWriter, _ *http.Request) {
  fmt.Printf("service_instance\n")

  w.Header().Add("Content-Type", "application/json")

  fmt.Fprintf(w, `{
    "entity": {
      "last_operation": {
         "type": "create",
         "state": "succeeded",
         "description": "Instance provisioning completed",
         "updated_at": "2016-07-13T16:49:05Z",
         "created_at": "2016-07-13T16:47:03Z"
      },
      "service_plan_url": "/v2/service_plans/service_plan_id"
    }
  }`)
}

func (r routes) notImplemented(w http.ResponseWriter, req *http.Request) {
  fmt.Printf("route not implemented '%s'\n", req.URL.Path)

  w.WriteHeader(http.StatusInternalServerError)
  fmt.Fprintf(w, "not implemented")
}
