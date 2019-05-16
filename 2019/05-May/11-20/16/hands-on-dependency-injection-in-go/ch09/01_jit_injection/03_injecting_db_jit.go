package jit_injection

import "errors"

type MyLoadPersonLogicJIT struct {
	dataSource DataSourceJIT
}

func (m *MyLoadPersonLogicJIT) Load(ID int) (Person, error) {
	return m.getDataSource().Load(ID)
}

func (m *MyLoadPersonLogicJIT) getDataSource() DataSourceJIT {
	if m.dataSource == nil {
		m.dataSource = NewMyDataSourceJIT()
	}
	return m.dataSource
}

type DataSourceJIT interface {
	Load(ID int) (Person, error)
}

func NewMyDataSourceJIT() *MyDataSourceJIT {
	return &MyDataSourceJIT{}
}

type MyDataSourceJIT struct {
}

func (m *MyDataSourceJIT) Load(ID int) (Person, error) {
	return Person{}, errors.New("not implemented yet")
}
