val snapshot = true

allprojects {
  var ver = "2.4.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
