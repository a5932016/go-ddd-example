package mBinding

var (
	Query = mTagBinding{}
)

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
