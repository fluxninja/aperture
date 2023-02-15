val snapshot = true

allprojects {
  var ver = "0.21.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
