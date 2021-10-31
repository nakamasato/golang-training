package main


type Dictionary map[string]string
type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}

const  (
	ErrWordDoesNotExist = DictionaryErr("could not find the word you were looking for")
	ErrWordExists = DictionaryErr("could not find the word you were looking for")
)

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrWordDoesNotExist
	}
	return definition, nil
}

func (d Dictionary) Add(key string, value string) {
	d[key] = value
}

func(d Dictionary) Update(key string, value string) error {
	_, err := d.Search(key)
	switch err {
	case ErrWordDoesNotExist:
		return ErrWordDoesNotExist
	case nil:
		d[key] = value
	default:
		return err
	}
	return nil
}
