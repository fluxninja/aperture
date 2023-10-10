---
title: Manually setting feature control points
sidebar_position: 1
slug: manually-setting-feature-control-points-using-python-sdk
keywords:
  - python
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
---

[Aperture Python SDK][pythonsdk] can be used to manually set feature control
points within a Go service.

To do so, first create an instance of ApertureClient:

```python
  from aperture_sdk import ApertureClient

  aperture_client = ApertureClient.new_client(endpoint="localhost:8089")
```

The created instance can then be used to start a flow:

```python
    # business logic produces labels
    labels = {
        "key": "value",
    }

    # start_flow performs a flowcontrol.v1.Check call to Aperture Agent.
    # It returns a Flow or raises an error if any.
    flow = aperture_client.start_flow(
      control_point="AwesomeFeature",
      explicit_labels=labels,
      check_timeout=timedelta(seconds=200),
    )

    # Check if flow check was successful.
    if not flow.success:
        logger.info("Flow check failed - will fail-open")

    # See whether flow was accepted by Aperture Agent.
    if flow.should_run():
        # do actual work
    else:
        # handle flow rejection by Aperture Agent
        flow.set_status(FlowStatus.Error)
    flow.end()
```

You can also use the flow as a context manager:

```python
  with aperture_client.start_flow(
    control_point="AwesomeFeature",
    explicit_labels=labels,
    check_timeout=timedelta(seconds=200),
  ) as flow:
    if flow.should_run():
      # do actual work
      # if you do not call flow.end() explicitly, it will be called automatically
      # when the context manager exits - with the status of the flow
      # depending on whether an error was raised or not
      pass
```

Additionally, you can decorate any function with aperture client. This will skip
running the function if the flow is rejected by Aperture Agent. This might be
helpful to handle specific routes in your service.

```python
  @app.get("/awesome-feature")
  @aperture_client.decorate("AwesomeFeature", check_timeout=timedelta(seconds=200), on_reject=lambda: ("Flow was rejected", 503))
  async def get_awesome_feature_handler():
    return "Flow was accepted", 202
```

For more context on using the Aperture Python SDK to set feature control points,
refer to the [example app][example] available in the repository.

<!--- TODO: Change to pypi package once it is published. --->

[pythonsdk]: https://github.com/fluxninja/aperture/tree/main/sdks/aperture-py
[example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-py/example
