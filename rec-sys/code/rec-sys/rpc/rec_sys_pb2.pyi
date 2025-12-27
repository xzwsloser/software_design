from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class GetRecResultReq(_message.Message):
    __slots__ = ("userId", "addressId", "touristType", "priceSensitive", "likeType", "targetType", "attentionType", "update", "limit")
    USERID_FIELD_NUMBER: _ClassVar[int]
    ADDRESSID_FIELD_NUMBER: _ClassVar[int]
    TOURISTTYPE_FIELD_NUMBER: _ClassVar[int]
    PRICESENSITIVE_FIELD_NUMBER: _ClassVar[int]
    LIKETYPE_FIELD_NUMBER: _ClassVar[int]
    TARGETTYPE_FIELD_NUMBER: _ClassVar[int]
    ATTENTIONTYPE_FIELD_NUMBER: _ClassVar[int]
    UPDATE_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    userId: int
    addressId: int
    touristType: int
    priceSensitive: int
    likeType: _containers.RepeatedScalarFieldContainer[int]
    targetType: _containers.RepeatedScalarFieldContainer[int]
    attentionType: _containers.RepeatedScalarFieldContainer[int]
    update: bool
    limit: int
    def __init__(self, userId: _Optional[int] = ..., addressId: _Optional[int] = ..., touristType: _Optional[int] = ..., priceSensitive: _Optional[int] = ..., likeType: _Optional[_Iterable[int]] = ..., targetType: _Optional[_Iterable[int]] = ..., attentionType: _Optional[_Iterable[int]] = ..., update: bool = ..., limit: _Optional[int] = ...) -> None: ...

class GetRecResultResp(_message.Message):
    __slots__ = ("siteIdxList",)
    SITEIDXLIST_FIELD_NUMBER: _ClassVar[int]
    siteIdxList: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, siteIdxList: _Optional[_Iterable[int]] = ...) -> None: ...
