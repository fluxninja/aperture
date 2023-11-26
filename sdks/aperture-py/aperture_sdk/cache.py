import enum

from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2 import (
    HIT,
    SUCCESS,
    CacheLookupStatus,
    CacheOperationStatus,
)


class LookupStatus(enum.Enum):
    HIT = "HIT"
    MISS = "MISS"


class OperationStatus(enum.Enum):
    SUCCESS = "SUCCESS"
    ERROR = "ERROR"


def convert_cache_lookup_status(status: CacheLookupStatus):
    return LookupStatus.HIT if status == HIT else LookupStatus.MISS


def convert_cache_operation_status(status: CacheOperationStatus):
    return OperationStatus.SUCCESS if status == SUCCESS else OperationStatus.ERROR


# convertCacheError converts a string error message to Python's Exception type.
# Returns None if the input string is empty.
def convert_cache_error(error_message):
    return None if not error_message else Exception(error_message)


# GetCachedValueResponse is the interface to read the response from a get cached value operation.
class GetCachedValueResponse:
    def __init__(self, value, lookup_status, operation_status, error):
        self.value = value
        self.lookup_status = lookup_status
        self.operation_status = operation_status
        self.error = error

    def get_value(self):
        return self.value

    def get_lookup_status(self):
        return self.lookup_status

    def get_operation_status(self):
        return self.operation_status

    def get_error(self):
        return self.error


# SetCachedValueResponse is the interface to read the response from a set cached value operation.
class SetCachedValueResponse:
    def __init__(self, operation_status, error):
        self.operation_status = operation_status
        self.error = error

    def get_error(self):
        return self.error

    def get_operation_status(self):
        return self.operation_status


# DeleteCachedValueResponse is the interface to read the response from a delete cached value operation.
class DeleteCachedValueResponse:
    def __init__(self, operation_status, error):
        self.operation_status = operation_status
        self.error = error

    def get_error(self):
        return self.error

    def get_operation_status(self):
        return self.operation_status
