package backend

type JsonStruct interface {
	Unmarshal([]byte) (JsonStruct, error)
}
