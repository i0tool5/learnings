package reflectnames

import "strings"

type call func() string
type decorator func(call) call

type Iface interface {
	CallFnOne() string
	CallFnTwo() string
	CallFnThree() string
}

type Some struct{}

func (s Some) CallFnOne() string {
	return "CallFnOne not decorated"
}
func (s Some) CallFnTwo() string {
	return "CallFnTwo not decorated"
}
func (s Some) CallFnThree() string {
	return "CallFnThree not decorated"
}

func SomeDecorator(name string) decorator {
	return func(fn call) call {
		return func() string {
			sb := new(strings.Builder)
			sb.WriteString("Calling decorated ")
			sb.WriteString(name)
			sb.WriteString("\t")
			sb.WriteString(
				strings.ReplaceAll(fn(), "not", "was"),
			)
			return sb.String()
		}
	}
}

type SomeDecorated struct {
	wraped map[string]call
}

func NewDecorated(i Iface) *SomeDecorated {
	namedWrapers := addNamedDeco(i)
	sd := new(SomeDecorated)
	sd.wraped = make(map[string]call, len(namedWrapers))
	sd.wraped["CallFnOne"] = namedWrapers["CallFnOne"](i.CallFnOne)
	sd.wraped["CallFnTwo"] = namedWrapers["CallFnTwo"](i.CallFnTwo)
	sd.wraped["CallFnThree"] = namedWrapers["CallFnThree"](i.CallFnThree)

	return sd
}

func (sd *SomeDecorated) CallFnOne() string {
	return sd.wraped["CallFnOne"]()
}

func (sd *SomeDecorated) CallFnTwo() string {
	return sd.wraped["CallFnTwo"]()
}
func (sd *SomeDecorated) CallFnThree() string {
	return sd.wraped["CallFnThree"]()
}

func addNamedDeco(struc any) map[string]decorator {
	dec := make(map[string]decorator)
	for _, name := range GetMethodNames(struc) {
		dec[name] = SomeDecorator(name)
	}

	return dec
}
