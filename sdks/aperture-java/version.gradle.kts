val snapshot = true

allprojects {
  var ver = "1.3.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
