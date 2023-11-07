{
  // Serialize a basic value to a string, adding double quotes around string values
  serializeValue: function(key, value)
    if (std.isString(value)) then
      '"%(value)s"' % { value: value }
    else
      if (std.isNumber(value) || std.isBoolean(value)) then
        std.toString(value)
      else
        error 'Unsupported key: %(key)s, value type: %(value)s' % { key: key, value: std.type(value) }
  ,
  // Convert a dict to a prometheus filter
  dictToPrometheusFilter: function(dict)
    std.join(
      ',',
      std.map(
        function(key) '%(key)s=%(value)s' % { key: key, value: $.serializeValue(key, dict[key]) },
        std.objectFields(dict),
      )
    ),
}
