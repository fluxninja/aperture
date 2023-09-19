val snapshot = true

allprojects {
  var ver = "2.15.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
