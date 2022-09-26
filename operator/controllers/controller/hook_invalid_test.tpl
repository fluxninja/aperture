{
  "apiVersion": "fluxninja.com/v1alpha1",
  "kind": "Controller",
  "metadata": {
    "name": "controller-sample"
  },
  "spec": {
    "image": {
      "pullPolicy": "test"
    },
    "config": {
      "etcd": {
        "endpoints": [
          "test.com"
        ]
      },
      "prometheus": {
        "address": "test.com"
      }
    }
  }
}
