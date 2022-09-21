local patch =
  {
    // Rename commonselectorv1ControlPoint to just ControlPoint
    local controlPoint = super.commonselectorv1ControlPoint,
    commonselectorv1ControlPoint:: null,
    ControlPoint: controlPoint,
  };

{
  v1+: patch,
}
