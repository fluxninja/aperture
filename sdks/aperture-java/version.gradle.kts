val snapshot = true

allprojects {
  var ver = "0.3.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
