# barista-coffeshop

Write a Go program that simulates the coffee shop scenario using goroutines. The program should:

1. Generate random coffee orders (e.g., latte, cappuccino, espresso) at random intervals.
2. Assign each order to a barista (represented by a goroutine).
3. Each barista should prepare the coffee order and print a message when the order is ready.
4. Use channels to communicate between the main goroutine and the barista goroutines.

*Requirements:*

- Use at least 3-5 baristas (goroutines) to handle orders concurrently.
- Use a channel to send orders to the baristas.
- Use a WaitGroup to ensure that the main goroutine waits for all orders to be completed before exiting.
- Print messages to indicate when an order is received, when a barista starts preparing an order, and when an order is ready.

*Example Output:*

```
Order received: Latte
Barista 1: Preparing Latte
Order received: Cappuccino
Barista 2: Preparing Cappuccino
Barista 1: Latte ready!
Order received: Espresso
Barista 3: Preparing Espresso
Barista 2: Cappuccino ready!
Barista 3: Espresso ready!
```
""""