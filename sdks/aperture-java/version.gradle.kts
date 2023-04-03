val snapshot = true

allprojects {
  var ver = "1.0.2"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
