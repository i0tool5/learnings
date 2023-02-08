# defaultdict in Python

`defaultdict` is very useful data structure. It allows developer to set default value for every key in the dictionary.
```python
from collections import defaultdict
d = defaultdict(int)
assert d['a'] == 0  # and it works fine 
```

One of the most coolest features of `defaultdict` (naturally in my opinion) is to pass *`list`* of *`set`* type as *default_factory* argument, which allows us to do something like this:

```python
d = defaultdict(list)
d[1].append('a')
d[1].append('A')
d[1].append('1')
```

This is a very simple example, but it shows, how powerfull this construction is.

But i found, that using methods of underlying key type is uncomfortable (but explicit, yes). That's why i wrote simple wrapper around `UserDict` and `defaultdict` (and ofcourse for some fun!).

## Why tests??

Tests are reflection of the fact that everything works as expected, and there is no violation of dict interfaces.

## WHY TYPES??!!

I am working with Go and Rust, which are strongly typed languages. So, i am using typing in Python because it's my habit (and this is a good practice).

### Explanation

```python
def __setitem__(self, k: Hashable, v: T) -> T:
    match self.data.default_factory.__name__:
        case 'list':
            self.data[k].append(v)
        case 'set':
            self.data[k].add(v)
        case _:
            self.data[k] = v

    return v
```
