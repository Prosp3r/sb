package main

import (
	"log"
	"math/rand"
	"time"
)

type Customer struct {
	Number int
	Name   string
}

type Shop struct {
	MaxSeats      int
	FreeSeats     int
	SeatCounter   chan int
	LockStatus    bool //true means you can't read seat availability
	CurrentLocker string
}

type OpenLocker struct {
	Sendertype string
	Name       string
	LockChan   bool
}

var Customers = make(chan Customer)

func main() {
	Sc := make(chan int)
	go makeCustomers()

	S := Shop{
		MaxSeats:    10,
		FreeSeats:   10,
		SeatCounter: Sc,
		LockStatus:  true,
	}

	Ol := make(chan OpenLocker)
	go RunShop(S, Ol)
	go RunBarber(Ol)

	for i := 0; ; i++ {
		Olock := <-Ol
		log.Printf("Gotten value from = %v Name: %v - - Free Seats : %v\n", Olock.Sendertype, Olock.Name, Olock.LockChan)

		if Olock.LockChan == true && Olock.Sendertype == S.CurrentLocker {
			S.LockStatus = Olock.LockChan
			log.Printf("Just locked up the Shop WR : -- -%v\n", Olock.Sendertype)
		} else {
			//do nothing
			log.Printf("Could not lock up the Shop WR : -- -%v\n", Olock.Sendertype)
		}
	}
}

func RunBarber(Ol chan OpenLocker) {
	//barbing delay +3
	//awake waiting delay +2
	//checkingwr delay +3

	//barber status /awake/sleeping/checkingwr/barbing
	for {
		sender := "barber"
		BName := "barberD"
		Openl := OpenLocker{Sendertype: sender, Name: BName}
		Openl.LockChan = true
		Ol <- Openl
		//awake
		time.Sleep(2 * time.Second)

		//check for shop lock status
		//Change status to barbing

		//Change status to sleeping
		//Change status to checkingwr
		//Change status to awake

	}
}

func RunShop(S Shop, Ol chan OpenLocker) {
	sender := "customer"
	var C Customer

	for i := 0; ; i++ {
		C = <-Customers
		if C.Number%2 == 0 {

			OpenL := OpenLocker{Sendertype: sender, Name: C.Name}
			OpenL.LockChan = false
			Ol <- OpenL
		} else {
			OpenL := OpenLocker{Sendertype: sender, Name: C.Name}
			OpenL.LockChan = true
			Ol <- OpenL
		}
		time.Sleep(1 * time.Second)
	}
}

func makeCustomers() {
	for i := 0; ; i++ {
		Customers <- Customer{Number: i, Name: RandomCode(6)}
	}
}

func RandomCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// package main

// import (
// 	"log"
// 	"math/rand"
// 	"sync"
// 	"time"
// )

// type Shop struct {
// 	WRSeats      []string  //maximum number of seats in waiting room
// 	VacantSeats  int       //total free/vacant seats in waiting room
// 	LockStatus   bool      //true means you can't read seat availability
// 	LockChan     chan bool //LockChannel reads into LockStatus from Barber and Customers
// 	BarberStatus string    //Barber status /sleeping/awake/barbing/checkingwr
// 	// Customers    []Customer       //total customers in wr
// }

// type Customer struct {
// 	Number int
// 	Name   string
// }

// func main() {
// 	var wg sync.WaitGroup

// 	var ShopLock chan bool
// 	seats := make([]string, 5)
// 	BarberShop := Shop{
// 		WRSeats:      seats,
// 		VacantSeats:  5,
// 		LockStatus:   false,
// 		LockChan:     ShopLock,
// 		BarberStatus: "awake",
// 	}

// 	Customers := makeCustomers()
// 	log.Print("Starting...")

// 	go RunBarber(BarberShop)
// 	// go OpenLockShop(BarberShop)
// 	wg.Add(1)
// 	go Run(BarberShop, Customers, &wg)
// 	go func() {
// 		for i := 0; ; i++ {
// 			f1 := <-BarberShop.LockChan
// 			BarberShop.LockStatus = f1
// 			log.Printf("Changing WR Lock Status from...%v", BarberShop.LockStatus)
// 			time.Sleep(10 * time.Millisecond)
// 		}

// 	}()

// 	time.Sleep(5 * time.Second)
// 	log.Println("Leaving program....bored")
// 	wg.Wait()
// }

// //Run -
// func Run(S Shop, C []Customer, wg *sync.WaitGroup) {
// 	log.Println("Starting Customer Run...")
// 	defer wg.Done()
// 	for _, Cu := range C {

// 		go func(Cu Customer) {

// 			if !S.LockStatus {
// 				S.LockChan <- true
// 				log.Printf("Customer...%v...LockWRStatus : %v", Cu.Name, S.LockStatus)
// 				//Barber is not sleeping and
// 				if S.VacantSeats > 0 && S.BarberStatus != "sleeping" {
// 					S.WRSeats = append(S.WRSeats, Cu.Name)
// 					S.VacantSeats--
// 					S.LockChan <- false
// 					log.Printf("Name is %v and I'm Taking a seat. %v seats are left", Cu.Name, S.VacantSeats)
// 				}

// 				//Barber is sleeping and there's empty seats- take a seat and - wake 'em up
// 				if S.VacantSeats > 0 && S.BarberStatus == "sleeping" {
// 					S.BarberStatus = "awake" //wake barber up
// 					S.WRSeats = append(S.WRSeats, Cu.Name)
// 					S.VacantSeats--
// 					S.LockChan <- false
// 					log.Printf("Name is %v and I'm Taking a seat. woke the barber up...oh and %v seats are left", Cu.Name, S.VacantSeats)
// 				}

// 				//Barber is not sleeping but seats are all occupied
// 				if S.VacantSeats < 1 {
// 					S.LockChan <- false
// 					log.Printf("Name is %v and I'm walking off. %v seats are left", Cu.Name, S.VacantSeats)
// 				}
// 			}
// 		}(Cu)
// 	}
// }

// //OpenLockShop - Always receive the channel data, prevents it from blocking and empty it's pipe into LockStatus.
// func OpenLockShop(S Shop) {
// 	log.Println("Starting OpenLocker...")
// 	for {
// 		S.LockStatus = <-S.LockChan
// 		log.Printf("Changing WR Lock Status from...%v", S.LockStatus)
// 		time.Sleep(1000 * time.Millisecond)
// 	}
// }

// func RunBarber(S Shop) {
// 	S.LockChan <- true
// 	log.Printf("Starting Barber...%v", S.BarberStatus)
// 	log.Printf("Starting Barber...LOCKSTATUS %v", S.LockStatus)
// 	BarbTime := 10 //Barbing delay time
// 	sleepCounter := 0

// 	for {

// 		if S.BarberStatus != "sleeping" {
// 			log.Printf("Starting Barber...%v --%v", S.BarberStatus, S.LockStatus)
// 			sleepCounter = 0 //reset sleep counter after non-sleep event

// 			if S.BarberStatus == "awake" && !S.LockStatus {
// 				S.LockChan <- true //lock the waiting room
// 				log.Printf("Starting Barber...%v --%v", S.BarberStatus, S.LockStatus)
// 				S.BarberStatus = "checkingwr" //checking waiting room
// 				if S.VacantSeats < len(S.WRSeats) {
// 					client := S.WRSeats[len(S.WRSeats)-1]
// 					S.WRSeats = S.WRSeats[:len(S.WRSeats)-1] //truncate
// 					S.BarberStatus = "barbing"
// 					log.Printf("Picked up a client. Name: %v. Going barbing", client)
// 					S.LockChan <- false
// 					//barbing delay
// 					time.Sleep(time.Duration(BarbTime) * time.Millisecond) //barber is barbing
// 					S.BarberStatus = "awake"                               //done barbing
// 					log.Printf("Done barbing client. Name: %v.", client)
// 				}
// 			}

// 			if S.BarberStatus == "awake" && S.LockStatus {
// 				//Awake but waiting room is in use
// 				time.Sleep(time.Duration(3) * time.Millisecond)
// 			}

// 		} else {
// 			//sleeping
// 			log.Println("Starting Barber...Sleeping")
// 			//sest more if sleep try times increase
// 			time.Sleep(time.Duration(sleepCounter) * time.Microsecond) //sleep more if sleep tries are more
// 			sleepCounter++
// 			continue
// 		}
// 		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
// 	}

// }
