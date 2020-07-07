// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common/networkingcommon (interfaces: LinkLayerDevice,LinkLayerAddress)

// Package instancepoller is a generated GoMock package.
package instancepoller

import (
	gomock "github.com/golang/mock/gomock"
	network "github.com/juju/juju/core/network"
	state "github.com/juju/juju/state"
	txn "gopkg.in/mgo.v2/txn"
	reflect "reflect"
)

// MockLinkLayerDevice is a mock of LinkLayerDevice interface
type MockLinkLayerDevice struct {
	ctrl     *gomock.Controller
	recorder *MockLinkLayerDeviceMockRecorder
}

// MockLinkLayerDeviceMockRecorder is the mock recorder for MockLinkLayerDevice
type MockLinkLayerDeviceMockRecorder struct {
	mock *MockLinkLayerDevice
}

// NewMockLinkLayerDevice creates a new mock instance
func NewMockLinkLayerDevice(ctrl *gomock.Controller) *MockLinkLayerDevice {
	mock := &MockLinkLayerDevice{ctrl: ctrl}
	mock.recorder = &MockLinkLayerDeviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLinkLayerDevice) EXPECT() *MockLinkLayerDeviceMockRecorder {
	return m.recorder
}

// DocID mocks base method
func (m *MockLinkLayerDevice) DocID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DocID")
	ret0, _ := ret[0].(string)
	return ret0
}

// DocID indicates an expected call of DocID
func (mr *MockLinkLayerDeviceMockRecorder) DocID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DocID", reflect.TypeOf((*MockLinkLayerDevice)(nil).DocID))
}

// MACAddress mocks base method
func (m *MockLinkLayerDevice) MACAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MACAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// MACAddress indicates an expected call of MACAddress
func (mr *MockLinkLayerDeviceMockRecorder) MACAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MACAddress", reflect.TypeOf((*MockLinkLayerDevice)(nil).MACAddress))
}

// Name mocks base method
func (m *MockLinkLayerDevice) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockLinkLayerDeviceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockLinkLayerDevice)(nil).Name))
}

// ParentID mocks base method
func (m *MockLinkLayerDevice) ParentID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParentID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ParentID indicates an expected call of ParentID
func (mr *MockLinkLayerDeviceMockRecorder) ParentID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParentID", reflect.TypeOf((*MockLinkLayerDevice)(nil).ParentID))
}

// ProviderID mocks base method
func (m *MockLinkLayerDevice) ProviderID() network.Id {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProviderID")
	ret0, _ := ret[0].(network.Id)
	return ret0
}

// ProviderID indicates an expected call of ProviderID
func (mr *MockLinkLayerDeviceMockRecorder) ProviderID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProviderID", reflect.TypeOf((*MockLinkLayerDevice)(nil).ProviderID))
}

// RemoveOps mocks base method
func (m *MockLinkLayerDevice) RemoveOps() []txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOps")
	ret0, _ := ret[0].([]txn.Op)
	return ret0
}

// RemoveOps indicates an expected call of RemoveOps
func (mr *MockLinkLayerDeviceMockRecorder) RemoveOps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOps", reflect.TypeOf((*MockLinkLayerDevice)(nil).RemoveOps))
}

// SetProviderIDOps mocks base method
func (m *MockLinkLayerDevice) SetProviderIDOps(arg0 network.Id) ([]txn.Op, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProviderIDOps", arg0)
	ret0, _ := ret[0].([]txn.Op)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetProviderIDOps indicates an expected call of SetProviderIDOps
func (mr *MockLinkLayerDeviceMockRecorder) SetProviderIDOps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProviderIDOps", reflect.TypeOf((*MockLinkLayerDevice)(nil).SetProviderIDOps), arg0)
}

// UpdateOps mocks base method
func (m *MockLinkLayerDevice) UpdateOps(arg0 state.LinkLayerDeviceArgs) []txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOps", arg0)
	ret0, _ := ret[0].([]txn.Op)
	return ret0
}

// UpdateOps indicates an expected call of UpdateOps
func (mr *MockLinkLayerDeviceMockRecorder) UpdateOps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOps", reflect.TypeOf((*MockLinkLayerDevice)(nil).UpdateOps), arg0)
}

// MockLinkLayerAddress is a mock of LinkLayerAddress interface
type MockLinkLayerAddress struct {
	ctrl     *gomock.Controller
	recorder *MockLinkLayerAddressMockRecorder
}

// MockLinkLayerAddressMockRecorder is the mock recorder for MockLinkLayerAddress
type MockLinkLayerAddressMockRecorder struct {
	mock *MockLinkLayerAddress
}

// NewMockLinkLayerAddress creates a new mock instance
func NewMockLinkLayerAddress(ctrl *gomock.Controller) *MockLinkLayerAddress {
	mock := &MockLinkLayerAddress{ctrl: ctrl}
	mock.recorder = &MockLinkLayerAddressMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLinkLayerAddress) EXPECT() *MockLinkLayerAddressMockRecorder {
	return m.recorder
}

// DeviceName mocks base method
func (m *MockLinkLayerAddress) DeviceName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeviceName")
	ret0, _ := ret[0].(string)
	return ret0
}

// DeviceName indicates an expected call of DeviceName
func (mr *MockLinkLayerAddressMockRecorder) DeviceName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeviceName", reflect.TypeOf((*MockLinkLayerAddress)(nil).DeviceName))
}

// Origin mocks base method
func (m *MockLinkLayerAddress) Origin() network.Origin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Origin")
	ret0, _ := ret[0].(network.Origin)
	return ret0
}

// Origin indicates an expected call of Origin
func (mr *MockLinkLayerAddressMockRecorder) Origin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Origin", reflect.TypeOf((*MockLinkLayerAddress)(nil).Origin))
}

// RemoveOps mocks base method
func (m *MockLinkLayerAddress) RemoveOps() []txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOps")
	ret0, _ := ret[0].([]txn.Op)
	return ret0
}

// RemoveOps indicates an expected call of RemoveOps
func (mr *MockLinkLayerAddressMockRecorder) RemoveOps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOps", reflect.TypeOf((*MockLinkLayerAddress)(nil).RemoveOps))
}

// SetOriginOps mocks base method
func (m *MockLinkLayerAddress) SetOriginOps(arg0 network.Origin) []txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOriginOps", arg0)
	ret0, _ := ret[0].([]txn.Op)
	return ret0
}

// SetOriginOps indicates an expected call of SetOriginOps
func (mr *MockLinkLayerAddressMockRecorder) SetOriginOps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOriginOps", reflect.TypeOf((*MockLinkLayerAddress)(nil).SetOriginOps), arg0)
}

// SetProviderIDOps mocks base method
func (m *MockLinkLayerAddress) SetProviderIDOps(arg0 network.Id) ([]txn.Op, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProviderIDOps", arg0)
	ret0, _ := ret[0].([]txn.Op)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetProviderIDOps indicates an expected call of SetProviderIDOps
func (mr *MockLinkLayerAddressMockRecorder) SetProviderIDOps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProviderIDOps", reflect.TypeOf((*MockLinkLayerAddress)(nil).SetProviderIDOps), arg0)
}

// SetProviderNetIDsOps mocks base method
func (m *MockLinkLayerAddress) SetProviderNetIDsOps(arg0, arg1 network.Id) []txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProviderNetIDsOps", arg0, arg1)
	ret0, _ := ret[0].([]txn.Op)
	return ret0
}

// SetProviderNetIDsOps indicates an expected call of SetProviderNetIDsOps
func (mr *MockLinkLayerAddressMockRecorder) SetProviderNetIDsOps(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProviderNetIDsOps", reflect.TypeOf((*MockLinkLayerAddress)(nil).SetProviderNetIDsOps), arg0, arg1)
}

// Value mocks base method
func (m *MockLinkLayerAddress) Value() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value")
	ret0, _ := ret[0].(string)
	return ret0
}

// Value indicates an expected call of Value
func (mr *MockLinkLayerAddressMockRecorder) Value() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockLinkLayerAddress)(nil).Value))
}
