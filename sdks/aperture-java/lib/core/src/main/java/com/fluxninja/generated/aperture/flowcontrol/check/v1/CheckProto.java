// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

public final class CheckProto {
  private CheckProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CheckRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CheckRequest_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CheckRequest_LabelsEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CheckRequest_LabelsEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CheckResponse_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CheckResponse_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CheckResponse_TelemetryFlowLabelsEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CheckResponse_TelemetryFlowLabelsEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CachedValue_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CachedValue_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CacheUpsertRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CacheUpsertRequest_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CacheUpsertResponse_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CacheUpsertResponse_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CacheDeleteRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CacheDeleteRequest_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_ClassifierInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_ClassifierInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_check_v1_FluxMeterInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_check_v1_FluxMeterInfo_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n)aperture/flowcontrol/check/v1/check.pr" +
      "oto\022\035aperture.flowcontrol.check.v1\032\036goog" +
      "le/protobuf/duration.proto\032\037google/proto" +
      "buf/timestamp.proto\"\371\001\n\014CheckRequest\022#\n\r" +
      "control_point\030\001 \001(\tR\014controlPoint\022O\n\006lab" +
      "els\030\002 \003(\01327.aperture.flowcontrol.check.v" +
      "1.CheckRequest.LabelsEntryR\006labels\022\033\n\tra" +
      "mp_mode\030\003 \001(\010R\010rampMode\022\033\n\tcache_key\030\004 \001" +
      "(\tR\010cacheKey\0329\n\013LabelsEntry\022\020\n\003key\030\001 \001(\t" +
      "R\003key\022\024\n\005value\030\002 \001(\tR\005value:\0028\001\"\314\n\n\rChec" +
      "kResponse\0220\n\005start\030\001 \001(\0132\032.google.protob" +
      "uf.TimestampR\005start\022,\n\003end\030\002 \001(\0132\032.googl" +
      "e.protobuf.TimestampR\003end\022\032\n\010services\030\004 " +
      "\003(\tR\010services\022#\n\rcontrol_point\030\005 \001(\tR\014co" +
      "ntrolPoint\022&\n\017flow_label_keys\030\006 \003(\tR\rflo" +
      "wLabelKeys\022y\n\025telemetry_flow_labels\030\007 \003(" +
      "\0132E.aperture.flowcontrol.check.v1.CheckR" +
      "esponse.TelemetryFlowLabelsEntryR\023teleme" +
      "tryFlowLabels\022^\n\rdecision_type\030\010 \001(\01629.a" +
      "perture.flowcontrol.check.v1.CheckRespon" +
      "se.DecisionTypeR\014decisionType\022^\n\rreject_" +
      "reason\030\t \001(\01629.aperture.flowcontrol.chec" +
      "k.v1.CheckResponse.RejectReasonR\014rejectR" +
      "eason\022X\n\020classifier_infos\030\n \003(\0132-.apertu" +
      "re.flowcontrol.check.v1.ClassifierInfoR\017" +
      "classifierInfos\022V\n\020flux_meter_infos\030\013 \003(" +
      "\0132,.aperture.flowcontrol.check.v1.FluxMe" +
      "terInfoR\016fluxMeterInfos\022[\n\021limiter_decis" +
      "ions\030\014 \003(\0132..aperture.flowcontrol.check." +
      "v1.LimiterDecisionR\020limiterDecisions\0226\n\t" +
      "wait_time\030\r \001(\0132\031.google.protobuf.Durati" +
      "onR\010waitTime\022h\n\033denied_response_status_c" +
      "ode\030\016 \001(\0162).aperture.flowcontrol.check.v" +
      "1.StatusCodeR\030deniedResponseStatusCode\022M" +
      "\n\014cached_value\030\017 \001(\0132*.aperture.flowcont" +
      "rol.check.v1.CachedValueR\013cachedValue\032F\n" +
      "\030TelemetryFlowLabelsEntry\022\020\n\003key\030\001 \001(\tR\003" +
      "key\022\024\n\005value\030\002 \001(\tR\005value:\0028\001\"\246\001\n\014Reject" +
      "Reason\022\026\n\022REJECT_REASON_NONE\020\000\022\036\n\032REJECT" +
      "_REASON_RATE_LIMITED\020\001\022\033\n\027REJECT_REASON_" +
      "NO_TOKENS\020\002\022\035\n\031REJECT_REASON_NOT_SAMPLED" +
      "\020\003\022\"\n\036REJECT_REASON_NO_MATCHING_RAMP\020\004\"F" +
      "\n\014DecisionType\022\032\n\026DECISION_TYPE_ACCEPTED" +
      "\020\000\022\032\n\026DECISION_TYPE_REJECTED\020\001\"\353\001\n\013Cache" +
      "dValue\022\024\n\005value\030\001 \001(\014R\005value\022U\n\rlookup_r" +
      "esult\030\002 \001(\01620.aperture.flowcontrol.check" +
      ".v1.CacheLookupResultR\014lookupResult\022U\n\rr" +
      "esponse_code\030\003 \001(\01620.aperture.flowcontro" +
      "l.check.v1.CacheResponseCodeR\014responseCo" +
      "de\022\030\n\007message\030\004 \001(\tR\007message\"\216\001\n\022CacheUp" +
      "sertRequest\022#\n\rcontrol_point\030\001 \001(\tR\014cont" +
      "rolPoint\022\020\n\003key\030\002 \001(\tR\003key\022\024\n\005value\030\003 \001(" +
      "\014R\005value\022+\n\003ttl\030\004 \001(\0132\031.google.protobuf." +
      "DurationR\003ttl\"u\n\023CacheUpsertResponse\022D\n\004" +
      "code\030\001 \001(\01620.aperture.flowcontrol.check." +
      "v1.CacheResponseCodeR\004code\022\030\n\007message\030\002 " +
      "\001(\tR\007message\"K\n\022CacheDeleteRequest\022#\n\rco" +
      "ntrol_point\030\001 \001(\tR\014controlPoint\022\020\n\003key\030\002" +
      " \001(\tR\003key\"u\n\023CacheDeleteResponse\022D\n\004code" +
      "\030\001 \001(\01620.aperture.flowcontrol.check.v1.C" +
      "acheResponseCodeR\004code\022\030\n\007message\030\002 \001(\tR" +
      "\007message\"\355\002\n\016ClassifierInfo\022\037\n\013policy_na" +
      "me\030\001 \001(\tR\npolicyName\022\037\n\013policy_hash\030\002 \001(" +
      "\tR\npolicyHash\022)\n\020classifier_index\030\003 \001(\003R" +
      "\017classifierIndex\022I\n\005error\030\005 \001(\01623.apertu" +
      "re.flowcontrol.check.v1.ClassifierInfo.E" +
      "rrorR\005error\"\242\001\n\005Error\022\016\n\nERROR_NONE\020\000\022\025\n" +
      "\021ERROR_EVAL_FAILED\020\001\022\031\n\025ERROR_EMPTY_RESU" +
      "LTSET\020\002\022\035\n\031ERROR_AMBIGUOUS_RESULTSET\020\003\022\032" +
      "\n\026ERROR_MULTI_EXPRESSION\020\004\022\034\n\030ERROR_EXPR" +
      "ESSION_NOT_MAP\020\005\"\246\014\n\017LimiterDecision\022\037\n\013" +
      "policy_name\030\001 \001(\tR\npolicyName\022\037\n\013policy_" +
      "hash\030\002 \001(\tR\npolicyHash\022!\n\014component_id\030\003" +
      " \001(\tR\013componentId\022\030\n\007dropped\030\004 \001(\010R\007drop" +
      "ped\022T\n\006reason\030\005 \001(\0162<.aperture.flowcontr" +
      "ol.check.v1.LimiterDecision.LimiterReaso" +
      "nR\006reason\022h\n\033denied_response_status_code" +
      "\030\n \001(\0162).aperture.flowcontrol.check.v1.S" +
      "tatusCodeR\030deniedResponseStatusCode\0226\n\tw" +
      "ait_time\030\013 \001(\0132\031.google.protobuf.Duratio" +
      "nR\010waitTime\022l\n\021rate_limiter_info\030\024 \001(\0132>" +
      ".aperture.flowcontrol.check.v1.LimiterDe" +
      "cision.RateLimiterInfoH\000R\017rateLimiterInf" +
      "o\022n\n\023load_scheduler_info\030\025 \001(\0132<.apertur" +
      "e.flowcontrol.check.v1.LimiterDecision.S" +
      "chedulerInfoH\000R\021loadSchedulerInfo\022_\n\014sam" +
      "pler_info\030\026 \001(\0132:.aperture.flowcontrol.c" +
      "heck.v1.LimiterDecision.SamplerInfoH\000R\013s" +
      "amplerInfo\022u\n\024quota_scheduler_info\030\027 \001(\013" +
      "2A.aperture.flowcontrol.check.v1.Limiter" +
      "Decision.QuotaSchedulerInfoH\000R\022quotaSche" +
      "dulerInfo\032`\n\nTokensInfo\022\034\n\tremaining\030\001 \001" +
      "(\001R\tremaining\022\030\n\007current\030\002 \001(\001R\007current\022" +
      "\032\n\010consumed\030\003 \001(\001R\010consumed\032\203\001\n\017RateLimi" +
      "terInfo\022\024\n\005label\030\001 \001(\tR\005label\022Z\n\013tokens_" +
      "info\030\002 \001(\01329.aperture.flowcontrol.check." +
      "v1.LimiterDecision.TokensInfoR\ntokensInf" +
      "o\032\256\001\n\rSchedulerInfo\022%\n\016workload_index\030\001 " +
      "\001(\tR\rworkloadIndex\022Z\n\013tokens_info\030\002 \001(\0132" +
      "9.aperture.flowcontrol.check.v1.LimiterD" +
      "ecision.TokensInfoR\ntokensInfo\022\032\n\010priori" +
      "ty\030\003 \001(\001R\010priority\032#\n\013SamplerInfo\022\024\n\005lab" +
      "el\030\001 \001(\tR\005label\032\311\001\n\022QuotaSchedulerInfo\022\024" +
      "\n\005label\030\001 \001(\tR\005label\022%\n\016workload_index\030\002" +
      " \001(\tR\rworkloadIndex\022Z\n\013tokens_info\030\003 \001(\013" +
      "29.aperture.flowcontrol.check.v1.Limiter" +
      "Decision.TokensInfoR\ntokensInfo\022\032\n\010prior" +
      "ity\030\004 \001(\001R\010priority\"Q\n\rLimiterReason\022\036\n\032" +
      "LIMITER_REASON_UNSPECIFIED\020\000\022 \n\034LIMITER_" +
      "REASON_KEY_NOT_FOUND\020\001B\t\n\007details\"7\n\rFlu" +
      "xMeterInfo\022&\n\017flux_meter_name\030\001 \001(\tR\rflu" +
      "xMeterName*&\n\021CacheLookupResult\022\007\n\003HIT\020\000" +
      "\022\010\n\004MISS\020\001*+\n\021CacheResponseCode\022\013\n\007SUCCE" +
      "SS\020\000\022\t\n\005ERROR\020\001*\265\t\n\nStatusCode\022\t\n\005Empty\020" +
      "\000\022\014\n\010Continue\020d\022\007\n\002OK\020\310\001\022\014\n\007Created\020\311\001\022\r" +
      "\n\010Accepted\020\312\001\022 \n\033NonAuthoritativeInforma" +
      "tion\020\313\001\022\016\n\tNoContent\020\314\001\022\021\n\014ResetContent\020" +
      "\315\001\022\023\n\016PartialContent\020\316\001\022\020\n\013MultiStatus\020\317" +
      "\001\022\024\n\017AlreadyReported\020\320\001\022\013\n\006IMUsed\020\342\001\022\024\n\017" +
      "MultipleChoices\020\254\002\022\025\n\020MovedPermanently\020\255" +
      "\002\022\n\n\005Found\020\256\002\022\r\n\010SeeOther\020\257\002\022\020\n\013NotModif" +
      "ied\020\260\002\022\r\n\010UseProxy\020\261\002\022\026\n\021TemporaryRedire" +
      "ct\020\263\002\022\026\n\021PermanentRedirect\020\264\002\022\017\n\nBadRequ" +
      "est\020\220\003\022\021\n\014Unauthorized\020\221\003\022\024\n\017PaymentRequ" +
      "ired\020\222\003\022\016\n\tForbidden\020\223\003\022\r\n\010NotFound\020\224\003\022\025" +
      "\n\020MethodNotAllowed\020\225\003\022\022\n\rNotAcceptable\020\226" +
      "\003\022 \n\033ProxyAuthenticationRequired\020\227\003\022\023\n\016R" +
      "equestTimeout\020\230\003\022\r\n\010Conflict\020\231\003\022\t\n\004Gone\020" +
      "\232\003\022\023\n\016LengthRequired\020\233\003\022\027\n\022PreconditionF" +
      "ailed\020\234\003\022\024\n\017PayloadTooLarge\020\235\003\022\017\n\nURIToo" +
      "Long\020\236\003\022\031\n\024UnsupportedMediaType\020\237\003\022\030\n\023Ra" +
      "ngeNotSatisfiable\020\240\003\022\026\n\021ExpectationFaile" +
      "d\020\241\003\022\027\n\022MisdirectedRequest\020\245\003\022\030\n\023Unproce" +
      "ssableEntity\020\246\003\022\013\n\006Locked\020\247\003\022\025\n\020FailedDe" +
      "pendency\020\250\003\022\024\n\017UpgradeRequired\020\252\003\022\031\n\024Pre" +
      "conditionRequired\020\254\003\022\024\n\017TooManyRequests\020" +
      "\255\003\022 \n\033RequestHeaderFieldsTooLarge\020\257\003\022\030\n\023" +
      "InternalServerError\020\364\003\022\023\n\016NotImplemented" +
      "\020\365\003\022\017\n\nBadGateway\020\366\003\022\027\n\022ServiceUnavailab" +
      "le\020\367\003\022\023\n\016GatewayTimeout\020\370\003\022\034\n\027HTTPVersio" +
      "nNotSupported\020\371\003\022\032\n\025VariantAlsoNegotiate" +
      "s\020\372\003\022\030\n\023InsufficientStorage\020\373\003\022\021\n\014LoopDe" +
      "tected\020\374\003\022\020\n\013NotExtended\020\376\003\022\"\n\035NetworkAu" +
      "thenticationRequired\020\377\0032\352\002\n\022FlowControlS" +
      "ervice\022d\n\005Check\022+.aperture.flowcontrol.c" +
      "heck.v1.CheckRequest\032,.aperture.flowcont" +
      "rol.check.v1.CheckResponse\"\000\022v\n\013CacheUps" +
      "ert\0221.aperture.flowcontrol.check.v1.Cach" +
      "eUpsertRequest\0322.aperture.flowcontrol.ch" +
      "eck.v1.CacheUpsertResponse\"\000\022v\n\013CacheDel" +
      "ete\0221.aperture.flowcontrol.check.v1.Cach" +
      "eDeleteRequest\0322.aperture.flowcontrol.ch" +
      "eck.v1.CacheDeleteResponse\"\000B\263\002\n5com.flu" +
      "xninja.generated.aperture.flowcontrol.ch" +
      "eck.v1B\nCheckProtoP\001ZWgithub.com/fluxnin" +
      "ja/aperture/v2/api/gen/proto/go/aperture" +
      "/flowcontrol/check/v1;checkv1\242\002\003AFC\252\002\035Ap" +
      "erture.Flowcontrol.Check.V1\312\002\035Aperture\\F" +
      "lowcontrol\\Check\\V1\342\002)Aperture\\Flowcontr" +
      "ol\\Check\\V1\\GPBMetadata\352\002 Aperture::Flow" +
      "control::Check::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.DurationProto.getDescriptor(),
          com.google.protobuf.TimestampProto.getDescriptor(),
        });
    internal_static_aperture_flowcontrol_check_v1_CheckRequest_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_flowcontrol_check_v1_CheckRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CheckRequest_descriptor,
        new java.lang.String[] { "ControlPoint", "Labels", "RampMode", "CacheKey", });
    internal_static_aperture_flowcontrol_check_v1_CheckRequest_LabelsEntry_descriptor =
      internal_static_aperture_flowcontrol_check_v1_CheckRequest_descriptor.getNestedTypes().get(0);
    internal_static_aperture_flowcontrol_check_v1_CheckRequest_LabelsEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CheckRequest_LabelsEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_aperture_flowcontrol_check_v1_CheckResponse_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_flowcontrol_check_v1_CheckResponse_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CheckResponse_descriptor,
        new java.lang.String[] { "Start", "End", "Services", "ControlPoint", "FlowLabelKeys", "TelemetryFlowLabels", "DecisionType", "RejectReason", "ClassifierInfos", "FluxMeterInfos", "LimiterDecisions", "WaitTime", "DeniedResponseStatusCode", "CachedValue", });
    internal_static_aperture_flowcontrol_check_v1_CheckResponse_TelemetryFlowLabelsEntry_descriptor =
      internal_static_aperture_flowcontrol_check_v1_CheckResponse_descriptor.getNestedTypes().get(0);
    internal_static_aperture_flowcontrol_check_v1_CheckResponse_TelemetryFlowLabelsEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CheckResponse_TelemetryFlowLabelsEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_aperture_flowcontrol_check_v1_CachedValue_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_flowcontrol_check_v1_CachedValue_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CachedValue_descriptor,
        new java.lang.String[] { "Value", "LookupResult", "ResponseCode", "Message", });
    internal_static_aperture_flowcontrol_check_v1_CacheUpsertRequest_descriptor =
      getDescriptor().getMessageTypes().get(3);
    internal_static_aperture_flowcontrol_check_v1_CacheUpsertRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CacheUpsertRequest_descriptor,
        new java.lang.String[] { "ControlPoint", "Key", "Value", "Ttl", });
    internal_static_aperture_flowcontrol_check_v1_CacheUpsertResponse_descriptor =
      getDescriptor().getMessageTypes().get(4);
    internal_static_aperture_flowcontrol_check_v1_CacheUpsertResponse_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CacheUpsertResponse_descriptor,
        new java.lang.String[] { "Code", "Message", });
    internal_static_aperture_flowcontrol_check_v1_CacheDeleteRequest_descriptor =
      getDescriptor().getMessageTypes().get(5);
    internal_static_aperture_flowcontrol_check_v1_CacheDeleteRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CacheDeleteRequest_descriptor,
        new java.lang.String[] { "ControlPoint", "Key", });
    internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_descriptor =
      getDescriptor().getMessageTypes().get(6);
    internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_descriptor,
        new java.lang.String[] { "Code", "Message", });
    internal_static_aperture_flowcontrol_check_v1_ClassifierInfo_descriptor =
      getDescriptor().getMessageTypes().get(7);
    internal_static_aperture_flowcontrol_check_v1_ClassifierInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_ClassifierInfo_descriptor,
        new java.lang.String[] { "PolicyName", "PolicyHash", "ClassifierIndex", "Error", });
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor =
      getDescriptor().getMessageTypes().get(8);
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor,
        new java.lang.String[] { "PolicyName", "PolicyHash", "ComponentId", "Dropped", "Reason", "DeniedResponseStatusCode", "WaitTime", "RateLimiterInfo", "LoadSchedulerInfo", "SamplerInfo", "QuotaSchedulerInfo", "Details", });
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo_descriptor =
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor.getNestedTypes().get(0);
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo_descriptor,
        new java.lang.String[] { "Remaining", "Current", "Consumed", });
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo_descriptor =
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor.getNestedTypes().get(1);
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo_descriptor,
        new java.lang.String[] { "Label", "TokensInfo", });
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo_descriptor =
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor.getNestedTypes().get(2);
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo_descriptor,
        new java.lang.String[] { "WorkloadIndex", "TokensInfo", "Priority", });
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo_descriptor =
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor.getNestedTypes().get(3);
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo_descriptor,
        new java.lang.String[] { "Label", });
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo_descriptor =
      internal_static_aperture_flowcontrol_check_v1_LimiterDecision_descriptor.getNestedTypes().get(4);
    internal_static_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo_descriptor,
        new java.lang.String[] { "Label", "WorkloadIndex", "TokensInfo", "Priority", });
    internal_static_aperture_flowcontrol_check_v1_FluxMeterInfo_descriptor =
      getDescriptor().getMessageTypes().get(9);
    internal_static_aperture_flowcontrol_check_v1_FluxMeterInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_check_v1_FluxMeterInfo_descriptor,
        new java.lang.String[] { "FluxMeterName", });
    com.google.protobuf.DurationProto.getDescriptor();
    com.google.protobuf.TimestampProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
