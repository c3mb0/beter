package b

import "runtime"

// B is the detailed wrapper around an error.
type B struct {
	Err  error
	Fn   string
	File string
	Line int
}

func (e *B) Error() string {
	return e.Err.Error()
}

// E either wraps the supplied error with a B or returns it as is if it's already of type B.
func E(r error) error {
	if r == nil {
		return nil
	}
	switch t := r.(type) {
	case *B:
		return t
	}
	p, f, l, o := runtime.Caller(1)
	if !o {
		return &B{Err: r, Fn: "runtime error", File: "runtime error", Line: 0}
	}
	fn := runtime.FuncForPC(p)
	return &B{Err: r, Fn: fn.Name(), File: f, Line: l}
}
