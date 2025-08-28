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
	fmt.Println("argument receive, initiated order for total of ", val)

	var wg sync.WaitGroup
	jobs := make(chan *Order, 100)
	resChan := make(chan *ResultOrder, 100)
	totalBarista := 5

	var completedOrder atomic.Int32
	var failedOrder atomic.Int32

	// Start barista workers
	for i := 1; i <= totalBarista; i++ {
		go DispatchWorker(&wg, i, jobs, resChan)
	}

	// Start order generator
	go func() {
		for i := 1; i <= val; i++ {
			wg.Add(1)
			orderData := generateOrder(i)
			jobs <- &orderData
		}
		close(jobs)
	}()

	// Start result processor
	go func() {
		for res := range resChan {
			fmt.Println("receiving result for id ", res.ID)
			if res.OrderCompleted {
				fmt.Printf("order completed for ID %v Barista No %v And Coffe Flavor %v \n", res.ID, res.BaristaNo, res.OrderFlavor)
				completedOrder.Add(1)
			} else {
				fmt.Printf("Order Completed but the customer already left because take too long for ID %v, Barista No %v and Coffee Flavor %v \n", res.ID, res.BaristaNo, res.OrderFlavor)
				failedOrder.Add(1)
			}
		}
	}()

	// Wait for all orders to be processed
	fmt.Println("waiting all worker to be done")
	wg.Wait()
	fmt.Println("closing result channel")
	// Close result channel after all workers are done
	close(resChan)

	// Give a moment for the result processor to finish
	time.Sleep(100 * time.Millisecond)

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
		fmt.Printf("Barista NO %v receiving order for ID %v and Coffe Flavor %v\n", baristaNo, job.ID, job.CoffeeFlavor)
		resultOrder := proceedJob(job, baristaNo)
		resData := new(ResultOrder)
		resData.ID = job.ID
		resData.BaristaNo = baristaNo
		resData.OrderCompleted = resultOrder
		resData.OrderFlavor = job.CoffeeFlavor
		fmt.Printf("BARISTA NO %v sending order ID %v to frontline \n", baristaNo, job.ID)
		result <- resData
		fmt.Printf("BARISTA NO %v sending signal done for order ID %v, and ready taking another order \n", baristaNo, job.ID)
		wg.Done()
	}
}

func proceedJob(order *Order, baristaNo int) bool {
	fmt.Printf("BARISTA NO %v processing order for ID %v and Flavor %v\n", baristaNo, order.ID, order.CoffeeFlavor)
	rand.Seed(time.Now().UnixNano())
	randomSecond := rand.Intn(15)
	time.Sleep(time.Duration(randomSecond) * time.Second)
	fmt.Printf("BARISTA NO %v processed order for id %v with duration of %v S \n", baristaNo, order.ID, randomSecond)
	if randomSecond > 5 {
		return false
	}
	return true
}
