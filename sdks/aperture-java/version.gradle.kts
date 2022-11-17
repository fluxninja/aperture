val snapshot = true

allprojects {
  var ver = "0.4.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
