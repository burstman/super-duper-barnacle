package datastorage

type Datastore struct {
	//name    string
	Data    interface{}
	Storage Datastorage
}
type Datastorage interface {
	List() string
	Load(string)
	Save(string, interface{})
}

func (d *Datastore) Load(name string) {
	d.Storage.Load(name)
}
func (d *Datastore) Save(name string) {
	d.Storage.Save(name, d.Data)
}
func (d *Datastore) List() string {
	return d.Storage.List()
}
