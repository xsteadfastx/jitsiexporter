// Code generated by mockery v1.0.0. DO NOT EDIT.

package jitsiexporter

import mock "github.com/stretchr/testify/mock"

// MockStater is an autogenerated mock type for the Stater type
type MockStater struct {
	mock.Mock
}

// Now provides a mock function with given fields: url
func (_m *MockStater) Now(url string) map[string]interface{} {
	ret := _m.Called(url)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(string) map[string]interface{}); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}
