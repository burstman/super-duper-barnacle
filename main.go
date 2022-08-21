package main

import (
	"fmt"
	csvdatabase "github/bustman/shops/csvDatabase"
	datastorage "github/bustman/shops/dataStorage"
	"sync"
)

type shops struct {
	mu   sync.RWMutex
	shop map[string]float64
	datastorage.Datastore
}

func (s *shops) Get(key string) (*float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.shop == nil {
		return nil, false
	}
	price, ok := s.shop[key]
	return &price, ok
}
func (s *shops) Put(key string, price float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.shop == nil {
		s.shop = make(map[string]float64)
		s.shop[key] = price
		return
	}
	s.shop[key] = price
}

// Load the data named
func (s *shops) Load(name string) {
	s.shop = make(map[string]float64)
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.Datastore.Load(name)
	s.shop = s.Data.(map[string]float64)
}

// Save the data named
func (s *shops) Save(name string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.Data = s.shop
	s.Datastore.Save(name)
	s.shop = nil
}
func (s *shops) List() string {
	return s.Storage.List()
}

// shop constructure
func Newshop(d datastorage.Datastorage) *shops {
	return &shops{shop: map[string]float64{},
		Datastore: datastorage.Datastore{Storage: d},
	}
}
func main() {
	csvdb := csvdatabase.NewCsvData()
	shop := Newshop(csvdb)
	shop.Put("hat", 45)
	shop.Put("shoes", 25)
	shop.Put("tshirt", 15)
	shop.Save("My_first_shop")
	shop.Put("trico", 20)
	shop.Put("chair", 200)
	shop.Save("My_Second_Shop")
	 shop.Load("My_first_shop")
	// shop.Load("My_Second_Shop")
	fmt.Println(shop.List())
	// shop.Load("My_Second_Shop")
	// fmt.Println(shop.shop)

}
