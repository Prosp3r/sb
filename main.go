package main

import (
	"log"
	"math/rand"
	"time"
)

type Shop struct {
	WRSeats      []string  //maximum number of seats in waiting room
	VacantSeats  int       //total free/vacant seats in waiting room
	LockStatus   bool      //true means you can't read seat availability
	LockChan     chan bool //LockChannel reads into LockStatus from Barber and Customers
	BarberStatus string    //Barber status /sleeping/awake/barbing/checkingwr
	// Customers    []Customer       //total customers in wr
}

type Customer struct {
	Number int
	Name   string
}

func main() {
	var ShopLock chan bool
	seats := make([]string, 5)
	BarberShop := Shop{
		WRSeats:      seats,
		VacantSeats:  5,
		LockStatus:   false,
		LockChan:     ShopLock,
		BarberStatus: "awake",
	}
	Customers := makeCustomers()

	go BarberShop.OpenLockShop()

	BarberShop.Run(Customers)

}

//OpenLockShop - Always receive the channel data, prevents it from blocking and empty it's pipe into LockStatus.
func (S Shop) OpenLockShop() {
	for {
		select {
		case <-S.LockChan:
			S.LockStatus = <-S.LockChan
		}
	}
}

func (S Shop) RunBarber() {
	BarbTime := 10 //Barbing delay time
	sleepCounter := 0
	for {

		if S.BarberStatus != "sleeping" {
			sleepCounter = 0 //reset sleep counter after non-sleep event

			if S.BarberStatus == "awake" && !S.LockStatus {
				S.LockChan <- true            //lock the waiting room
				S.BarberStatus = "checkingwr" //checking waiting room
				if S.VacantSeats < len(S.WRSeats) {
					client := S.WRSeats[len(S.WRSeats)-1]
					S.WRSeats = S.WRSeats[:len(S.WRSeats)-1] //truncate
					S.BarberStatus = "barbing"
					log.Printf("Picked up a client. Name: %v. Going barbing", client)
					S.LockChan <- false
					//barbing delay
					time.Sleep(time.Duration(BarbTime) * time.Millisecond) //barber is barbing
					S.BarberStatus = "awake"                               //done barbing
					log.Printf("Done barbing client. Name: %v.", client)
				}
			}

			if S.BarberStatus == "awake" && S.LockStatus {
				//Awake but waiting room is in use
				time.Sleep(time.Duration(3) * time.Millisecond)
			}

		} else {
			//sleeping
			//sest more if sleep try times increase
			time.Sleep(time.Duration(sleepCounter) * time.Microsecond) //sleep more if sleep tries are more
			sleepCounter++
			continue
		}
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

//Run -
func (S Shop) Run(C []Customer) {
	for _, Cu := range C {

		go func(Cu Customer) {
			if !S.LockStatus {
				S.LockChan <- true
				//Barber is not sleeping and
				if S.VacantSeats > 0 && S.BarberStatus != "sleeping" {
					S.WRSeats = append(S.WRSeats, Cu.Name)
					S.VacantSeats--
					S.LockChan <- false
					log.Printf("Name is %v and I'm Taking a seat. %v seats are left", Cu.Name, S.VacantSeats)
				}

				//Barber is sleeping and there's empty seats- take a seat and - wake 'em up
				if S.VacantSeats > 0 && S.BarberStatus == "sleeping" {
					S.BarberStatus = "awake" //wake barber up
					S.WRSeats = append(S.WRSeats, Cu.Name)
					S.VacantSeats--
					S.LockChan <- false
					log.Printf("Name is %v and I'm Taking a seat. woke the barber up...oh and %v seats are left", Cu.Name, S.VacantSeats)
				}

				//Barber is not sleeping but seats are all occupied
				if S.VacantSeats < 1 {
					S.LockChan <- false
					log.Printf("Name is %v and I'm walking off. %v seats are left", Cu.Name, S.VacantSeats)
				}
			}
		}(Cu)
	}
}

func makeCustomers() []Customer {
	Customers := []Customer{
		{Number: 0, Name: "John"}, {Number: 1, Name: "Jane"}, {Number: 2, Name: "Google"}, {Number: 3, Name: "Pringle"}, {Number: 4, Name: "Mr.Mouse"}, {Number: 5, Name: "Viv"}, {Number: 6, Name: "Walid"},
	}
	return Customers
}
