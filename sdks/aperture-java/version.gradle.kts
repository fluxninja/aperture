val snapshot = true

allprojects {
  var ver = "99.99.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
