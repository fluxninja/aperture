val snapshot = true

allprojects {
  var ver = "0.26.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
