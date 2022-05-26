package types

type (
	T interface{}
	R interface{}
	K interface{}

	KV struct {
		KEY   T
		VALUE T
	}
)
