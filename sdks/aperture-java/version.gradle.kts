val snapshot = true

allprojects {
  var ver = "0.5.0"
  if (snapshot) {
    ver += "-SNAPSHOT"
  }
  version = ver
}
