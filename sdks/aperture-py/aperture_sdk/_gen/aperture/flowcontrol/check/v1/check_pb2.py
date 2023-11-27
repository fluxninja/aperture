# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: aperture/flowcontrol/check/v1/check.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n)aperture/flowcontrol/check/v1/check.proto\x12\x1d\x61perture.flowcontrol.check.v1\x1a\x1egoogle/protobuf/duration.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xc1\x02\n\x0c\x43heckRequest\x12#\n\rcontrol_point\x18\x01 \x01(\tR\x0c\x63ontrolPoint\x12O\n\x06labels\x18\x02 \x03(\x0b\x32\x37.aperture.flowcontrol.check.v1.CheckRequest.LabelsEntryR\x06labels\x12\x1b\n\tramp_mode\x18\x03 \x01(\x08R\x08rampMode\x12\x63\n\x14\x63\x61\x63he_lookup_request\x18\x04 \x01(\x0b\x32\x31.aperture.flowcontrol.check.v1.CacheLookupRequestR\x12\x63\x61\x63heLookupRequest\x1a\x39\n\x0bLabelsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\"\xe5\n\n\rCheckResponse\x12\x30\n\x05start\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\x05start\x12,\n\x03\x65nd\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\x03\x65nd\x12\x1a\n\x08services\x18\x04 \x03(\tR\x08services\x12#\n\rcontrol_point\x18\x05 \x01(\tR\x0c\x63ontrolPoint\x12&\n\x0f\x66low_label_keys\x18\x06 \x03(\tR\rflowLabelKeys\x12y\n\x15telemetry_flow_labels\x18\x07 \x03(\x0b\x32\x45.aperture.flowcontrol.check.v1.CheckResponse.TelemetryFlowLabelsEntryR\x13telemetryFlowLabels\x12^\n\rdecision_type\x18\x08 \x01(\x0e\x32\x39.aperture.flowcontrol.check.v1.CheckResponse.DecisionTypeR\x0c\x64\x65\x63isionType\x12^\n\rreject_reason\x18\t \x01(\x0e\x32\x39.aperture.flowcontrol.check.v1.CheckResponse.RejectReasonR\x0crejectReason\x12X\n\x10\x63lassifier_infos\x18\n \x03(\x0b\x32-.aperture.flowcontrol.check.v1.ClassifierInfoR\x0f\x63lassifierInfos\x12V\n\x10\x66lux_meter_infos\x18\x0b \x03(\x0b\x32,.aperture.flowcontrol.check.v1.FluxMeterInfoR\x0e\x66luxMeterInfos\x12[\n\x11limiter_decisions\x18\x0c \x03(\x0b\x32..aperture.flowcontrol.check.v1.LimiterDecisionR\x10limiterDecisions\x12\x36\n\twait_time\x18\r \x01(\x0b\x32\x19.google.protobuf.DurationR\x08waitTime\x12h\n\x1b\x64\x65nied_response_status_code\x18\x0e \x01(\x0e\x32).aperture.flowcontrol.check.v1.StatusCodeR\x18\x64\x65niedResponseStatusCode\x12\x66\n\x15\x63\x61\x63he_lookup_response\x18\x0f \x01(\x0b\x32\x32.aperture.flowcontrol.check.v1.CacheLookupResponseR\x13\x63\x61\x63heLookupResponse\x1a\x46\n\x18TelemetryFlowLabelsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\"\xa6\x01\n\x0cRejectReason\x12\x16\n\x12REJECT_REASON_NONE\x10\x00\x12\x1e\n\x1aREJECT_REASON_RATE_LIMITED\x10\x01\x12\x1b\n\x17REJECT_REASON_NO_TOKENS\x10\x02\x12\x1d\n\x19REJECT_REASON_NOT_SAMPLED\x10\x03\x12\"\n\x1eREJECT_REASON_NO_MATCHING_RAMP\x10\x04\"F\n\x0c\x44\x65\x63isionType\x12\x1a\n\x16\x44\x45\x43ISION_TYPE_ACCEPTED\x10\x00\x12\x1a\n\x16\x44\x45\x43ISION_TYPE_REJECTED\x10\x01\"\x88\x02\n\x11KeyLookupResponse\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\x0cR\x05value\x12U\n\rlookup_status\x18\x03 \x01(\x0e\x32\x30.aperture.flowcontrol.check.v1.CacheLookupStatusR\x0clookupStatus\x12^\n\x10operation_status\x18\x04 \x01(\x0e\x32\x33.aperture.flowcontrol.check.v1.CacheOperationStatusR\x0foperationStatus\x12\x14\n\x05\x65rror\x18\x05 \x01(\tR\x05\x65rror\"\xfd\x02\n\x12\x43\x61\x63heUpsertRequest\x12#\n\rcontrol_point\x18\x01 \x01(\tR\x0c\x63ontrolPoint\x12W\n\x12result_cache_entry\x18\x02 \x01(\x0b\x32).aperture.flowcontrol.check.v1.CacheEntryR\x10resultCacheEntry\x12x\n\x13state_cache_entries\x18\x03 \x03(\x0b\x32H.aperture.flowcontrol.check.v1.CacheUpsertRequest.StateCacheEntriesEntryR\x11stateCacheEntries\x1ao\n\x16StateCacheEntriesEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12?\n\x05value\x18\x02 \x01(\x0b\x32).aperture.flowcontrol.check.v1.CacheEntryR\x05value:\x02\x38\x01\"a\n\nCacheEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\x0cR\x05value\x12+\n\x03ttl\x18\x03 \x01(\x0b\x32\x19.google.protobuf.DurationR\x03ttl\"\xf6\x02\n\x13\x43\x61\x63heUpsertResponse\x12\x64\n\x15result_cache_response\x18\x01 \x01(\x0b\x32\x30.aperture.flowcontrol.check.v1.KeyUpsertResponseR\x13resultCacheResponse\x12\x7f\n\x15state_cache_responses\x18\x02 \x03(\x0b\x32K.aperture.flowcontrol.check.v1.CacheUpsertResponse.StateCacheResponsesEntryR\x13stateCacheResponses\x1ax\n\x18StateCacheResponsesEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x46\n\x05value\x18\x02 \x01(\x0b\x32\x30.aperture.flowcontrol.check.v1.KeyUpsertResponseR\x05value:\x02\x38\x01\"\x9b\x01\n\x11KeyUpsertResponse\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12^\n\x10operation_status\x18\x02 \x01(\x0e\x32\x33.aperture.flowcontrol.check.v1.CacheOperationStatusR\x0foperationStatus\x12\x14\n\x05\x65rror\x18\x03 \x01(\tR\x05\x65rror\"\x8d\x01\n\x12\x43\x61\x63heDeleteRequest\x12#\n\rcontrol_point\x18\x01 \x01(\tR\x0c\x63ontrolPoint\x12(\n\x10result_cache_key\x18\x02 \x01(\tR\x0eresultCacheKey\x12(\n\x10state_cache_keys\x18\x03 \x03(\tR\x0estateCacheKeys\"\xf6\x02\n\x13\x43\x61\x63heDeleteResponse\x12\x64\n\x15result_cache_response\x18\x01 \x01(\x0b\x32\x30.aperture.flowcontrol.check.v1.KeyDeleteResponseR\x13resultCacheResponse\x12\x7f\n\x15state_cache_responses\x18\x02 \x03(\x0b\x32K.aperture.flowcontrol.check.v1.CacheDeleteResponse.StateCacheResponsesEntryR\x13stateCacheResponses\x1ax\n\x18StateCacheResponsesEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x46\n\x05value\x18\x02 \x01(\x0b\x32\x30.aperture.flowcontrol.check.v1.KeyDeleteResponseR\x05value:\x02\x38\x01\"\x9b\x01\n\x11KeyDeleteResponse\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12^\n\x10operation_status\x18\x02 \x01(\x0e\x32\x33.aperture.flowcontrol.check.v1.CacheOperationStatusR\x0foperationStatus\x12\x14\n\x05\x65rror\x18\x03 \x01(\tR\x05\x65rror\"\x8d\x01\n\x12\x43\x61\x63heLookupRequest\x12#\n\rcontrol_point\x18\x01 \x01(\tR\x0c\x63ontrolPoint\x12(\n\x10result_cache_key\x18\x02 \x01(\tR\x0eresultCacheKey\x12(\n\x10state_cache_keys\x18\x03 \x03(\tR\x0estateCacheKeys\"\xf6\x02\n\x13\x43\x61\x63heLookupResponse\x12\x64\n\x15result_cache_response\x18\x01 \x01(\x0b\x32\x30.aperture.flowcontrol.check.v1.KeyLookupResponseR\x13resultCacheResponse\x12\x7f\n\x15state_cache_responses\x18\x02 \x03(\x0b\x32K.aperture.flowcontrol.check.v1.CacheLookupResponse.StateCacheResponsesEntryR\x13stateCacheResponses\x1ax\n\x18StateCacheResponsesEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x46\n\x05value\x18\x02 \x01(\x0b\x32\x30.aperture.flowcontrol.check.v1.KeyLookupResponseR\x05value:\x02\x38\x01\"\xed\x02\n\x0e\x43lassifierInfo\x12\x1f\n\x0bpolicy_name\x18\x01 \x01(\tR\npolicyName\x12\x1f\n\x0bpolicy_hash\x18\x02 \x01(\tR\npolicyHash\x12)\n\x10\x63lassifier_index\x18\x03 \x01(\x03R\x0f\x63lassifierIndex\x12I\n\x05\x65rror\x18\x05 \x01(\x0e\x32\x33.aperture.flowcontrol.check.v1.ClassifierInfo.ErrorR\x05\x65rror\"\xa2\x01\n\x05\x45rror\x12\x0e\n\nERROR_NONE\x10\x00\x12\x15\n\x11\x45RROR_EVAL_FAILED\x10\x01\x12\x19\n\x15\x45RROR_EMPTY_RESULTSET\x10\x02\x12\x1d\n\x19\x45RROR_AMBIGUOUS_RESULTSET\x10\x03\x12\x1a\n\x16\x45RROR_MULTI_EXPRESSION\x10\x04\x12\x1c\n\x18\x45RROR_EXPRESSION_NOT_MAP\x10\x05\"\xa6\x0c\n\x0fLimiterDecision\x12\x1f\n\x0bpolicy_name\x18\x01 \x01(\tR\npolicyName\x12\x1f\n\x0bpolicy_hash\x18\x02 \x01(\tR\npolicyHash\x12!\n\x0c\x63omponent_id\x18\x03 \x01(\tR\x0b\x63omponentId\x12\x18\n\x07\x64ropped\x18\x04 \x01(\x08R\x07\x64ropped\x12T\n\x06reason\x18\x05 \x01(\x0e\x32<.aperture.flowcontrol.check.v1.LimiterDecision.LimiterReasonR\x06reason\x12h\n\x1b\x64\x65nied_response_status_code\x18\n \x01(\x0e\x32).aperture.flowcontrol.check.v1.StatusCodeR\x18\x64\x65niedResponseStatusCode\x12\x36\n\twait_time\x18\x0b \x01(\x0b\x32\x19.google.protobuf.DurationR\x08waitTime\x12l\n\x11rate_limiter_info\x18\x14 \x01(\x0b\x32>.aperture.flowcontrol.check.v1.LimiterDecision.RateLimiterInfoH\x00R\x0frateLimiterInfo\x12n\n\x13load_scheduler_info\x18\x15 \x01(\x0b\x32<.aperture.flowcontrol.check.v1.LimiterDecision.SchedulerInfoH\x00R\x11loadSchedulerInfo\x12_\n\x0csampler_info\x18\x16 \x01(\x0b\x32:.aperture.flowcontrol.check.v1.LimiterDecision.SamplerInfoH\x00R\x0bsamplerInfo\x12u\n\x14quota_scheduler_info\x18\x17 \x01(\x0b\x32\x41.aperture.flowcontrol.check.v1.LimiterDecision.QuotaSchedulerInfoH\x00R\x12quotaSchedulerInfo\x1a`\n\nTokensInfo\x12\x1c\n\tremaining\x18\x01 \x01(\x01R\tremaining\x12\x18\n\x07\x63urrent\x18\x02 \x01(\x01R\x07\x63urrent\x12\x1a\n\x08\x63onsumed\x18\x03 \x01(\x01R\x08\x63onsumed\x1a\x83\x01\n\x0fRateLimiterInfo\x12\x14\n\x05label\x18\x01 \x01(\tR\x05label\x12Z\n\x0btokens_info\x18\x02 \x01(\x0b\x32\x39.aperture.flowcontrol.check.v1.LimiterDecision.TokensInfoR\ntokensInfo\x1a\xae\x01\n\rSchedulerInfo\x12%\n\x0eworkload_index\x18\x01 \x01(\tR\rworkloadIndex\x12Z\n\x0btokens_info\x18\x02 \x01(\x0b\x32\x39.aperture.flowcontrol.check.v1.LimiterDecision.TokensInfoR\ntokensInfo\x12\x1a\n\x08priority\x18\x03 \x01(\x01R\x08priority\x1a#\n\x0bSamplerInfo\x12\x14\n\x05label\x18\x01 \x01(\tR\x05label\x1a\xc9\x01\n\x12QuotaSchedulerInfo\x12\x14\n\x05label\x18\x01 \x01(\tR\x05label\x12%\n\x0eworkload_index\x18\x02 \x01(\tR\rworkloadIndex\x12Z\n\x0btokens_info\x18\x03 \x01(\x0b\x32\x39.aperture.flowcontrol.check.v1.LimiterDecision.TokensInfoR\ntokensInfo\x12\x1a\n\x08priority\x18\x04 \x01(\x01R\x08priority\"Q\n\rLimiterReason\x12\x1e\n\x1aLIMITER_REASON_UNSPECIFIED\x10\x00\x12 \n\x1cLIMITER_REASON_KEY_NOT_FOUND\x10\x01\x42\t\n\x07\x64\x65tails\"7\n\rFluxMeterInfo\x12&\n\x0f\x66lux_meter_name\x18\x01 \x01(\tR\rfluxMeterName*&\n\x11\x43\x61\x63heLookupStatus\x12\x07\n\x03HIT\x10\x00\x12\x08\n\x04MISS\x10\x01*.\n\x14\x43\x61\x63heOperationStatus\x12\x0b\n\x07SUCCESS\x10\x00\x12\t\n\x05\x45RROR\x10\x01*\xb5\t\n\nStatusCode\x12\t\n\x05\x45mpty\x10\x00\x12\x0c\n\x08\x43ontinue\x10\x64\x12\x07\n\x02OK\x10\xc8\x01\x12\x0c\n\x07\x43reated\x10\xc9\x01\x12\r\n\x08\x41\x63\x63\x65pted\x10\xca\x01\x12 \n\x1bNonAuthoritativeInformation\x10\xcb\x01\x12\x0e\n\tNoContent\x10\xcc\x01\x12\x11\n\x0cResetContent\x10\xcd\x01\x12\x13\n\x0ePartialContent\x10\xce\x01\x12\x10\n\x0bMultiStatus\x10\xcf\x01\x12\x14\n\x0f\x41lreadyReported\x10\xd0\x01\x12\x0b\n\x06IMUsed\x10\xe2\x01\x12\x14\n\x0fMultipleChoices\x10\xac\x02\x12\x15\n\x10MovedPermanently\x10\xad\x02\x12\n\n\x05\x46ound\x10\xae\x02\x12\r\n\x08SeeOther\x10\xaf\x02\x12\x10\n\x0bNotModified\x10\xb0\x02\x12\r\n\x08UseProxy\x10\xb1\x02\x12\x16\n\x11TemporaryRedirect\x10\xb3\x02\x12\x16\n\x11PermanentRedirect\x10\xb4\x02\x12\x0f\n\nBadRequest\x10\x90\x03\x12\x11\n\x0cUnauthorized\x10\x91\x03\x12\x14\n\x0fPaymentRequired\x10\x92\x03\x12\x0e\n\tForbidden\x10\x93\x03\x12\r\n\x08NotFound\x10\x94\x03\x12\x15\n\x10MethodNotAllowed\x10\x95\x03\x12\x12\n\rNotAcceptable\x10\x96\x03\x12 \n\x1bProxyAuthenticationRequired\x10\x97\x03\x12\x13\n\x0eRequestTimeout\x10\x98\x03\x12\r\n\x08\x43onflict\x10\x99\x03\x12\t\n\x04Gone\x10\x9a\x03\x12\x13\n\x0eLengthRequired\x10\x9b\x03\x12\x17\n\x12PreconditionFailed\x10\x9c\x03\x12\x14\n\x0fPayloadTooLarge\x10\x9d\x03\x12\x0f\n\nURITooLong\x10\x9e\x03\x12\x19\n\x14UnsupportedMediaType\x10\x9f\x03\x12\x18\n\x13RangeNotSatisfiable\x10\xa0\x03\x12\x16\n\x11\x45xpectationFailed\x10\xa1\x03\x12\x17\n\x12MisdirectedRequest\x10\xa5\x03\x12\x18\n\x13UnprocessableEntity\x10\xa6\x03\x12\x0b\n\x06Locked\x10\xa7\x03\x12\x15\n\x10\x46\x61iledDependency\x10\xa8\x03\x12\x14\n\x0fUpgradeRequired\x10\xaa\x03\x12\x19\n\x14PreconditionRequired\x10\xac\x03\x12\x14\n\x0fTooManyRequests\x10\xad\x03\x12 \n\x1bRequestHeaderFieldsTooLarge\x10\xaf\x03\x12\x18\n\x13InternalServerError\x10\xf4\x03\x12\x13\n\x0eNotImplemented\x10\xf5\x03\x12\x0f\n\nBadGateway\x10\xf6\x03\x12\x17\n\x12ServiceUnavailable\x10\xf7\x03\x12\x13\n\x0eGatewayTimeout\x10\xf8\x03\x12\x1c\n\x17HTTPVersionNotSupported\x10\xf9\x03\x12\x1a\n\x15VariantAlsoNegotiates\x10\xfa\x03\x12\x18\n\x13InsufficientStorage\x10\xfb\x03\x12\x11\n\x0cLoopDetected\x10\xfc\x03\x12\x10\n\x0bNotExtended\x10\xfe\x03\x12\"\n\x1dNetworkAuthenticationRequired\x10\xff\x03\x32\xe2\x03\n\x12\x46lowControlService\x12\x64\n\x05\x43heck\x12+.aperture.flowcontrol.check.v1.CheckRequest\x1a,.aperture.flowcontrol.check.v1.CheckResponse\"\x00\x12v\n\x0b\x43\x61\x63heUpsert\x12\x31.aperture.flowcontrol.check.v1.CacheUpsertRequest\x1a\x32.aperture.flowcontrol.check.v1.CacheUpsertResponse\"\x00\x12v\n\x0b\x43\x61\x63heDelete\x12\x31.aperture.flowcontrol.check.v1.CacheDeleteRequest\x1a\x32.aperture.flowcontrol.check.v1.CacheDeleteResponse\"\x00\x12v\n\x0b\x43\x61\x63heLookup\x12\x31.aperture.flowcontrol.check.v1.CacheLookupRequest\x1a\x32.aperture.flowcontrol.check.v1.CacheLookupResponse\"\x00\x42\xb3\x02\n5com.fluxninja.generated.aperture.flowcontrol.check.v1B\nCheckProtoP\x01ZWgithub.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1;checkv1\xa2\x02\x03\x41\x46\x43\xaa\x02\x1d\x41perture.Flowcontrol.Check.V1\xca\x02\x1d\x41perture\\Flowcontrol\\Check\\V1\xe2\x02)Aperture\\Flowcontrol\\Check\\V1\\GPBMetadata\xea\x02 Aperture::Flowcontrol::Check::V1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'aperture.flowcontrol.check.v1.check_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n5com.fluxninja.generated.aperture.flowcontrol.check.v1B\nCheckProtoP\001ZWgithub.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1;checkv1\242\002\003AFC\252\002\035Aperture.Flowcontrol.Check.V1\312\002\035Aperture\\Flowcontrol\\Check\\V1\342\002)Aperture\\Flowcontrol\\Check\\V1\\GPBMetadata\352\002 Aperture::Flowcontrol::Check::V1'
  _globals['_CHECKREQUEST_LABELSENTRY']._options = None
  _globals['_CHECKREQUEST_LABELSENTRY']._serialized_options = b'8\001'
  _globals['_CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY']._options = None
  _globals['_CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY']._serialized_options = b'8\001'
  _globals['_CACHEUPSERTREQUEST_STATECACHEENTRIESENTRY']._options = None
  _globals['_CACHEUPSERTREQUEST_STATECACHEENTRIESENTRY']._serialized_options = b'8\001'
  _globals['_CACHEUPSERTRESPONSE_STATECACHERESPONSESENTRY']._options = None
  _globals['_CACHEUPSERTRESPONSE_STATECACHERESPONSESENTRY']._serialized_options = b'8\001'
  _globals['_CACHEDELETERESPONSE_STATECACHERESPONSESENTRY']._options = None
  _globals['_CACHEDELETERESPONSE_STATECACHERESPONSESENTRY']._serialized_options = b'8\001'
  _globals['_CACHELOOKUPRESPONSE_STATECACHERESPONSESENTRY']._options = None
  _globals['_CACHELOOKUPRESPONSE_STATECACHERESPONSESENTRY']._serialized_options = b'8\001'
  _globals['_CACHELOOKUPSTATUS']._serialized_start=6336
  _globals['_CACHELOOKUPSTATUS']._serialized_end=6374
  _globals['_CACHEOPERATIONSTATUS']._serialized_start=6376
  _globals['_CACHEOPERATIONSTATUS']._serialized_end=6422
  _globals['_STATUSCODE']._serialized_start=6425
  _globals['_STATUSCODE']._serialized_end=7630
  _globals['_CHECKREQUEST']._serialized_start=142
  _globals['_CHECKREQUEST']._serialized_end=463
  _globals['_CHECKREQUEST_LABELSENTRY']._serialized_start=406
  _globals['_CHECKREQUEST_LABELSENTRY']._serialized_end=463
  _globals['_CHECKRESPONSE']._serialized_start=466
  _globals['_CHECKRESPONSE']._serialized_end=1847
  _globals['_CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY']._serialized_start=1536
  _globals['_CHECKRESPONSE_TELEMETRYFLOWLABELSENTRY']._serialized_end=1606
  _globals['_CHECKRESPONSE_REJECTREASON']._serialized_start=1609
  _globals['_CHECKRESPONSE_REJECTREASON']._serialized_end=1775
  _globals['_CHECKRESPONSE_DECISIONTYPE']._serialized_start=1777
  _globals['_CHECKRESPONSE_DECISIONTYPE']._serialized_end=1847
  _globals['_KEYLOOKUPRESPONSE']._serialized_start=1850
  _globals['_KEYLOOKUPRESPONSE']._serialized_end=2114
  _globals['_CACHEUPSERTREQUEST']._serialized_start=2117
  _globals['_CACHEUPSERTREQUEST']._serialized_end=2498
  _globals['_CACHEUPSERTREQUEST_STATECACHEENTRIESENTRY']._serialized_start=2387
  _globals['_CACHEUPSERTREQUEST_STATECACHEENTRIESENTRY']._serialized_end=2498
  _globals['_CACHEENTRY']._serialized_start=2500
  _globals['_CACHEENTRY']._serialized_end=2597
  _globals['_CACHEUPSERTRESPONSE']._serialized_start=2600
  _globals['_CACHEUPSERTRESPONSE']._serialized_end=2974
  _globals['_CACHEUPSERTRESPONSE_STATECACHERESPONSESENTRY']._serialized_start=2854
  _globals['_CACHEUPSERTRESPONSE_STATECACHERESPONSESENTRY']._serialized_end=2974
  _globals['_KEYUPSERTRESPONSE']._serialized_start=2977
  _globals['_KEYUPSERTRESPONSE']._serialized_end=3132
  _globals['_CACHEDELETEREQUEST']._serialized_start=3135
  _globals['_CACHEDELETEREQUEST']._serialized_end=3276
  _globals['_CACHEDELETERESPONSE']._serialized_start=3279
  _globals['_CACHEDELETERESPONSE']._serialized_end=3653
  _globals['_CACHEDELETERESPONSE_STATECACHERESPONSESENTRY']._serialized_start=3533
  _globals['_CACHEDELETERESPONSE_STATECACHERESPONSESENTRY']._serialized_end=3653
  _globals['_KEYDELETERESPONSE']._serialized_start=3656
  _globals['_KEYDELETERESPONSE']._serialized_end=3811
  _globals['_CACHELOOKUPREQUEST']._serialized_start=3814
  _globals['_CACHELOOKUPREQUEST']._serialized_end=3955
  _globals['_CACHELOOKUPRESPONSE']._serialized_start=3958
  _globals['_CACHELOOKUPRESPONSE']._serialized_end=4332
  _globals['_CACHELOOKUPRESPONSE_STATECACHERESPONSESENTRY']._serialized_start=4212
  _globals['_CACHELOOKUPRESPONSE_STATECACHERESPONSESENTRY']._serialized_end=4332
  _globals['_CLASSIFIERINFO']._serialized_start=4335
  _globals['_CLASSIFIERINFO']._serialized_end=4700
  _globals['_CLASSIFIERINFO_ERROR']._serialized_start=4538
  _globals['_CLASSIFIERINFO_ERROR']._serialized_end=4700
  _globals['_LIMITERDECISION']._serialized_start=4703
  _globals['_LIMITERDECISION']._serialized_end=6277
  _globals['_LIMITERDECISION_TOKENSINFO']._serialized_start=5535
  _globals['_LIMITERDECISION_TOKENSINFO']._serialized_end=5631
  _globals['_LIMITERDECISION_RATELIMITERINFO']._serialized_start=5634
  _globals['_LIMITERDECISION_RATELIMITERINFO']._serialized_end=5765
  _globals['_LIMITERDECISION_SCHEDULERINFO']._serialized_start=5768
  _globals['_LIMITERDECISION_SCHEDULERINFO']._serialized_end=5942
  _globals['_LIMITERDECISION_SAMPLERINFO']._serialized_start=5944
  _globals['_LIMITERDECISION_SAMPLERINFO']._serialized_end=5979
  _globals['_LIMITERDECISION_QUOTASCHEDULERINFO']._serialized_start=5982
  _globals['_LIMITERDECISION_QUOTASCHEDULERINFO']._serialized_end=6183
  _globals['_LIMITERDECISION_LIMITERREASON']._serialized_start=6185
  _globals['_LIMITERDECISION_LIMITERREASON']._serialized_end=6266
  _globals['_FLUXMETERINFO']._serialized_start=6279
  _globals['_FLUXMETERINFO']._serialized_end=6334
  _globals['_FLOWCONTROLSERVICE']._serialized_start=7633
  _globals['_FLOWCONTROLSERVICE']._serialized_end=8115
# @@protoc_insertion_point(module_scope)
