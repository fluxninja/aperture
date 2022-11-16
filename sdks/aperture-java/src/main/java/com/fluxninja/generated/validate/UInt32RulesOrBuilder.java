// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: validate/validate.proto

package com.fluxninja.generated.validate;

public interface UInt32RulesOrBuilder extends
    // @@protoc_insertion_point(interface_extends:validate.UInt32Rules)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Const specifies that this field must be exactly the specified value
   * </pre>
   *
   * <code>optional uint32 const = 1 [json_name = "const"];</code>
   * @return Whether the const field is set.
   */
  boolean hasConst();
  /**
   * <pre>
   * Const specifies that this field must be exactly the specified value
   * </pre>
   *
   * <code>optional uint32 const = 1 [json_name = "const"];</code>
   * @return The const.
   */
  int getConst();

  /**
   * <pre>
   * Lt specifies that this field must be less than the specified value,
   * exclusive
   * </pre>
   *
   * <code>optional uint32 lt = 2 [json_name = "lt"];</code>
   * @return Whether the lt field is set.
   */
  boolean hasLt();
  /**
   * <pre>
   * Lt specifies that this field must be less than the specified value,
   * exclusive
   * </pre>
   *
   * <code>optional uint32 lt = 2 [json_name = "lt"];</code>
   * @return The lt.
   */
  int getLt();

  /**
   * <pre>
   * Lte specifies that this field must be less than or equal to the
   * specified value, inclusive
   * </pre>
   *
   * <code>optional uint32 lte = 3 [json_name = "lte"];</code>
   * @return Whether the lte field is set.
   */
  boolean hasLte();
  /**
   * <pre>
   * Lte specifies that this field must be less than or equal to the
   * specified value, inclusive
   * </pre>
   *
   * <code>optional uint32 lte = 3 [json_name = "lte"];</code>
   * @return The lte.
   */
  int getLte();

  /**
   * <pre>
   * Gt specifies that this field must be greater than the specified value,
   * exclusive. If the value of Gt is larger than a specified Lt or Lte, the
   * range is reversed.
   * </pre>
   *
   * <code>optional uint32 gt = 4 [json_name = "gt"];</code>
   * @return Whether the gt field is set.
   */
  boolean hasGt();
  /**
   * <pre>
   * Gt specifies that this field must be greater than the specified value,
   * exclusive. If the value of Gt is larger than a specified Lt or Lte, the
   * range is reversed.
   * </pre>
   *
   * <code>optional uint32 gt = 4 [json_name = "gt"];</code>
   * @return The gt.
   */
  int getGt();

  /**
   * <pre>
   * Gte specifies that this field must be greater than or equal to the
   * specified value, inclusive. If the value of Gte is larger than a
   * specified Lt or Lte, the range is reversed.
   * </pre>
   *
   * <code>optional uint32 gte = 5 [json_name = "gte"];</code>
   * @return Whether the gte field is set.
   */
  boolean hasGte();
  /**
   * <pre>
   * Gte specifies that this field must be greater than or equal to the
   * specified value, inclusive. If the value of Gte is larger than a
   * specified Lt or Lte, the range is reversed.
   * </pre>
   *
   * <code>optional uint32 gte = 5 [json_name = "gte"];</code>
   * @return The gte.
   */
  int getGte();

  /**
   * <pre>
   * In specifies that this field must be equal to one of the specified
   * values
   * </pre>
   *
   * <code>repeated uint32 in = 6 [json_name = "in"];</code>
   * @return A list containing the in.
   */
  java.util.List<java.lang.Integer> getInList();
  /**
   * <pre>
   * In specifies that this field must be equal to one of the specified
   * values
   * </pre>
   *
   * <code>repeated uint32 in = 6 [json_name = "in"];</code>
   * @return The count of in.
   */
  int getInCount();
  /**
   * <pre>
   * In specifies that this field must be equal to one of the specified
   * values
   * </pre>
   *
   * <code>repeated uint32 in = 6 [json_name = "in"];</code>
   * @param index The index of the element to return.
   * @return The in at the given index.
   */
  int getIn(int index);

  /**
   * <pre>
   * NotIn specifies that this field cannot be equal to one of the specified
   * values
   * </pre>
   *
   * <code>repeated uint32 not_in = 7 [json_name = "notIn"];</code>
   * @return A list containing the notIn.
   */
  java.util.List<java.lang.Integer> getNotInList();
  /**
   * <pre>
   * NotIn specifies that this field cannot be equal to one of the specified
   * values
   * </pre>
   *
   * <code>repeated uint32 not_in = 7 [json_name = "notIn"];</code>
   * @return The count of notIn.
   */
  int getNotInCount();
  /**
   * <pre>
   * NotIn specifies that this field cannot be equal to one of the specified
   * values
   * </pre>
   *
   * <code>repeated uint32 not_in = 7 [json_name = "notIn"];</code>
   * @param index The index of the element to return.
   * @return The notIn at the given index.
   */
  int getNotIn(int index);

  /**
   * <pre>
   * IgnoreEmpty specifies that the validation rules of this field should be
   * evaluated only if the field is not empty
   * </pre>
   *
   * <code>optional bool ignore_empty = 8 [json_name = "ignoreEmpty"];</code>
   * @return Whether the ignoreEmpty field is set.
   */
  boolean hasIgnoreEmpty();
  /**
   * <pre>
   * IgnoreEmpty specifies that the validation rules of this field should be
   * evaluated only if the field is not empty
   * </pre>
   *
   * <code>optional bool ignore_empty = 8 [json_name = "ignoreEmpty"];</code>
   * @return The ignoreEmpty.
   */
  boolean getIgnoreEmpty();
}
