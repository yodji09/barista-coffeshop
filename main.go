package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Order struct {
	ID           int
	Size         string
	CoffeeFlavor string
	CoffeeType   string
}

type ResultOrder struct {
	ID             int
	OrderCompleted bool
	BaristaNo      int
	OrderFlavor    string
}

var AvailableSize = []string{"Grande", "Medium", "Small"}
var AvailableFlavor = []string{"Cappucino", "Frappucino", "Latte", "Long Black", "Americano", "Matcha Latte"}
var AvailableType = []string{"Hot", "Ice"}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./barista-coffeshop <integer>")
		fmt.Println("Example: ./barista-coffeshop 200")
		os.Exit(1)
	}
	val, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("argument must be integer")
		os.Exit(1)
	}
	if val < 1 {
		fmt.Println("argument must be greater than 1")
		os.Exit(1)
	}
	fmt.Println("argument receive, initiated order id ", val)
	var wg sync.WaitGroup
	jobs := make(chan *Order, 100)
	resChan := make(chan *ResultOrder, 100)
	totalBarista := 5
	for i := 0; i < totalBarista; i++ {
		go DispatchWorker(&wg, i, jobs, resChan)
	}
	var completedOrder atomic.Int32
	var failedOrder atomic.Int32
	go func() {
		for i := 1; i <= val; i++ {
			wg.Add(1)
			orderData := generateOrder(i)
			jobs <- &orderData
		}
		close(jobs)
	}()

	for res := range resChan {
		if res.OrderCompleted {
			fmt.Printf("order completed for ID %v Barista No %v And Coffe Flavor %v \n", res.ID, res.BaristaNo, res.OrderFlavor)
			completedOrder.Add(1)
		} else {
			fmt.Printf("Order Completed but the customer already left because take too long for ID %v, Barista No %v and Coffee Flavor %v \n", res.ID, res.BaristaNo, res.OrderFlavor)
			failedOrder.Add(1)
		}
	}
	wg.Wait()
	close(resChan)
	fmt.Println("total completed order ", completedOrder.Load())
	fmt.Println("total failed order", failedOrder.Load())
}

func generateOrder(id int) Order {
	size := AvailableSize[rand.Intn(len(AvailableSize))]
	orderedType := AvailableType[rand.Intn(len(AvailableType))]
	orderedFlavor := AvailableFlavor[rand.Intn(len(AvailableFlavor))]
	ordered := Order{Size: size, CoffeeFlavor: orderedFlavor, CoffeeType: orderedType, ID: id}
	return ordered
}

func DispatchWorker(wg *sync.WaitGroup, baristaNo int, jobs <-chan *Order, result chan<- *ResultOrder) {
	for job := range jobs {
		fmt.Printf("order received for ID %v barista No %v and Coffe Flavor %v", job.ID, baristaNo, job.CoffeeFlavor)
		resultOrder := proceedJob(job)
		resData := new(ResultOrder)
		resData.ID = job.ID
		resData.BaristaNo = baristaNo
		resData.OrderCompleted = resultOrder
		resData.OrderFlavor = job.CoffeeFlavor
		result <- resData
		wg.Done()
	}
}

func proceedJob(order *Order) bool {
	fmt.Printf("order processed for %v", order.CoffeeFlavor)
	rand.Seed(time.Now().UnixNano())
	randomSecond := rand.Intn(10)
	time.Sleep(time.Duration(randomSecond) * time.Millisecond)
	if randomSecond > 5 {
		return false
	}
	return true
}
