val snapshot = true

allprojects {
  var ver = "2.5.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
