/**
* @param (dry_run: bool) Dynamic configuration for setting dry run mode at runtime without restarting this policy. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. It is useful for observing the behavior of load scheduler without disrupting any real traffic.
*/
{
  dry_run: '__REQUIRED_FIELD__',
}