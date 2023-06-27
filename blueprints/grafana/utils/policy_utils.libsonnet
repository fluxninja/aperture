{
  // Serialize a basic value to a string, adding double quotes around string values
  serializeValue: function(value)
    if (std.isString(value)) then
      '"%(value)s"' % { value: value }
    else
      if (std.isNumber(value) || std.isBoolean(value)) then
        std.toString(value)
      else
        error 'Unsupported value type: %(value)s' % { value: std.typeOf(value) }
  ,
  // Convert a dict to a prometheus filter
  dictToPrometheusFilter: function(dict)
    std.join(
      ',',
      std.map(
        function(key) '%(key)s=%(value)s' % { key: key, value: $.serializeValue(dict[key]) },
        std.objectFields(dict),
      )
    ),
}
