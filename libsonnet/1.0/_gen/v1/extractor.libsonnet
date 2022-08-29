{
  new():: {
  },
  withAddress(address):: {
    address: address,
  },
  withAddressMixin(address):: {
    address+: address,
  },
  withFrom(from):: {
    from: from,
  },
  withFromMixin(from):: {
    from+: from,
  },
  withJson(json):: {
    json: json,
  },
  withJsonMixin(json):: {
    json+: json,
  },
  withJwt(jwt):: {
    jwt: jwt,
  },
  withJwtMixin(jwt):: {
    jwt+: jwt,
  },
  withPathTemplates(path_templates):: {
    path_templates: path_templates,
  },
  withPathTemplatesMixin(path_templates):: {
    path_templates+: path_templates,
  },
}
