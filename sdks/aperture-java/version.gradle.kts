val snapshot = true

allprojects {
  var ver = "2.24.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
