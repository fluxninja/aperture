// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/label_matcher.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public interface MatchExpressionOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.MatchExpression)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * The expression negates the result of subexpression.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression not = 1 [json_name = "not"];</code>
   * @return Whether the not field is set.
   */
  boolean hasNot();
  /**
   * <pre>
   * The expression negates the result of subexpression.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression not = 1 [json_name = "not"];</code>
   * @return The not.
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchExpression getNot();
  /**
   * <pre>
   * The expression negates the result of subexpression.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression not = 1 [json_name = "not"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchExpressionOrBuilder getNotOrBuilder();

  /**
   * <pre>
   * The expression is true when all subexpressions are true.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression.List all = 2 [json_name = "all"];</code>
   * @return Whether the all field is set.
   */
  boolean hasAll();
  /**
   * <pre>
   * The expression is true when all subexpressions are true.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression.List all = 2 [json_name = "all"];</code>
   * @return The all.
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchExpression.List getAll();
  /**
   * <pre>
   * The expression is true when all subexpressions are true.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression.List all = 2 [json_name = "all"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchExpression.ListOrBuilder getAllOrBuilder();

  /**
   * <pre>
   * The expression is true when any subexpression is true.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression.List any = 3 [json_name = "any"];</code>
   * @return Whether the any field is set.
   */
  boolean hasAny();
  /**
   * <pre>
   * The expression is true when any subexpression is true.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression.List any = 3 [json_name = "any"];</code>
   * @return The any.
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchExpression.List getAny();
  /**
   * <pre>
   * The expression is true when any subexpression is true.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchExpression.List any = 3 [json_name = "any"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchExpression.ListOrBuilder getAnyOrBuilder();

  /**
   * <pre>
   * The expression is true when label with given name exists.
   * </pre>
   *
   * <code>string label_exists = 4 [json_name = "labelExists"];</code>
   * @return Whether the labelExists field is set.
   */
  boolean hasLabelExists();
  /**
   * <pre>
   * The expression is true when label with given name exists.
   * </pre>
   *
   * <code>string label_exists = 4 [json_name = "labelExists"];</code>
   * @return The labelExists.
   */
  java.lang.String getLabelExists();
  /**
   * <pre>
   * The expression is true when label with given name exists.
   * </pre>
   *
   * <code>string label_exists = 4 [json_name = "labelExists"];</code>
   * @return The bytes for labelExists.
   */
  com.google.protobuf.ByteString
      getLabelExistsBytes();

  /**
   * <pre>
   * The expression is true when label value equals given value.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.EqualsMatchExpression label_equals = 5 [json_name = "labelEquals"];</code>
   * @return Whether the labelEquals field is set.
   */
  boolean hasLabelEquals();
  /**
   * <pre>
   * The expression is true when label value equals given value.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.EqualsMatchExpression label_equals = 5 [json_name = "labelEquals"];</code>
   * @return The labelEquals.
   */
  com.fluxninja.generated.aperture.policy.language.v1.EqualsMatchExpression getLabelEquals();
  /**
   * <pre>
   * The expression is true when label value equals given value.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.EqualsMatchExpression label_equals = 5 [json_name = "labelEquals"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.EqualsMatchExpressionOrBuilder getLabelEqualsOrBuilder();

  /**
   * <pre>
   * The expression is true when label matches given regex.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchesMatchExpression label_matches = 6 [json_name = "labelMatches"];</code>
   * @return Whether the labelMatches field is set.
   */
  boolean hasLabelMatches();
  /**
   * <pre>
   * The expression is true when label matches given regex.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchesMatchExpression label_matches = 6 [json_name = "labelMatches"];</code>
   * @return The labelMatches.
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchesMatchExpression getLabelMatches();
  /**
   * <pre>
   * The expression is true when label matches given regex.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.MatchesMatchExpression label_matches = 6 [json_name = "labelMatches"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.MatchesMatchExpressionOrBuilder getLabelMatchesOrBuilder();

  public com.fluxninja.generated.aperture.policy.language.v1.MatchExpression.VariantCase getVariantCase();
}
