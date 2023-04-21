val snapshot = true

allprojects {
  var ver = "1.4.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
