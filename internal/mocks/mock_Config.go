// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	language "golang.org/x/text/language"

	time "time"
)

// MockConfig is an autogenerated mock type for the Config type
type MockConfig struct {
	mock.Mock
}

type MockConfig_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConfig) EXPECT() *MockConfig_Expecter {
	return &MockConfig_Expecter{mock: &_m.Mock}
}

// All provides a mock function with given fields:
func (_m *MockConfig) All() map[string]interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for All")
	}

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// MockConfig_All_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'All'
type MockConfig_All_Call struct {
	*mock.Call
}

// All is a helper method to define mock.On call
func (_e *MockConfig_Expecter) All() *MockConfig_All_Call {
	return &MockConfig_All_Call{Call: _e.mock.On("All")}
}

func (_c *MockConfig_All_Call) Run(run func()) *MockConfig_All_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_All_Call) Return(_a0 map[string]interface{}) *MockConfig_All_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_All_Call) RunAndReturn(run func() map[string]interface{}) *MockConfig_All_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0
func (_m *MockConfig) Get(_a0 string) interface{} {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// MockConfig_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockConfig_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 string
func (_e *MockConfig_Expecter) Get(_a0 interface{}) *MockConfig_Get_Call {
	return &MockConfig_Get_Call{Call: _e.mock.On("Get", _a0)}
}

func (_c *MockConfig_Get_Call) Run(run func(_a0 string)) *MockConfig_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockConfig_Get_Call) Return(_a0 interface{}) *MockConfig_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_Get_Call) RunAndReturn(run func(string) interface{}) *MockConfig_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetBool provides a mock function with given fields: _a0
func (_m *MockConfig) GetBool(_a0 string) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetBool")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockConfig_GetBool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBool'
type MockConfig_GetBool_Call struct {
	*mock.Call
}

// GetBool is a helper method to define mock.On call
//   - _a0 string
func (_e *MockConfig_Expecter) GetBool(_a0 interface{}) *MockConfig_GetBool_Call {
	return &MockConfig_GetBool_Call{Call: _e.mock.On("GetBool", _a0)}
}

func (_c *MockConfig_GetBool_Call) Run(run func(_a0 string)) *MockConfig_GetBool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockConfig_GetBool_Call) Return(_a0 bool) *MockConfig_GetBool_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_GetBool_Call) RunAndReturn(run func(string) bool) *MockConfig_GetBool_Call {
	_c.Call.Return(run)
	return _c
}

// GetInt provides a mock function with given fields: _a0
func (_m *MockConfig) GetInt(_a0 string) int {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetInt")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockConfig_GetInt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetInt'
type MockConfig_GetInt_Call struct {
	*mock.Call
}

// GetInt is a helper method to define mock.On call
//   - _a0 string
func (_e *MockConfig_Expecter) GetInt(_a0 interface{}) *MockConfig_GetInt_Call {
	return &MockConfig_GetInt_Call{Call: _e.mock.On("GetInt", _a0)}
}

func (_c *MockConfig_GetInt_Call) Run(run func(_a0 string)) *MockConfig_GetInt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockConfig_GetInt_Call) Return(_a0 int) *MockConfig_GetInt_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_GetInt_Call) RunAndReturn(run func(string) int) *MockConfig_GetInt_Call {
	_c.Call.Return(run)
	return _c
}

// GetString provides a mock function with given fields: _a0
func (_m *MockConfig) GetString(_a0 string) string {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetString")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockConfig_GetString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetString'
type MockConfig_GetString_Call struct {
	*mock.Call
}

// GetString is a helper method to define mock.On call
//   - _a0 string
func (_e *MockConfig_Expecter) GetString(_a0 interface{}) *MockConfig_GetString_Call {
	return &MockConfig_GetString_Call{Call: _e.mock.On("GetString", _a0)}
}

func (_c *MockConfig_GetString_Call) Run(run func(_a0 string)) *MockConfig_GetString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockConfig_GetString_Call) Return(_a0 string) *MockConfig_GetString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_GetString_Call) RunAndReturn(run func(string) string) *MockConfig_GetString_Call {
	_c.Call.Return(run)
	return _c
}

// GetStringSlice provides a mock function with given fields: _a0
func (_m *MockConfig) GetStringSlice(_a0 string) []string {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetStringSlice")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MockConfig_GetStringSlice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStringSlice'
type MockConfig_GetStringSlice_Call struct {
	*mock.Call
}

// GetStringSlice is a helper method to define mock.On call
//   - _a0 string
func (_e *MockConfig_Expecter) GetStringSlice(_a0 interface{}) *MockConfig_GetStringSlice_Call {
	return &MockConfig_GetStringSlice_Call{Call: _e.mock.On("GetStringSlice", _a0)}
}

func (_c *MockConfig_GetStringSlice_Call) Run(run func(_a0 string)) *MockConfig_GetStringSlice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockConfig_GetStringSlice_Call) Return(_a0 []string) *MockConfig_GetStringSlice_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_GetStringSlice_Call) RunAndReturn(run func(string) []string) *MockConfig_GetStringSlice_Call {
	_c.Call.Return(run)
	return _c
}

// GetWorkWeekdays provides a mock function with given fields:
func (_m *MockConfig) GetWorkWeekdays() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetWorkWeekdays")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MockConfig_GetWorkWeekdays_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorkWeekdays'
type MockConfig_GetWorkWeekdays_Call struct {
	*mock.Call
}

// GetWorkWeekdays is a helper method to define mock.On call
func (_e *MockConfig_Expecter) GetWorkWeekdays() *MockConfig_GetWorkWeekdays_Call {
	return &MockConfig_GetWorkWeekdays_Call{Call: _e.mock.On("GetWorkWeekdays")}
}

func (_c *MockConfig_GetWorkWeekdays_Call) Run(run func()) *MockConfig_GetWorkWeekdays_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_GetWorkWeekdays_Call) Return(_a0 []string) *MockConfig_GetWorkWeekdays_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_GetWorkWeekdays_Call) RunAndReturn(run func() []string) *MockConfig_GetWorkWeekdays_Call {
	_c.Call.Return(run)
	return _c
}

// InteractivePageSize provides a mock function with given fields:
func (_m *MockConfig) InteractivePageSize() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for InteractivePageSize")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockConfig_InteractivePageSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InteractivePageSize'
type MockConfig_InteractivePageSize_Call struct {
	*mock.Call
}

// InteractivePageSize is a helper method to define mock.On call
func (_e *MockConfig_Expecter) InteractivePageSize() *MockConfig_InteractivePageSize_Call {
	return &MockConfig_InteractivePageSize_Call{Call: _e.mock.On("InteractivePageSize")}
}

func (_c *MockConfig_InteractivePageSize_Call) Run(run func()) *MockConfig_InteractivePageSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_InteractivePageSize_Call) Return(_a0 int) *MockConfig_InteractivePageSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_InteractivePageSize_Call) RunAndReturn(run func() int) *MockConfig_InteractivePageSize_Call {
	_c.Call.Return(run)
	return _c
}

// IsAllowNameForID provides a mock function with given fields:
func (_m *MockConfig) IsAllowNameForID() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsAllowNameForID")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockConfig_IsAllowNameForID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsAllowNameForID'
type MockConfig_IsAllowNameForID_Call struct {
	*mock.Call
}

// IsAllowNameForID is a helper method to define mock.On call
func (_e *MockConfig_Expecter) IsAllowNameForID() *MockConfig_IsAllowNameForID_Call {
	return &MockConfig_IsAllowNameForID_Call{Call: _e.mock.On("IsAllowNameForID")}
}

func (_c *MockConfig_IsAllowNameForID_Call) Run(run func()) *MockConfig_IsAllowNameForID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_IsAllowNameForID_Call) Return(_a0 bool) *MockConfig_IsAllowNameForID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_IsAllowNameForID_Call) RunAndReturn(run func() bool) *MockConfig_IsAllowNameForID_Call {
	_c.Call.Return(run)
	return _c
}

// IsDebuging provides a mock function with given fields:
func (_m *MockConfig) IsDebuging() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsDebuging")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockConfig_IsDebuging_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsDebuging'
type MockConfig_IsDebuging_Call struct {
	*mock.Call
}

// IsDebuging is a helper method to define mock.On call
func (_e *MockConfig_Expecter) IsDebuging() *MockConfig_IsDebuging_Call {
	return &MockConfig_IsDebuging_Call{Call: _e.mock.On("IsDebuging")}
}

func (_c *MockConfig_IsDebuging_Call) Run(run func()) *MockConfig_IsDebuging_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_IsDebuging_Call) Return(_a0 bool) *MockConfig_IsDebuging_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_IsDebuging_Call) RunAndReturn(run func() bool) *MockConfig_IsDebuging_Call {
	_c.Call.Return(run)
	return _c
}

// IsInteractive provides a mock function with given fields:
func (_m *MockConfig) IsInteractive() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsInteractive")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockConfig_IsInteractive_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsInteractive'
type MockConfig_IsInteractive_Call struct {
	*mock.Call
}

// IsInteractive is a helper method to define mock.On call
func (_e *MockConfig_Expecter) IsInteractive() *MockConfig_IsInteractive_Call {
	return &MockConfig_IsInteractive_Call{Call: _e.mock.On("IsInteractive")}
}

func (_c *MockConfig_IsInteractive_Call) Run(run func()) *MockConfig_IsInteractive_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_IsInteractive_Call) Return(_a0 bool) *MockConfig_IsInteractive_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_IsInteractive_Call) RunAndReturn(run func() bool) *MockConfig_IsInteractive_Call {
	_c.Call.Return(run)
	return _c
}

// IsSearchProjectWithClientsName provides a mock function with given fields:
func (_m *MockConfig) IsSearchProjectWithClientsName() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsSearchProjectWithClientsName")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockConfig_IsSearchProjectWithClientsName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsSearchProjectWithClientsName'
type MockConfig_IsSearchProjectWithClientsName_Call struct {
	*mock.Call
}

// IsSearchProjectWithClientsName is a helper method to define mock.On call
func (_e *MockConfig_Expecter) IsSearchProjectWithClientsName() *MockConfig_IsSearchProjectWithClientsName_Call {
	return &MockConfig_IsSearchProjectWithClientsName_Call{Call: _e.mock.On("IsSearchProjectWithClientsName")}
}

func (_c *MockConfig_IsSearchProjectWithClientsName_Call) Run(run func()) *MockConfig_IsSearchProjectWithClientsName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_IsSearchProjectWithClientsName_Call) Return(_a0 bool) *MockConfig_IsSearchProjectWithClientsName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_IsSearchProjectWithClientsName_Call) RunAndReturn(run func() bool) *MockConfig_IsSearchProjectWithClientsName_Call {
	_c.Call.Return(run)
	return _c
}

// Language provides a mock function with given fields:
func (_m *MockConfig) Language() language.Tag {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Language")
	}

	var r0 language.Tag
	if rf, ok := ret.Get(0).(func() language.Tag); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(language.Tag)
	}

	return r0
}

// MockConfig_Language_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Language'
type MockConfig_Language_Call struct {
	*mock.Call
}

// Language is a helper method to define mock.On call
func (_e *MockConfig_Expecter) Language() *MockConfig_Language_Call {
	return &MockConfig_Language_Call{Call: _e.mock.On("Language")}
}

func (_c *MockConfig_Language_Call) Run(run func()) *MockConfig_Language_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_Language_Call) Return(_a0 language.Tag) *MockConfig_Language_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_Language_Call) RunAndReturn(run func() language.Tag) *MockConfig_Language_Call {
	_c.Call.Return(run)
	return _c
}

// LogLevel provides a mock function with given fields:
func (_m *MockConfig) LogLevel() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LogLevel")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockConfig_LogLevel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LogLevel'
type MockConfig_LogLevel_Call struct {
	*mock.Call
}

// LogLevel is a helper method to define mock.On call
func (_e *MockConfig_Expecter) LogLevel() *MockConfig_LogLevel_Call {
	return &MockConfig_LogLevel_Call{Call: _e.mock.On("LogLevel")}
}

func (_c *MockConfig_LogLevel_Call) Run(run func()) *MockConfig_LogLevel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_LogLevel_Call) Return(_a0 string) *MockConfig_LogLevel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_LogLevel_Call) RunAndReturn(run func() string) *MockConfig_LogLevel_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields:
func (_m *MockConfig) Save() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockConfig_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockConfig_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
func (_e *MockConfig_Expecter) Save() *MockConfig_Save_Call {
	return &MockConfig_Save_Call{Call: _e.mock.On("Save")}
}

func (_c *MockConfig_Save_Call) Run(run func()) *MockConfig_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_Save_Call) Return(_a0 error) *MockConfig_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_Save_Call) RunAndReturn(run func() error) *MockConfig_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SetBool provides a mock function with given fields: _a0, _a1
func (_m *MockConfig) SetBool(_a0 string, _a1 bool) {
	_m.Called(_a0, _a1)
}

// MockConfig_SetBool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetBool'
type MockConfig_SetBool_Call struct {
	*mock.Call
}

// SetBool is a helper method to define mock.On call
//   - _a0 string
//   - _a1 bool
func (_e *MockConfig_Expecter) SetBool(_a0 interface{}, _a1 interface{}) *MockConfig_SetBool_Call {
	return &MockConfig_SetBool_Call{Call: _e.mock.On("SetBool", _a0, _a1)}
}

func (_c *MockConfig_SetBool_Call) Run(run func(_a0 string, _a1 bool)) *MockConfig_SetBool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(bool))
	})
	return _c
}

func (_c *MockConfig_SetBool_Call) Return() *MockConfig_SetBool_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockConfig_SetBool_Call) RunAndReturn(run func(string, bool)) *MockConfig_SetBool_Call {
	_c.Call.Return(run)
	return _c
}

// SetInt provides a mock function with given fields: _a0, _a1
func (_m *MockConfig) SetInt(_a0 string, _a1 int) {
	_m.Called(_a0, _a1)
}

// MockConfig_SetInt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetInt'
type MockConfig_SetInt_Call struct {
	*mock.Call
}

// SetInt is a helper method to define mock.On call
//   - _a0 string
//   - _a1 int
func (_e *MockConfig_Expecter) SetInt(_a0 interface{}, _a1 interface{}) *MockConfig_SetInt_Call {
	return &MockConfig_SetInt_Call{Call: _e.mock.On("SetInt", _a0, _a1)}
}

func (_c *MockConfig_SetInt_Call) Run(run func(_a0 string, _a1 int)) *MockConfig_SetInt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(int))
	})
	return _c
}

func (_c *MockConfig_SetInt_Call) Return() *MockConfig_SetInt_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockConfig_SetInt_Call) RunAndReturn(run func(string, int)) *MockConfig_SetInt_Call {
	_c.Call.Return(run)
	return _c
}

// SetLanguage provides a mock function with given fields: _a0
func (_m *MockConfig) SetLanguage(_a0 language.Tag) {
	_m.Called(_a0)
}

// MockConfig_SetLanguage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetLanguage'
type MockConfig_SetLanguage_Call struct {
	*mock.Call
}

// SetLanguage is a helper method to define mock.On call
//   - _a0 language.Tag
func (_e *MockConfig_Expecter) SetLanguage(_a0 interface{}) *MockConfig_SetLanguage_Call {
	return &MockConfig_SetLanguage_Call{Call: _e.mock.On("SetLanguage", _a0)}
}

func (_c *MockConfig_SetLanguage_Call) Run(run func(_a0 language.Tag)) *MockConfig_SetLanguage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(language.Tag))
	})
	return _c
}

func (_c *MockConfig_SetLanguage_Call) Return() *MockConfig_SetLanguage_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockConfig_SetLanguage_Call) RunAndReturn(run func(language.Tag)) *MockConfig_SetLanguage_Call {
	_c.Call.Return(run)
	return _c
}

// SetString provides a mock function with given fields: _a0, _a1
func (_m *MockConfig) SetString(_a0 string, _a1 string) {
	_m.Called(_a0, _a1)
}

// MockConfig_SetString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetString'
type MockConfig_SetString_Call struct {
	*mock.Call
}

// SetString is a helper method to define mock.On call
//   - _a0 string
//   - _a1 string
func (_e *MockConfig_Expecter) SetString(_a0 interface{}, _a1 interface{}) *MockConfig_SetString_Call {
	return &MockConfig_SetString_Call{Call: _e.mock.On("SetString", _a0, _a1)}
}

func (_c *MockConfig_SetString_Call) Run(run func(_a0 string, _a1 string)) *MockConfig_SetString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockConfig_SetString_Call) Return() *MockConfig_SetString_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockConfig_SetString_Call) RunAndReturn(run func(string, string)) *MockConfig_SetString_Call {
	_c.Call.Return(run)
	return _c
}

// SetStringSlice provides a mock function with given fields: _a0, _a1
func (_m *MockConfig) SetStringSlice(_a0 string, _a1 []string) {
	_m.Called(_a0, _a1)
}

// MockConfig_SetStringSlice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStringSlice'
type MockConfig_SetStringSlice_Call struct {
	*mock.Call
}

// SetStringSlice is a helper method to define mock.On call
//   - _a0 string
//   - _a1 []string
func (_e *MockConfig_Expecter) SetStringSlice(_a0 interface{}, _a1 interface{}) *MockConfig_SetStringSlice_Call {
	return &MockConfig_SetStringSlice_Call{Call: _e.mock.On("SetStringSlice", _a0, _a1)}
}

func (_c *MockConfig_SetStringSlice_Call) Run(run func(_a0 string, _a1 []string)) *MockConfig_SetStringSlice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]string))
	})
	return _c
}

func (_c *MockConfig_SetStringSlice_Call) Return() *MockConfig_SetStringSlice_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockConfig_SetStringSlice_Call) RunAndReturn(run func(string, []string)) *MockConfig_SetStringSlice_Call {
	_c.Call.Return(run)
	return _c
}

// SetTimeZone provides a mock function with given fields: _a0
func (_m *MockConfig) SetTimeZone(_a0 *time.Location) {
	_m.Called(_a0)
}

// MockConfig_SetTimeZone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTimeZone'
type MockConfig_SetTimeZone_Call struct {
	*mock.Call
}

// SetTimeZone is a helper method to define mock.On call
//   - _a0 *time.Location
func (_e *MockConfig_Expecter) SetTimeZone(_a0 interface{}) *MockConfig_SetTimeZone_Call {
	return &MockConfig_SetTimeZone_Call{Call: _e.mock.On("SetTimeZone", _a0)}
}

func (_c *MockConfig_SetTimeZone_Call) Run(run func(_a0 *time.Location)) *MockConfig_SetTimeZone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*time.Location))
	})
	return _c
}

func (_c *MockConfig_SetTimeZone_Call) Return() *MockConfig_SetTimeZone_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockConfig_SetTimeZone_Call) RunAndReturn(run func(*time.Location)) *MockConfig_SetTimeZone_Call {
	_c.Call.Return(run)
	return _c
}

// TimeZone provides a mock function with given fields:
func (_m *MockConfig) TimeZone() *time.Location {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TimeZone")
	}

	var r0 *time.Location
	if rf, ok := ret.Get(0).(func() *time.Location); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*time.Location)
		}
	}

	return r0
}

// MockConfig_TimeZone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TimeZone'
type MockConfig_TimeZone_Call struct {
	*mock.Call
}

// TimeZone is a helper method to define mock.On call
func (_e *MockConfig_Expecter) TimeZone() *MockConfig_TimeZone_Call {
	return &MockConfig_TimeZone_Call{Call: _e.mock.On("TimeZone")}
}

func (_c *MockConfig_TimeZone_Call) Run(run func()) *MockConfig_TimeZone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfig_TimeZone_Call) Return(_a0 *time.Location) *MockConfig_TimeZone_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfig_TimeZone_Call) RunAndReturn(run func() *time.Location) *MockConfig_TimeZone_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConfig creates a new instance of MockConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConfig(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConfig {
	mock := &MockConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
