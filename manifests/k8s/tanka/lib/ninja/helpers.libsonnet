{
  // Returns the root of the chart folder
  //
  // Usage: `local helm = tanka.helm.new(helpers.helmChartsRoot);`
  // NOTE:
  // Usually, helm renderer is initialized with `std.thisFile`
  // which points to current file.
  // Then the last element on the path is removed to get the directory path.
  // That's why we append '/marker' as this "file".
  helmChartsRoot: std.extVar('projectRoot') + '/marker',
}
