import enum

from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2 import (
    HIT,
    CacheLookupStatus,
)


class LookupStatus(enum.Enum):
    HIT = "HIT"
    MISS = "MISS"


def convert_cache_lookup_status(status: CacheLookupStatus):
    return LookupStatus.HIT if status == HIT else LookupStatus.MISS


# convertCacheError converts a string error message to Python's Exception type.
# Returns None if the input string is empty.
def convert_cache_error(error_message):
    return None if not error_message else Exception(error_message)


# KeyLookupResponse is the interface that represents a cache value lookup.
class KeyLookupResponse:
    def __init__(self, value, lookup_status: LookupStatus, error):
        self.value = value
        self.lookup_status = lookup_status
        self.error = error

    def get_value(self):
        return self.value

    def get_lookup_status(self) -> LookupStatus:
        return self.lookup_status

    def get_error(self):
        return self.error


# KeyUpsertResponse represents response of updating or inserting a cache key.
class KeyUpsertResponse:
    def __init__(self, error):
        self.error = error

    def get_error(self):
        return self.error


# DeleteCachedValueResponse represents the response of deleting a cache key.
class KeyDeleteResponse:
    def __init__(self, error):
        self.error = error

    def get_error(self):
        return self.error
