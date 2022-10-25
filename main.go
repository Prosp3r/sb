package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type WaitRoom struct {
	MaxSeats   int
	UsedSeats  map[int]string
	LockStatus string //open/closed
	Keyholder  string
}

func (Wr *WaitRoom) Unlock(Keyholder string) bool {
	if Wr.Keyholder == Keyholder {
		Wr.LockStatus = "open"
		Wr.Keyholder = ""
		return true
	}
	return false
}

func (Wr *WaitRoom) Lock(Keyholder string) bool {
	if Wr.Keyholder == "" {
		Wr.LockStatus = "closed"
		Wr.Keyholder = Keyholder
		return true
	}
	return false
}

func (Wr *WaitRoom) SeatStatus() map[int]string {
	log.Printf("Waiting room seat status: %v\n", Wr.UsedSeats)
	return Wr.UsedSeats
}

func (Wr *WaitRoom) TakeASeat(Cu Customer) bool {
	CustomerID := <-Cu
	CustomerIDString := strconv.Itoa(CustomerID)
	if len(Wr.UsedSeats) < Wr.MaxSeats {
		Wr.UsedSeats[CustomerID] = CustomerIDString
		log.Printf("Waiting room seat taken by Customer -> %v\n", CustomerID)
		_ = Wr.SeatStatus()
		return true
	}
	return false
}

//TakeCustomerFromWr -
func (Wr *WaitRoom) TakeCustomerFromWr(Brb *Barber, Cu Customer) bool {
	Bd := len(Wr.UsedSeats)
	delete(Wr.UsedSeats, <-Cu)

	if len(Wr.UsedSeats) < Bd {
		return true
	}
	return false
}

type Barber struct {
	Status            string //sleeping/awake/checkingWR/barbing
	PermittedStatuses map[string]string
}

//SetState - wake/sleep/checking/barbing
func (Brb *Barber) SetStatus(state string) bool {
	if v, ok := Brb.PermittedStatuses[state]; ok {
		Brb.Status = v
		return true
	}
	return false
}

//Sleep - Sets the Barber's state to sleeping
func (Brb *Barber) Sleep() bool {
	if Brb.Status != "barbing" && Brb.Status != "checkingWR" {
		return Brb.SetStatus("sleeping")
	}
	return false
}

func (Brb *Barber) CheckWR() bool {
	return false
}

type Shop chan WaitRoom

func (S *Shop) Run() bool {

	return false
}

func (S *Shop) FreeSeats() int {
	// Sh := S

	return 0
}

type Customer chan int

func (Cu Customer) WalkIn(End chan bool, Brb *Barber, Wr *WaitRoom) {
	for {
		select {
		case <-End:
			fmt.Printf("End Signaled closing doors")
			return
		default:
			CustomerID := <-Cu
			fmt.Printf("New Customer at the door ID: %v\n", CustomerID)
			for Tries := 1; ; Tries++ {
				if Wr.LockStatus == "locked" && Wr.Keyholder == "barber" {
					fmt.Printf("Waiting room is locked by barber try No: %v. I %v will try again\n", Tries, CustomerID)
					continue
				}

				if Wr.LockStatus == "open" {
					Wr.Lock("customer")
					if Brb.Status == "sleeping" {
						Brb.SetStatus("wake") //wake barber
					}
					//if seats are vacant, take a seat
					if len(Wr.UsedSeats) < Wr.MaxSeats {
						//take a seat
						if !Wr.TakeASeat(Cu) {
							log.Printf("Could not take waiting room seat... customer %v is leaving\n", CustomerID)
						}
						Wr.Unlock("customer")
					} else {
						//walk away
						fmt.Printf("Waiting room is full try No: %v. I, %v will walk away now\n", Tries, CustomerID)
						Wr.Unlock("customer")
						time.Sleep(time.Millisecond * 2000)
						continue
					}
				}
			}
		}
	}
}

//NewCustomerGen - Generates new customers randomly conitnuously
func NewCustomerGen(Ch chan<- int, End chan bool, limit int) int {
	var TotalSent int = 0
	if limit > 0 {
		for i := 1; i < limit+1; i++ {
			Ch <- i
			TotalSent = TotalSent + 1
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}

		fmt.Printf("No more customers, closing doors. sent total: %v\n", TotalSent)
		close(Ch)
		close(End)

	} else {
		for i := 1; ; i++ {
			select {
			case End <- true:
				fmt.Println("Stop allowing more customers please closing doors")
				close(Ch)
				close(End)
				return TotalSent

			default:
				Ch <- i
				TotalSent = TotalSent + 1
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			}
		}
	}
	return TotalSent
}

func main() {
	BarberStates := map[string]string{"sleeping": "sleeping", "awake": "awake", "checkingWR": "checkingWR", "barbing": "barbing"}
	var Cu Customer = make(chan int)
	var End = make(chan bool)
	var limit int = 10 //limitless

	seats := make(map[int]string, 10)

	Wr := WaitRoom{
		MaxSeats:   10,
		UsedSeats:  seats,
		LockStatus: "open",
		Keyholder:  "",
	}
	Brb := Barber{
		Status:            "awake", //awake/checkingwr/barbing/sleeping
		PermittedStatuses: BarberStates,
	}

	go NewCustomerGen(Cu, End, limit)
	go Cu.WalkIn(End, &Brb, &Wr)
	// time.Sleep(time.Duration(rand.Intn(1e3)) * time.Second)
	time.Sleep(15 * time.Second)
	fmt.Println("Leaving...")
}
