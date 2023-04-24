# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: aperture/flowcontrol/check/v1/check.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n)aperture/flowcontrol/check/v1/check.proto\x12\x1d\x61perture.flowcontrol.check.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\xbf\x01\n\x0c\x43heckRequest\x12#\n\rcontrol_point\x18\x01 \x01(\tR\x0c\x63ontrolPoint\x12O\n\x06labels\x18\x02 \x03(\x0b\x32\x37.aperture.flowcontrol.check.v1.CheckRequest.LabelsEntryR\x06labels\x1a\x39\n\x0bLabelsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\"\xc4\x08\n\rCheckResponse\x12\x30\n\x05start\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\x05start\x12,\n\x03\x65nd\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\x03\x65nd\x12\x1a\n\x08services\x18\x04 \x03(\tR\x08services\x12#\n\rcontrol_point\x18\x05 \x01(\tR\x0c\x63ontrolPoint\x12&\n\x0f\x66low_label_keys\x18\x06 \x03(\tR\rflowLabelKeys\x12y\n\x15telemetry_flow_labels\x18\x07 \x03(\x0b\x32\x45.aperture.flowcontrol.check.v1.CheckResponse.TelemetryFlowLabelsEntryR\x13telemetryFlowLabels\x12^\n\rdecision_type\x18\x08 \x01(\x0e\x32\x39.aperture.flowcontrol.check.v1.CheckResponse.DecisionTypeR\x0c\x64\x65\x63isionType\x12^\n\rreject_reason\x18\t \x01(\x0e\x32\x39.aperture.flowcontrol.check.v1.CheckResponse.RejectReasonR\x0crejectReason\x12X\n\x10\x63lassifier_infos\x18\n \x03(\x0b\x32-.aperture.flowcontrol.check.v1.ClassifierInfoR\x0f\x63lassifierInfos\x12V\n\x10\x66lux_meter_infos\x18\x0b \x03(\x0b\x32,.aperture.flowcontrol.check.v1.FluxMeterInfoR\x0e\x66luxMeterInfos\x12[\n\x11limiter_decisions\x18\x0c \x03(\x0b\x32..aperture.flowcontrol.check.v1.LimiterDecisionR\x10limiterDecisions\x1a\x46\n\x18TelemetryFlowLabelsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\"\x8f\x01\n\x0cRejectReason\x12\x16\n\x12REJECT_REASON_NONE\x10\x00\x12\x1e\n\x1aREJECT_REASON_RATE_LIMITED\x10\x01\x12%\n!REJECT_REASON_CONCURRENCY_LIMITED\x10\x02\x12 \n\x1cREJECT_REASON_LOAD_REGULATED\x10\x03\"F\n\x0c\x44\x65\x63isionType\x12\x1a\n\x16\x44\x45\x43ISION_TYPE_ACCEPTED\x10\x00\x12\x1a\n\x16\x44\x45\x43ISION_TYPE_REJECTED\x10\x01\"\xed\x02\n\x0e\x43lassifierInfo\x12\x1f\n\x0bpolicy_name\x18\x01 \x01(\tR\npolicyName\x12\x1f\n\x0bpolicy_hash\x18\x02 \x01(\tR\npolicyHash\x12)\n\x10\x63lassifier_index\x18\x03 \x01(\x03R\x0f\x63lassifierIndex\x12I\n\x05\x65rror\x18\x05 \x01(\x0e\x32\x33.aperture.flowcontrol.check.v1.ClassifierInfo.ErrorR\x05\x65rror\"\xa2\x01\n\x05\x45rror\x12\x0e\n\nERROR_NONE\x10\x00\x12\x15\n\x11\x45RROR_EVAL_FAILED\x10\x01\x12\x19\n\x15\x45RROR_EMPTY_RESULTSET\x10\x02\x12\x1d\n\x19\x45RROR_AMBIGUOUS_RESULTSET\x10\x03\x12\x1a\n\x16\x45RROR_MULTI_EXPRESSION\x10\x04\x12\x1c\n\x18\x45RROR_EXPRESSION_NOT_MAP\x10\x05\"\xca\x07\n\x0fLimiterDecision\x12\x1f\n\x0bpolicy_name\x18\x01 \x01(\tR\npolicyName\x12\x1f\n\x0bpolicy_hash\x18\x02 \x01(\tR\npolicyHash\x12!\n\x0c\x63omponent_id\x18\x03 \x01(\tR\x0b\x63omponentId\x12\x18\n\x07\x64ropped\x18\x04 \x01(\x08R\x07\x64ropped\x12T\n\x06reason\x18\x05 \x01(\x0e\x32<.aperture.flowcontrol.check.v1.LimiterDecision.LimiterReasonR\x06reason\x12l\n\x11rate_limiter_info\x18\x06 \x01(\x0b\x32>.aperture.flowcontrol.check.v1.LimiterDecision.RateLimiterInfoH\x00R\x0frateLimiterInfo\x12\x81\x01\n\x18\x63oncurrency_limiter_info\x18\x07 \x01(\x0b\x32\x45.aperture.flowcontrol.check.v1.LimiterDecision.ConcurrencyLimiterInfoH\x00R\x16\x63oncurrencyLimiterInfo\x12r\n\x13load_regulator_info\x18\x08 \x01(\x0b\x32@.aperture.flowcontrol.check.v1.LimiterDecision.LoadRegulatorInfoH\x00R\x11loadRegulatorInfo\x1a\x88\x01\n\x0fRateLimiterInfo\x12\x1c\n\tremaining\x18\x01 \x01(\x03R\tremaining\x12\x18\n\x07\x63urrent\x18\x02 \x01(\x03R\x07\x63urrent\x12\x14\n\x05label\x18\x03 \x01(\tR\x05label\x12\'\n\x0ftokens_consumed\x18\x04 \x01(\x04R\x0etokensConsumed\x1ah\n\x16\x43oncurrencyLimiterInfo\x12%\n\x0eworkload_index\x18\x01 \x01(\tR\rworkloadIndex\x12\'\n\x0ftokens_consumed\x18\x02 \x01(\x04R\x0etokensConsumed\x1a)\n\x11LoadRegulatorInfo\x12\x14\n\x05label\x18\x01 \x01(\tR\x05label\"Q\n\rLimiterReason\x12\x1e\n\x1aLIMITER_REASON_UNSPECIFIED\x10\x00\x12 \n\x1cLIMITER_REASON_KEY_NOT_FOUND\x10\x01\x42\t\n\x07\x64\x65tails\"7\n\rFluxMeterInfo\x12&\n\x0f\x66lux_meter_name\x18\x01 \x01(\tR\rfluxMeterName2z\n\x12\x46lowControlService\x12\x64\n\x05\x43heck\x12+.aperture.flowcontrol.check.v1.CheckRequest\x1a,.aperture.flowcontrol.check.v1.CheckResponse\"\x00\x42\xb0\x02\n5com.fluxninja.generated.aperture.flowcontrol.check.v1B\nCheckProtoP\x01ZTgithub.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1;checkv1\xa2\x02\x03\x41\x46\x43\xaa\x02\x1d\x41perture.Flowcontrol.Check.V1\xca\x02\x1d\x41perture\\Flowcontrol\\Check\\V1\xe2\x02)Aperture\\Flowcontrol\\Check\\V1\\GPBMetadata\xea\x02 Aperture::Flowcontrol::Check::V1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'aperture.flowcontrol.check.v1.check_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n5com.fluxninja.generated.aperture.flowcontrol.check.v1B\nCheckProtoP\001ZTgithub.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1;checkv1\242\002\003AFC\252\002\035Aperture.Flowcontrol.Check.V1\312\002\035Aperture\\Flowcontrol\\Check\\V1\342\002)Aperture\\Flowcontrol\\Check\\V1\\GPBMetadata\352\002 Aperture::Flowcontrol::Check::V1'
  _CHECKREQUEST_LABELSENTRY._options = None
  _CHECKREQUEST_LABELSENTRY._serialized_options = b'8\001'
  _CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY._options = None
  _CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY._serialized_options = b'8\001'
  _globals['_CHECKREQUEST']._serialized_start=110
  _globals['_CHECKREQUEST']._serialized_end=301
  _globals['_CHECKREQUEST_LABELSENTRY']._serialized_start=244
  _globals['_CHECKREQUEST_LABELSENTRY']._serialized_end=301
  _globals['_CHECKRESPONSE']._serialized_start=304
  _globals['_CHECKRESPONSE']._serialized_end=1396
  _globals['_CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY']._serialized_start=1108
  _globals['_CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY']._serialized_end=1178
  _globals['_CHECKRESPONSE_REJECTREASON']._serialized_start=1181
  _globals['_CHECKRESPONSE_REJECTREASON']._serialized_end=1324
  _globals['_CHECKRESPONSE_DECISIONTYPE']._serialized_start=1326
  _globals['_CHECKRESPONSE_DECISIONTYPE']._serialized_end=1396
  _globals['_CLASSIFIERINFO']._serialized_start=1399
  _globals['_CLASSIFIERINFO']._serialized_end=1764
  _globals['_CLASSIFIERINFO_ERROR']._serialized_start=1602
  _globals['_CLASSIFIERINFO_ERROR']._serialized_end=1764
  _globals['_LIMITERDECISION']._serialized_start=1767
  _globals['_LIMITERDECISION']._serialized_end=2737
  _globals['_LIMITERDECISION_RATELIMITERINFO']._serialized_start=2358
  _globals['_LIMITERDECISION_RATELIMITERINFO']._serialized_end=2494
  _globals['_LIMITERDECISION_CONCURRENCYLIMITERINFO']._serialized_start=2496
  _globals['_LIMITERDECISION_CONCURRENCYLIMITERINFO']._serialized_end=2600
  _globals['_LIMITERDECISION_LOADREGULATORINFO']._serialized_start=2602
  _globals['_LIMITERDECISION_LOADREGULATORINFO']._serialized_end=2643
  _globals['_LIMITERDECISION_LIMITERREASON']._serialized_start=2645
  _globals['_LIMITERDECISION_LIMITERREASON']._serialized_end=2726
  _globals['_FLUXMETERINFO']._serialized_start=2739
  _globals['_FLUXMETERINFO']._serialized_end=2794
  _globals['_FLOWCONTROLSERVICE']._serialized_start=2796
  _globals['_FLOWCONTROLSERVICE']._serialized_end=2918
# @@protoc_insertion_point(module_scope)
