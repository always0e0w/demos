package main

import "testing"

func TestNewLRU(t *testing.T) {
	lru := NewLRU(3)
	if lru == nil {
		t.Error("NewLRU failed!")
		t.FailNow()
	}
	t.Logf("lru: %+v", *lru)
}

func TestGet(t *testing.T) {
	lru := NewLRU(3)
	v, e := lru.Get("k1")
	if e || v != "" {
		t.Errorf("k1 not exists, but got %s\n", v)
		t.FailNow()
	}
}

func TestSet(t *testing.T) {
	lru := NewLRU(3)
	lru.Set("k1", "v1")
	v, e := lru.Get("k1")
	if !e || v != "v1" {
		t.Errorf("do not got the expedted value. got %s\n", v)
		t.FailNow()
	}
	lru.Set("k2", "v2")
	lru.Set("k3", "v3")
	lru.Set("k4", "v4")
	v, e = lru.Get("k1")
	if e || v != "" {
		t.Errorf("k1 should not exists. got %s\n", v)
		t.FailNow()
	}
	v, e = lru.Get("k4")
	if !e || v != "v4" {
		t.Errorf("do not got the expedted value. got %s\n", v)
		t.FailNow()
	}
}
