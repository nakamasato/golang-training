package main

type Dictionary map[string]string
type DictionaryErr string

var (
	ErrNotFound         = DictionaryErr("could not find the word you were looking for")
	ErrWordExists       = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("word doesn't not exist")
)

func (d Dictionary) Search(key string) (string, error) {
	definition, ok := d[key]
	if !ok {
		return d[key], ErrNotFound
	}
	return definition, nil
}

func (d Dictionary) Add(key, value string) error {
	_, err := d.Search(key)
	switch err {
	case ErrNotFound:
		d[key] = value
	case nil:
		return ErrWordExists
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(key, value string) error {
	switch _, err := d.Search(key); err {
	case ErrNotFound:
		return ErrWordDoesNotExist
	case nil:
		d[key] = value
		return nil
	default:
		return err
	}
}

func (d Dictionary) Delete(key string) {
	delete(d, key)
}

func (e DictionaryErr) Error() string {
	return string(e)
}
