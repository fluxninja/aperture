val snapshot = true

allprojects {
  var ver = "2.2.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
