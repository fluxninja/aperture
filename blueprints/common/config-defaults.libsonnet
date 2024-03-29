local selectors_defaults = [{
  control_point: '__REQUIRED_FIELD__',
}];

{
  /**
  * @param (policy.policy_name: string) Name of the policy.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  * @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
  */
  policy: {
    policy_name: '__REQUIRED_FIELD__',
    components: [],
    resources: {
      flow_control: {
        classifiers: [],
      },
    },
  },

  // defaults for the schemas
  selectors_defaults: selectors_defaults,
}
