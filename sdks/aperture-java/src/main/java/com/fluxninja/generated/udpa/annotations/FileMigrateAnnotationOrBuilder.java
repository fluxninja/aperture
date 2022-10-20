// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: udpa/annotations/migrate.proto

package com.fluxninja.generated.udpa.annotations;

public interface FileMigrateAnnotationOrBuilder extends
    // @@protoc_insertion_point(interface_extends:udpa.annotations.FileMigrateAnnotation)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Move all types in the file to another package, this implies changing proto
   * file path.
   * </pre>
   *
   * <code>string move_to_package = 2 [json_name = "moveToPackage"];</code>
   * @return The moveToPackage.
   */
  java.lang.String getMoveToPackage();
  /**
   * <pre>
   * Move all types in the file to another package, this implies changing proto
   * file path.
   * </pre>
   *
   * <code>string move_to_package = 2 [json_name = "moveToPackage"];</code>
   * @return The bytes for moveToPackage.
   */
  com.google.protobuf.ByteString
      getMoveToPackageBytes();
}
