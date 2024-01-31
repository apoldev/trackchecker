// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	models "github.com/apoldev/trackchecker/internal/app/models"
	mock "github.com/stretchr/testify/mock"
)

// Crawler is an autogenerated mock type for the Crawler type
type Crawler struct {
	mock.Mock
}

// Start provides a mock function with given fields: number
func (_m *Crawler) Start(number *models.TrackingNumber) (*models.Crawler, error) {
	ret := _m.Called(number)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 *models.Crawler
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.TrackingNumber) (*models.Crawler, error)); ok {
		return rf(number)
	}
	if rf, ok := ret.Get(0).(func(*models.TrackingNumber) *models.Crawler); ok {
		r0 = rf(number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Crawler)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.TrackingNumber) error); ok {
		r1 = rf(number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCrawler creates a new instance of Crawler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCrawler(t interface {
	mock.TestingT
	Cleanup(func())
}) *Crawler {
	mock := &Crawler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
