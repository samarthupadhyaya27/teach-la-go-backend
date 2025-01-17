package db

import (
	"context"
	"errors"
)

type MockDB struct {
	// "Users, Programs, Class" collection
	db map[string]map[string]interface{}
}

func (d *MockDB) LoadProgram(_ context.Context, pid string) (Program, error) {
	p, ok := d.db[programsPath][pid].(Program)
	if !ok {
		return Program{}, errors.New("program has not been created")
	}
	return p, nil
}

func (d *MockDB) StoreProgram(_ context.Context, p Program) error {
	d.db[programsPath][p.UID] = p
	return nil
}

func (d *MockDB) RemoveProgram(_ context.Context, pid string) error {
	delete(d.db[programsPath], pid)
	return nil
}

func (d *MockDB) LoadClass(_ context.Context, cid string) (c Class, err error) {
	c, ok := d.db[classesPath][cid].(Class)
	if !ok {
		err = errors.New("invalid class ID")
	}
	return
}

func (d *MockDB) StoreClass(_ context.Context, c Class) error {
	d.db[classesPath][c.CID] = c
	return nil
}

func (d *MockDB) DeleteClass(_ context.Context, cid string) error {
	delete(d.db[classesPath], cid)
	return nil
}

func (d *MockDB) LoadUser(_ context.Context, uid string) (u User, err error) {
	u, ok := d.db[usersPath][uid].(User)
	if !ok {
		err = errors.New("invalid user ID")
	}
	return
}

func (d *MockDB) StoreUser(_ context.Context, u User) error {
	d.db[usersPath][u.UID] = u
	return nil
}

func (d *MockDB) DeleteUser(_ context.Context, uid string) error {
	delete(d.db[usersPath], uid)
	return nil
}

// Creates a new MockDB.
func OpenMock() *MockDB {
	m := MockDB{db: make(map[string]map[string]interface{})}
	m.db[usersPath] = make(map[string]interface{})
	m.db[programsPath] = make(map[string]interface{})
	m.db[classesPath] = make(map[string]interface{})
	return &m
}
