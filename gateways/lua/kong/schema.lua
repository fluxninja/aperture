local typedefs = require "kong.db.schema.typedefs"

return {
    name = "aperture-plugin",
    fields = { {
        config = {
            type = "record",
            fields = { {
                control_point = {
                    type = "string",
                    required = true
                }
            } }
        }
    } }
}
