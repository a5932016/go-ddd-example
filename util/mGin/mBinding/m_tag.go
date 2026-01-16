package mBinding

import "net/http"

const mTag = "mTag"

type mTagBinding struct{}

func (mTagBinding) Name() string {
	return mTag
}

func (mTagBinding) Bind(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	if err := mapFormByTag(obj, values, mTag); err != nil {
		return err
	}
	return validate(obj)
}
