// package: aperture.flowcontrol.check.v1
// file: api/aperture/flowcontrol/check/v1/check.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class CheckRequest extends jspb.Message { 
    getControlPoint(): string;
    setControlPoint(value: string): CheckRequest;

    getLabelsMap(): jspb.Map<string, string>;
    clearLabelsMap(): void;
    getRampMode(): boolean;
    setRampMode(value: boolean): CheckRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CheckRequest.AsObject;
    static toObject(includeInstance: boolean, msg: CheckRequest): CheckRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CheckRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CheckRequest;
    static deserializeBinaryFromReader(message: CheckRequest, reader: jspb.BinaryReader): CheckRequest;
}

export namespace CheckRequest {
    export type AsObject = {
        controlPoint: string,

        labelsMap: Array<[string, string]>,
        rampMode: boolean,
    }
}

export class CheckResponse extends jspb.Message { 

    hasStart(): boolean;
    clearStart(): void;
    getStart(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setStart(value?: google_protobuf_timestamp_pb.Timestamp): CheckResponse;

    hasEnd(): boolean;
    clearEnd(): void;
    getEnd(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setEnd(value?: google_protobuf_timestamp_pb.Timestamp): CheckResponse;
    clearServicesList(): void;
    getServicesList(): Array<string>;
    setServicesList(value: Array<string>): CheckResponse;
    addServices(value: string, index?: number): string;
    getControlPoint(): string;
    setControlPoint(value: string): CheckResponse;
    clearFlowLabelKeysList(): void;
    getFlowLabelKeysList(): Array<string>;
    setFlowLabelKeysList(value: Array<string>): CheckResponse;
    addFlowLabelKeys(value: string, index?: number): string;

    getTelemetryFlowLabelsMap(): jspb.Map<string, string>;
    clearTelemetryFlowLabelsMap(): void;
    getDecisionType(): CheckResponse.DecisionType;
    setDecisionType(value: CheckResponse.DecisionType): CheckResponse;
    getRejectReason(): CheckResponse.RejectReason;
    setRejectReason(value: CheckResponse.RejectReason): CheckResponse;
    clearClassifierInfosList(): void;
    getClassifierInfosList(): Array<ClassifierInfo>;
    setClassifierInfosList(value: Array<ClassifierInfo>): CheckResponse;
    addClassifierInfos(value?: ClassifierInfo, index?: number): ClassifierInfo;
    clearFluxMeterInfosList(): void;
    getFluxMeterInfosList(): Array<FluxMeterInfo>;
    setFluxMeterInfosList(value: Array<FluxMeterInfo>): CheckResponse;
    addFluxMeterInfos(value?: FluxMeterInfo, index?: number): FluxMeterInfo;
    clearLimiterDecisionsList(): void;
    getLimiterDecisionsList(): Array<LimiterDecision>;
    setLimiterDecisionsList(value: Array<LimiterDecision>): CheckResponse;
    addLimiterDecisions(value?: LimiterDecision, index?: number): LimiterDecision;

    hasWaitTime(): boolean;
    clearWaitTime(): void;
    getWaitTime(): google_protobuf_duration_pb.Duration | undefined;
    setWaitTime(value?: google_protobuf_duration_pb.Duration): CheckResponse;
    getDeniedResponseStatusCode(): StatusCode;
    setDeniedResponseStatusCode(value: StatusCode): CheckResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CheckResponse.AsObject;
    static toObject(includeInstance: boolean, msg: CheckResponse): CheckResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CheckResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CheckResponse;
    static deserializeBinaryFromReader(message: CheckResponse, reader: jspb.BinaryReader): CheckResponse;
}

export namespace CheckResponse {
    export type AsObject = {
        start?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        end?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        servicesList: Array<string>,
        controlPoint: string,
        flowLabelKeysList: Array<string>,

        telemetryFlowLabelsMap: Array<[string, string]>,
        decisionType: CheckResponse.DecisionType,
        rejectReason: CheckResponse.RejectReason,
        classifierInfosList: Array<ClassifierInfo.AsObject>,
        fluxMeterInfosList: Array<FluxMeterInfo.AsObject>,
        limiterDecisionsList: Array<LimiterDecision.AsObject>,
        waitTime?: google_protobuf_duration_pb.Duration.AsObject,
        deniedResponseStatusCode: StatusCode,
    }

    export enum RejectReason {
    REJECT_REASON_NONE = 0,
    REJECT_REASON_RATE_LIMITED = 1,
    REJECT_REASON_NO_TOKENS = 2,
    REJECT_REASON_NOT_SAMPLED = 3,
    REJECT_REASON_NO_MATCHING_RAMP = 4,
    }

    export enum DecisionType {
    DECISION_TYPE_ACCEPTED = 0,
    DECISION_TYPE_REJECTED = 1,
    }

}

export class ClassifierInfo extends jspb.Message { 
    getPolicyName(): string;
    setPolicyName(value: string): ClassifierInfo;
    getPolicyHash(): string;
    setPolicyHash(value: string): ClassifierInfo;
    getClassifierIndex(): number;
    setClassifierIndex(value: number): ClassifierInfo;
    getError(): ClassifierInfo.Error;
    setError(value: ClassifierInfo.Error): ClassifierInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ClassifierInfo.AsObject;
    static toObject(includeInstance: boolean, msg: ClassifierInfo): ClassifierInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ClassifierInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ClassifierInfo;
    static deserializeBinaryFromReader(message: ClassifierInfo, reader: jspb.BinaryReader): ClassifierInfo;
}

export namespace ClassifierInfo {
    export type AsObject = {
        policyName: string,
        policyHash: string,
        classifierIndex: number,
        error: ClassifierInfo.Error,
    }

    export enum Error {
    ERROR_NONE = 0,
    ERROR_EVAL_FAILED = 1,
    ERROR_EMPTY_RESULTSET = 2,
    ERROR_AMBIGUOUS_RESULTSET = 3,
    ERROR_MULTI_EXPRESSION = 4,
    ERROR_EXPRESSION_NOT_MAP = 5,
    }

}

export class LimiterDecision extends jspb.Message { 
    getPolicyName(): string;
    setPolicyName(value: string): LimiterDecision;
    getPolicyHash(): string;
    setPolicyHash(value: string): LimiterDecision;
    getComponentId(): string;
    setComponentId(value: string): LimiterDecision;
    getDropped(): boolean;
    setDropped(value: boolean): LimiterDecision;
    getReason(): LimiterDecision.LimiterReason;
    setReason(value: LimiterDecision.LimiterReason): LimiterDecision;
    getDeniedResponseStatusCode(): StatusCode;
    setDeniedResponseStatusCode(value: StatusCode): LimiterDecision;

    hasWaitTime(): boolean;
    clearWaitTime(): void;
    getWaitTime(): google_protobuf_duration_pb.Duration | undefined;
    setWaitTime(value?: google_protobuf_duration_pb.Duration): LimiterDecision;

    hasRateLimiterInfo(): boolean;
    clearRateLimiterInfo(): void;
    getRateLimiterInfo(): LimiterDecision.RateLimiterInfo | undefined;
    setRateLimiterInfo(value?: LimiterDecision.RateLimiterInfo): LimiterDecision;

    hasLoadSchedulerInfo(): boolean;
    clearLoadSchedulerInfo(): void;
    getLoadSchedulerInfo(): LimiterDecision.SchedulerInfo | undefined;
    setLoadSchedulerInfo(value?: LimiterDecision.SchedulerInfo): LimiterDecision;

    hasSamplerInfo(): boolean;
    clearSamplerInfo(): void;
    getSamplerInfo(): LimiterDecision.SamplerInfo | undefined;
    setSamplerInfo(value?: LimiterDecision.SamplerInfo): LimiterDecision;

    hasQuotaSchedulerInfo(): boolean;
    clearQuotaSchedulerInfo(): void;
    getQuotaSchedulerInfo(): LimiterDecision.QuotaSchedulerInfo | undefined;
    setQuotaSchedulerInfo(value?: LimiterDecision.QuotaSchedulerInfo): LimiterDecision;

    getDetailsCase(): LimiterDecision.DetailsCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LimiterDecision.AsObject;
    static toObject(includeInstance: boolean, msg: LimiterDecision): LimiterDecision.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LimiterDecision, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LimiterDecision;
    static deserializeBinaryFromReader(message: LimiterDecision, reader: jspb.BinaryReader): LimiterDecision;
}

export namespace LimiterDecision {
    export type AsObject = {
        policyName: string,
        policyHash: string,
        componentId: string,
        dropped: boolean,
        reason: LimiterDecision.LimiterReason,
        deniedResponseStatusCode: StatusCode,
        waitTime?: google_protobuf_duration_pb.Duration.AsObject,
        rateLimiterInfo?: LimiterDecision.RateLimiterInfo.AsObject,
        loadSchedulerInfo?: LimiterDecision.SchedulerInfo.AsObject,
        samplerInfo?: LimiterDecision.SamplerInfo.AsObject,
        quotaSchedulerInfo?: LimiterDecision.QuotaSchedulerInfo.AsObject,
    }


    export class TokensInfo extends jspb.Message { 
        getRemaining(): number;
        setRemaining(value: number): TokensInfo;
        getCurrent(): number;
        setCurrent(value: number): TokensInfo;
        getConsumed(): number;
        setConsumed(value: number): TokensInfo;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): TokensInfo.AsObject;
        static toObject(includeInstance: boolean, msg: TokensInfo): TokensInfo.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: TokensInfo, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): TokensInfo;
        static deserializeBinaryFromReader(message: TokensInfo, reader: jspb.BinaryReader): TokensInfo;
    }

    export namespace TokensInfo {
        export type AsObject = {
            remaining: number,
            current: number,
            consumed: number,
        }
    }

    export class RateLimiterInfo extends jspb.Message { 
        getLabel(): string;
        setLabel(value: string): RateLimiterInfo;

        hasTokensInfo(): boolean;
        clearTokensInfo(): void;
        getTokensInfo(): LimiterDecision.TokensInfo | undefined;
        setTokensInfo(value?: LimiterDecision.TokensInfo): RateLimiterInfo;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): RateLimiterInfo.AsObject;
        static toObject(includeInstance: boolean, msg: RateLimiterInfo): RateLimiterInfo.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: RateLimiterInfo, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): RateLimiterInfo;
        static deserializeBinaryFromReader(message: RateLimiterInfo, reader: jspb.BinaryReader): RateLimiterInfo;
    }

    export namespace RateLimiterInfo {
        export type AsObject = {
            label: string,
            tokensInfo?: LimiterDecision.TokensInfo.AsObject,
        }
    }

    export class SchedulerInfo extends jspb.Message { 
        getWorkloadIndex(): string;
        setWorkloadIndex(value: string): SchedulerInfo;

        hasTokensInfo(): boolean;
        clearTokensInfo(): void;
        getTokensInfo(): LimiterDecision.TokensInfo | undefined;
        setTokensInfo(value?: LimiterDecision.TokensInfo): SchedulerInfo;
        getPriority(): number;
        setPriority(value: number): SchedulerInfo;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): SchedulerInfo.AsObject;
        static toObject(includeInstance: boolean, msg: SchedulerInfo): SchedulerInfo.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: SchedulerInfo, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): SchedulerInfo;
        static deserializeBinaryFromReader(message: SchedulerInfo, reader: jspb.BinaryReader): SchedulerInfo;
    }

    export namespace SchedulerInfo {
        export type AsObject = {
            workloadIndex: string,
            tokensInfo?: LimiterDecision.TokensInfo.AsObject,
            priority: number,
        }
    }

    export class SamplerInfo extends jspb.Message { 
        getLabel(): string;
        setLabel(value: string): SamplerInfo;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): SamplerInfo.AsObject;
        static toObject(includeInstance: boolean, msg: SamplerInfo): SamplerInfo.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: SamplerInfo, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): SamplerInfo;
        static deserializeBinaryFromReader(message: SamplerInfo, reader: jspb.BinaryReader): SamplerInfo;
    }

    export namespace SamplerInfo {
        export type AsObject = {
            label: string,
        }
    }

    export class QuotaSchedulerInfo extends jspb.Message { 
        getLabel(): string;
        setLabel(value: string): QuotaSchedulerInfo;
        getWorkloadIndex(): string;
        setWorkloadIndex(value: string): QuotaSchedulerInfo;

        hasTokensInfo(): boolean;
        clearTokensInfo(): void;
        getTokensInfo(): LimiterDecision.TokensInfo | undefined;
        setTokensInfo(value?: LimiterDecision.TokensInfo): QuotaSchedulerInfo;
        getPriority(): number;
        setPriority(value: number): QuotaSchedulerInfo;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): QuotaSchedulerInfo.AsObject;
        static toObject(includeInstance: boolean, msg: QuotaSchedulerInfo): QuotaSchedulerInfo.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: QuotaSchedulerInfo, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): QuotaSchedulerInfo;
        static deserializeBinaryFromReader(message: QuotaSchedulerInfo, reader: jspb.BinaryReader): QuotaSchedulerInfo;
    }

    export namespace QuotaSchedulerInfo {
        export type AsObject = {
            label: string,
            workloadIndex: string,
            tokensInfo?: LimiterDecision.TokensInfo.AsObject,
            priority: number,
        }
    }


    export enum LimiterReason {
    LIMITER_REASON_UNSPECIFIED = 0,
    LIMITER_REASON_KEY_NOT_FOUND = 1,
    }


    export enum DetailsCase {
        DETAILS_NOT_SET = 0,
        RATE_LIMITER_INFO = 20,
        LOAD_SCHEDULER_INFO = 21,
        SAMPLER_INFO = 22,
        QUOTA_SCHEDULER_INFO = 23,
    }

}

export class FluxMeterInfo extends jspb.Message { 
    getFluxMeterName(): string;
    setFluxMeterName(value: string): FluxMeterInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): FluxMeterInfo.AsObject;
    static toObject(includeInstance: boolean, msg: FluxMeterInfo): FluxMeterInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: FluxMeterInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): FluxMeterInfo;
    static deserializeBinaryFromReader(message: FluxMeterInfo, reader: jspb.BinaryReader): FluxMeterInfo;
}

export namespace FluxMeterInfo {
    export type AsObject = {
        fluxMeterName: string,
    }
}

export enum StatusCode {
    EMPTY = 0,
    CONTINUE = 100,
    OK = 200,
    CREATED = 201,
    ACCEPTED = 202,
    NONAUTHORITATIVEINFORMATION = 203,
    NOCONTENT = 204,
    RESETCONTENT = 205,
    PARTIALCONTENT = 206,
    MULTISTATUS = 207,
    ALREADYREPORTED = 208,
    IMUSED = 226,
    MULTIPLECHOICES = 300,
    MOVEDPERMANENTLY = 301,
    FOUND = 302,
    SEEOTHER = 303,
    NOTMODIFIED = 304,
    USEPROXY = 305,
    TEMPORARYREDIRECT = 307,
    PERMANENTREDIRECT = 308,
    BADREQUEST = 400,
    UNAUTHORIZED = 401,
    PAYMENTREQUIRED = 402,
    FORBIDDEN = 403,
    NOTFOUND = 404,
    METHODNOTALLOWED = 405,
    NOTACCEPTABLE = 406,
    PROXYAUTHENTICATIONREQUIRED = 407,
    REQUESTTIMEOUT = 408,
    CONFLICT = 409,
    GONE = 410,
    LENGTHREQUIRED = 411,
    PRECONDITIONFAILED = 412,
    PAYLOADTOOLARGE = 413,
    URITOOLONG = 414,
    UNSUPPORTEDMEDIATYPE = 415,
    RANGENOTSATISFIABLE = 416,
    EXPECTATIONFAILED = 417,
    MISDIRECTEDREQUEST = 421,
    UNPROCESSABLEENTITY = 422,
    LOCKED = 423,
    FAILEDDEPENDENCY = 424,
    UPGRADEREQUIRED = 426,
    PRECONDITIONREQUIRED = 428,
    TOOMANYREQUESTS = 429,
    REQUESTHEADERFIELDSTOOLARGE = 431,
    INTERNALSERVERERROR = 500,
    NOTIMPLEMENTED = 501,
    BADGATEWAY = 502,
    SERVICEUNAVAILABLE = 503,
    GATEWAYTIMEOUT = 504,
    HTTPVERSIONNOTSUPPORTED = 505,
    VARIANTALSONEGOTIATES = 506,
    INSUFFICIENTSTORAGE = 507,
    LOOPDETECTED = 508,
    NOTEXTENDED = 510,
    NETWORKAUTHENTICATIONREQUIRED = 511,
}
