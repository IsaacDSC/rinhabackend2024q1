job "application-job" {
  type = "service"

  group "application-group" {
    count = 1

    scaling {
      enabled = true
      min     = 1
      max     = 10

      policy {
        check "96pct" {
          strategy "app-sizing-percentile" {
            percentile = "70"
          }
        }
      }
    }

    network {
      port "ingress" {
        to = 80
        static = 8080
      }
    }

    restart {
     attempts = 10
     delay    = "30s"
     mode     = "fail"
    }

    service {
      name     = "rinhabackend"
      port     = "ingress"
      provider = "nomad"

      // check_restart {
      //   limit           = 3
      //   grace           = "10s"
      //   ignore_warnings = false
      // }

      // check {
      //   name     = "healtcheck"
      //   type     = "http"
      //   path     = "/api/v1/status"
      //   interval = "5s"
      //   timeout  = "1s"
      //   method   = "GET"
      //   header {
      //     Authorization = ["Basic ZWxhc3RpYzpjaGFuZ2VtZQ=="]
      //   }
      // }
    }

    // update {
    //   max_parallel     = 1
    //   canary           = 1
    //   min_healthy_time = "30s"
    //   auto_revert      = true
    //   auto_promote     = false
    // }

    task "application-task" {
      driver = "docker"
      
      env {
        DB_HOST=pgdb
        DB_PORT=4321
      }

      config {
        image = "isaacdsc/rinhabackend:latest"
        ports = ["application"]
      }
      resources {
        cpu    = 500 # MHz
        memory = 500 # MB
      }
    }

    task "ingress-application-task" {
      driver = "docker"

      config {
        image = "nginx:latest"
        ports = ["ingress"]
        volumes = [
            "local:/etc/nginx/conf.d",
        ]
      }

      resources {
        cpu    = 50 # MHz
        memory = 50 # MB
      }

      template {
        data = <<EOF
        server {
          listen 80;

          location / {
              proxy_pass http://localhost:3000;
          }
        }
        EOF

        destination   = "local/load-balancer.conf"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }
    }

  }
}