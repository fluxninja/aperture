{
  "apiVersion": "fluxninja.com/v1alpha1",
  "kind": "Policy",
  "metadata": {
    "name": "policy-sample"
  },
  "spec": {
    "circuit": {
      "components": [
        {
          "constant": {
            "in_ports": {
              "output": {
                "signal_name": "EMA_LIMIT_MULTIPLIER"
              }
            }
          }
        }
      ]
    },
    "resources": {
      "flux_meters": {
        "service1-demo-app": {
          "selector": {
            "agent_group": "default",
            "control_point": {
              "traffic": "ingress"
            },
            "service": "service1-demo-app.demoapp.svc.cluster.local"
          }
        }
      },
      "classifiers": [
        {
          "rules": {
            "user_type": {
              "extractor": {
                "from": "request.http.headers.user_type"
              }
            },
          },
          "selector": {
            "agent_group": "default",
            "control_point": {
              "traffic": "ingress"
            },
            "service": "service1-demo-app.demoapp.svc.cluster.local"
          }
        }
      ]
    }
  }
}
