package loex

type Erroreable = func() error

func TryAll(f ...Erroreable) error {
	for _, oneF := range f {
		if err := oneF(); err != nil {
			return err
		}
	}

	return nil
}
