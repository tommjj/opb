package interfaces

type IStorage interface {
	Load(v any) error
	Sync(v any) error
}
