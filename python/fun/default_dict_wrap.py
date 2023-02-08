from collections import UserDict, defaultdict
from typing import (
    Any,
    Callable,
    Hashable,
    TypeVar
)


T = TypeVar('T')


class WrapDefaultDict(UserDict):
    def __init__(self, defaultfactory: Callable = list):
        self.data: defaultdict = defaultdict(defaultfactory)

    def __setitem__(self, k: Hashable, v: T) -> T:
        match self.data.default_factory.__name__:
            case 'list':
                self.data[k].append(v)
            case 'set':
                self.data[k].add(v)
            case _:
                self.data[k] = v

        return v

    def __getitem__(self, k: Hashable) -> Any:
        return self.data.__getitem__(k)

    def __missing__(self, key: Hashable) -> Any:
        return self.data.__missing__(key)

    @property
    def default_factory(self) -> Any:
        return self.data.default_factory