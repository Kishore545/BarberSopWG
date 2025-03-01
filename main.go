//With Goroutine channel
//-------------------------

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	haircutTime     = 2 * time.Second // Time taken for a haircut
	maxWaitingSeats = 5               // Maximum number of customers that can wait
)

type BarberShop struct {
	seats      chan struct{}
	customers  chan int
	barberDone chan bool
	wg         sync.WaitGroup
}

func NewBarberShop() *BarberShop {
	return &BarberShop{
		seats:      make(chan struct{}, maxWaitingSeats),
		customers:  make(chan int),
		barberDone: make(chan bool),
	}
}

func (shop *BarberShop) Barber() {
	for customer := range shop.customers {
		fmt.Printf("Barber is cutting hair for Customer %d\n", customer)
		time.Sleep(haircutTime)
		fmt.Printf("Customer %d is done with the haircut\n", customer)
		<-shop.seats // Remove customer from waiting seats
		shop.barberDone <- true
	}
}

func (shop *BarberShop) Customer(id int) {
	shop.wg.Add(1)
	defer shop.wg.Done()

	select {
	case shop.seats <- struct{}{}: // Check if there's space in the waiting area
		fmt.Printf("Customer %d is waiting.\n", id)
		shop.customers <- id
		<-shop.barberDone // Wait until the barber is done
	default:
		fmt.Printf("Customer %d leaves as the shop is full.\n", id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	shop := NewBarberShop()
	go shop.Barber()

	for i := 1; i <= 20; i++ { // Simulating 20 customers
		go shop.Customer(i)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}

	shop.wg.Wait()        // Wait for all customers to finish
	close(shop.customers) // Close channels
	close(shop.barberDone)
	fmt.Println("Barber shop is closing.")
}
