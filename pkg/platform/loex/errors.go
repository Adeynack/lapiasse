package loex

type Erroreable = func() error

func GetAllOrErr(f ...Erroreable) error {
	for _, oneF := range f {
		if err := oneF(); err != nil {
			return err
		}
	}

	return nil
}

func GetAllOrErr2[T1, T2 any](
	f1 func() (T1, error),
	f2 func() (T2, error),
) (T1, T2, error) {
	v1, err := f1()
	if err != nil {
		var zeroT2 T2
		return v1, zeroT2, err
	}

	v2, err := f2()
	if err != nil {
		return v1, v2, err
	}

	return v1, v2, nil
}
