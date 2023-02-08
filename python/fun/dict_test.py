from collections import defaultdict

from default_dict_wrap import WrapDefaultDict


def test_list():
    d = WrapDefaultDict()
    d[1] = 'a'
    assert d[1] == ['a']
    d[1] = 'A'
    assert d[1] == ['a', 'A']
    d[1].append('B')


def test_set():
    d = WrapDefaultDict(set)
    d[1] = 'a'
    assert d[1] == {'a'}
    d[1] = 'A'
    assert d[1] == {'a', 'A'}
    d[1].add('A')
    assert d[1] == {'a', 'A'}


def test_other_types():
    d = WrapDefaultDict(int)
    assert d[1] == 0
    d = WrapDefaultDict(str)
    assert d[1] == ""


def test_dict_methods():
    d = WrapDefaultDict()
    d[1] = 'a'
    d[1] = 'A'
    d[2] = 'b'
    d[3] = 'c'
    d[4] = 'd'
    keys = d.keys()
    assert tuple(keys) == (1, 2, 3, 4)
    values = d.values()
    assert tuple(values) == (['a', 'A'], ['b'], ['c'], ['d'])
    assert d.pop(3) == ['c']
    assert d.popitem() == (4, 'd')


def test_defaultdict_methods():
    wd = WrapDefaultDict(list)
    dd = defaultdict(list)
    for m in dir(dd):
        assert hasattr(wd, m)

    assert wd.default_factory == dd.default_factory
