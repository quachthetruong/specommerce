package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CreateOrderRequest struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	TotalAmount  float64 `json:"total_amount"`
}

// 200 unique real customer names
var customers = []struct {
	ID   string
	Name string
}{
	{"CUST001", "Alexander Johnson"}, {"CUST002", "Sophia Williams"}, {"CUST003", "Benjamin Brown"}, {"CUST004", "Isabella Davis"},
	{"CUST005", "Lucas Miller"}, {"CUST006", "Mia Wilson"}, {"CUST007", "Henry Moore"}, {"CUST008", "Charlotte Taylor"},
	{"CUST009", "Sebastian Anderson"}, {"CUST010", "Amelia Thomas"}, {"CUST011", "Oliver Jackson"}, {"CUST012", "Harper White"},
	{"CUST013", "Theodore Harris"}, {"CUST014", "Evelyn Martin"}, {"CUST015", "Jack Thompson"}, {"CUST016", "Abigail Garcia"},
	{"CUST017", "Owen Martinez"}, {"CUST018", "Emily Robinson"}, {"CUST019", "Levi Clark"}, {"CUST020", "Elizabeth Rodriguez"},
	{"CUST021", "Daniel Lewis"}, {"CUST022", "Sofia Lee"}, {"CUST023", "Samuel Walker"}, {"CUST024", "Avery Hall"},
	{"CUST025", "Matthew Allen"}, {"CUST026", "Ella Young"}, {"CUST027", "Jackson Hernandez"}, {"CUST028", "Scarlett King"},
	{"CUST029", "David Wright"}, {"CUST030", "Victoria Lopez"}, {"CUST031", "Joseph Hill"}, {"CUST032", "Grace Scott"},
	{"CUST033", "Carter Green"}, {"CUST034", "Chloe Adams"}, {"CUST035", "Wyatt Baker"}, {"CUST036", "Camila Gonzalez"},
	{"CUST037", "Julian Nelson"}, {"CUST038", "Aria Carter"}, {"CUST039", "Luke Mitchell"}, {"CUST040", "Madison Perez"},
	{"CUST041", "Grayson Roberts"}, {"CUST042", "Layla Turner"}, {"CUST043", "Liam Phillips"}, {"CUST044", "Penelope Campbell"},
	{"CUST045", "Isaac Parker"}, {"CUST046", "Riley Evans"}, {"CUST047", "Jayden Edwards"}, {"CUST048", "Nora Collins"},
	{"CUST049", "Connor Stewart"}, {"CUST050", "Hazel Sanchez"}, {"CUST051", "Aaron Morris"}, {"CUST052", "Luna Rogers"},
	{"CUST053", "Eli Reed"}, {"CUST054", "Aurora Cook"}, {"CUST055", "Nathan Morgan"}, {"CUST056", "Savannah Bell"},
	{"CUST057", "Caleb Murphy"}, {"CUST058", "Brooklyn Bailey"}, {"CUST059", "Ryan Rivera"}, {"CUST060", "Lillian Cooper"},
	{"CUST061", "Asher Richardson"}, {"CUST062", "Samantha Cox"}, {"CUST063", "Thomas Howard"}, {"CUST064", "Leah Ward"},
	{"CUST065", "Leo Torres"}, {"CUST066", "Audrey Peterson"}, {"CUST067", "Charles Gray"}, {"CUST068", "Bella Ramirez"},
	{"CUST069", "Christopher James"}, {"CUST070", "Claire Watson"}, {"CUST071", "Joshua Brooks"}, {"CUST072", "Addison Kelly"},
	{"CUST073", "Andrew Sanders"}, {"CUST074", "Natalie Price"}, {"CUST075", "Ezra Bennett"}, {"CUST076", "Paisley Wood"},
	{"CUST077", "John Barnes"}, {"CUST078", "Naomi Ross"}, {"CUST079", "Hudson Henderson"}, {"CUST080", "Elena Coleman"},
	{"CUST081", "Christian Jenkins"}, {"CUST082", "Sarah Perry"}, {"CUST083", "Jaxon Powell"}, {"CUST084", "Anna Long"},
	{"CUST085", "Jameson Patterson"}, {"CUST086", "Quinn Hughes"}, {"CUST087", "Cooper Flores"}, {"CUST088", "Nevaeh Washington"},
	{"CUST089", "Jeremiah Butler"}, {"CUST090", "Maya Simmons"}, {"CUST091", "Easton Foster"}, {"CUST092", "Willow Gonzales"},
	{"CUST093", "Nolan Bryant"}, {"CUST094", "Kinsley Alexander"}, {"CUST095", "Adrian Russell"}, {"CUST096", "Aaliyah Griffin"},
	{"CUST097", "Cameron Diaz"}, {"CUST098", "Genesis Hayes"}, {"CUST099", "Jordan Myers"}, {"CUST100", "Ariana Ford"},
	{"CUST101", "Ian Hamilton"}, {"CUST102", "Allison Graham"}, {"CUST103", "Carson Sullivan"}, {"CUST104", "Gabriella Wallace"},
	{"CUST105", "Jaxson Stone"}, {"CUST106", "Serenity West"}, {"CUST107", "Colton Cole"}, {"CUST108", "Stella Jordan"},
	{"CUST109", "Brayden Reed"}, {"CUST110", "Delilah Fox"}, {"CUST111", "Robert McDonald"}, {"CUST112", "Isla Gardner"},
	{"CUST113", "Greyson Lawrence"}, {"CUST114", "Ivy Fuller"}, {"CUST115", "Rowan Ellis"}, {"CUST116", "Autumn Boyd"},
	{"CUST117", "Adam Carpenter"}, {"CUST118", "Piper Mason"}, {"CUST119", "Bentley Moreno"}, {"CUST120", "Lydia Warren"},
	{"CUST121", "Collin Harvey"}, {"CUST122", "Nova Gilbert"}, {"CUST123", "Ryder Graham"}, {"CUST124", "Emilia Simpson"},
	{"CUST125", "Kingston Knight"}, {"CUST126", "Violet Butler"}, {"CUST127", "Hunter Ellis"}, {"CUST128", "Iris Alexander"},
	{"CUST129", "Jose Hoffman"}, {"CUST130", "Adalynn Griffin"}, {"CUST131", "Dominic Mason"}, {"CUST132", "Eden Diaz"},
	{"CUST133", "Jace Butler"}, {"CUST134", "Emery Hayes"}, {"CUST135", "Kai Myers"}, {"CUST136", "Remi Ford"},
	{"CUST137", "Braxton Hamilton"}, {"CUST138", "Mackenzie Graham"}, {"CUST139", "Maximus Sullivan"}, {"CUST140", "Brielle Wallace"},
	{"CUST141", "Silas Stone"}, {"CUST142", "Tessa West"}, {"CUST143", "Maddox Cole"}, {"CUST144", "Juniper Jordan"},
	{"CUST145", "Iwan Reed"}, {"CUST146", "Catalina Fox"}, {"CUST147", "Kaiden McDonald"}, {"CUST148", "Rosalie Gardner"},
	{"CUST149", "Antonio Lawrence"}, {"CUST150", "Sloane Fuller"}, {"CUST151", "Colson Ellis"}, {"CUST152", "Vera Boyd"},
	{"CUST153", "Maverick Carpenter"}, {"CUST154", "Cecilia Mason"}, {"CUST155", "Rhett Moreno"}, {"CUST156", "Margot Warren"},
	{"CUST157", "Knox Harvey"}, {"CUST158", "Anastasia Gilbert"}, {"CUST159", "Beckett Graham"}, {"CUST160", "Daphne Simpson"},
	{"CUST161", "Finn Knight"}, {"CUST162", "Eloise Butler"}, {"CUST163", "Dash Ellis"}, {"CUST164", "Ophelia Alexander"},
	{"CUST165", "Dean Hoffman"}, {"CUST166", "Lennox Griffin"}, {"CUST167", "Holden Mason"}, {"CUST168", "Celeste Diaz"},
	{"CUST169", "Julius Hayes"}, {"CUST170", "Ember Myers"}, {"CUST171", "Creed Ford"}, {"CUST172", "Sage Hamilton"},
	{"CUST173", "Phoenix Graham"}, {"CUST174", "Wren Sullivan"}, {"CUST175", "Zion Wallace"}, {"CUST176", "Thea Stone"},
	{"CUST177", "Archer West"}, {"CUST178", "Azalea Cole"}, {"CUST179", "Orion Jordan"}, {"CUST180", "Magnolia Reed"},
	{"CUST181", "Soren Fox"}, {"CUST182", "Poppy McDonald"}, {"CUST183", "Atticus Gardner"}, {"CUST184", "Jasmine Lawrence"},
	{"CUST185", "Atlas Fuller"}, {"CUST186", "Marigold Ellis"}, {"CUST187", "Felix Boyd"}, {"CUST188", "Clementine Carpenter"},
	{"CUST189", "Jasper Mason"}, {"CUST190", "Dahlia Moreno"}, {"CUST191", "Enzo Warren"}, {"CUST192", "Sage Harvey"},
	{"CUST193", "Cash Gilbert"}, {"CUST194", "Willa Graham"}, {"CUST195", "Axel Simpson"}, {"CUST196", "Indigo Knight"},
	{"CUST197", "Knox Butler"}, {"CUST198", "Meadow Ellis"}, {"CUST199", "Zander Alexander"}, {"CUST200", "River Hoffman"},
}

const (
	orderServiceURL = "http://localhost:8080/api/v1/orders"
	totalOrders     = 500
	ordersBelow200  = 400 // Orders less than $200
	orders200to400  = 100 // Orders from $200 to $400
	maxConcurrency  = 50  // Maximum concurrent requests
)

func createOrder(client *http.Client, order CreateOrderRequest, wg *sync.WaitGroup, orderNum int) {
	defer wg.Done()

	sleepDuration := time.Duration(2000+rand.Intn(8000)) * time.Millisecond
	time.Sleep(sleepDuration)

	jsonData, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", orderServiceURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error creating order #%d for %s: %v\n", orderNum, order.CustomerID, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Order #%d - %d %s: $%.2f (waited %.1fs)\n", orderNum, resp.StatusCode, order.CustomerID, order.TotalAmount, sleepDuration.Seconds())
}

// Step 1: Create two separate lists of transactions with proper ratios
// Step 2: Generate 400 transactions below $200 and 100 transactions between $200 and $400
// Step 3: Shuffle both customer and transaction lists to ensure randomness
// Step 4: Assign the first 200 transactions to the first 200 unique customers
// Step 5: Distribute the remaining 300 transactions randomly among all customers
// Step 6: Combine and shuffle all orders
// Step 7: Process all orders concurrently with random delays
func main() {
	fmt.Println("Simple Order Mock Generator")
	fmt.Printf("Target: %d orders (%d < $200, %d $200-$400)\n", totalOrders, ordersBelow200, orders200to400)
	fmt.Printf("Using %d unique customers\n", len(customers))
	fmt.Printf("Running orders concurrently with 3-10 second random delays\n")
	fmt.Println(strings.Repeat("=", 50))

	client := &http.Client{Timeout: 30 * time.Second}

	semaphore := make(chan struct{}, maxConcurrency)

	var belowTransactions []CreateOrderRequest
	var aboveTransactions []CreateOrderRequest

	for i := 0; i < ordersBelow200; i++ {
		amount := 50.0 + rand.Float64()*149.0
		belowTransactions = append(belowTransactions, CreateOrderRequest{
			TotalAmount: float64(int(amount*100)) / 100,
		})
	}

	for i := 0; i < orders200to400; i++ {
		amount := 200.0 + rand.Float64()*200.0
		aboveTransactions = append(aboveTransactions, CreateOrderRequest{
			TotalAmount: float64(int(amount*100)) / 100,
		})
	}

	allTransactions := append(belowTransactions, aboveTransactions...)

	shuffledCustomers := make([]struct{ ID, Name string }, len(customers))
	copy(shuffledCustomers, customers)
	rand.Shuffle(len(shuffledCustomers), func(i, j int) {
		shuffledCustomers[i], shuffledCustomers[j] = shuffledCustomers[j], shuffledCustomers[i]
	})

	rand.Shuffle(len(allTransactions), func(i, j int) {
		allTransactions[i], allTransactions[j] = allTransactions[j], allTransactions[i]
	})

	var guaranteedOrders []CreateOrderRequest
	for i := 0; i < 200; i++ {
		customer := shuffledCustomers[i]
		transaction := allTransactions[i]
		transaction.CustomerID = customer.ID
		transaction.CustomerName = customer.Name
		guaranteedOrders = append(guaranteedOrders, transaction)
	}

	var remainingOrders []CreateOrderRequest
	for i := 200; i < len(allTransactions); i++ {
		customer := customers[rand.Intn(len(customers))]
		transaction := allTransactions[i]
		transaction.CustomerID = customer.ID
		transaction.CustomerName = customer.Name
		remainingOrders = append(remainingOrders, transaction)
	}

	var allOrders []CreateOrderRequest
	allOrders = append(allOrders, guaranteedOrders...)
	allOrders = append(allOrders, remainingOrders...)

	rand.Shuffle(len(allOrders), func(i, j int) {
		allOrders[i], allOrders[j] = allOrders[j], allOrders[i]
	})

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("Starting %d orders with max %d concurrent requests...\n", len(allOrders), maxConcurrency)

	for i, order := range allOrders {
		wg.Add(1)
		go func(ord CreateOrderRequest, orderNum int) {
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore
			createOrder(client, ord, &wg, orderNum)
		}(order, i+1)
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Completed %d orders in %v\n", totalOrders, duration)
	fmt.Printf("Rate: %.1f orders/second\n", float64(totalOrders)/duration.Seconds())
}
