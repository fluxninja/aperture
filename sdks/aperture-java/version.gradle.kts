val snapshot = true

allprojects {
  var ver = "2.20.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
