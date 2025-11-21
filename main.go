package main
import ("fmt"
	"strings"
	"os"
	"github.com/joho/godotenv"
	"Transaction_system/requests"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	//load .env file
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}

	//Get api key
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return fmt.Errorf("API key is not set")
	}

	fmt.Println("API key loaded successfully")

        var num string
	var amt int
	var r string
	url1 := "https://demo.campay.net/api/collect/"

for{
        fmt.Printf("Please enter your mobile money number: ")
        fmt.Scan(&num)

	momoNum := make([]int, len(num))
	for i, ch := range num {
		momoNum[i] = int(ch - '0')
	}

        if len(num) == 9 && momoNum[0] == 6 && (numValidator(momoNum) == "MTN" || numValidator(momoNum) == "Orange"){
			break
	}else {
                fmt.Println("invalid number.\nPlease enter a 9 digit number.\nNumber must start with 6")
        }

        fmt.Println(numValidator(momoNum))

}

        printDivider()

        fmt.Printf("\nPlease enter amount: ")
        fmt.Scanln(&amt)

        printDivider()

        if (amt > 0) {
        	fmt.Printf("\nReference: ")
		fmt.Scan(&r)

	printDivider()

		if r != ""{
		        body := map[string]string{ "from" : "237" + num,
                                      		   "amount" : fmt.Sprintf("%d", amt), 
						   "description": r,
			}

		fmt.Println("Sending payment request....")
        	reference := requests.PaymentRequest(url1 , apiKey, body)

		fmt.Println()
		printDivider()

		if reference != "" {
			url2 := "https://demo.campay.net/api/transaction/"+ reference +"/"
			requests.GetTransactionStatus(url2, apiKey)
		}else {
			fmt.Println("Could not retrieve reference")
		}
                printDivider()

		}
        }else {
		fmt.Println("amount must be greater than zero")
	}

	return nil
}


func numValidator(num []int) string {
	n := num[2]
	//mtn number validator
	if num[1] == 5 {
	switch n {
	case 0, 1, 2, 3 , 4:
		return  "MTN"
	case 5, 6, 7, 8:
		return "Orange"
	}
	}else if num[1] == 7{
        switch n {
        case 0, 1, 2, 3 , 4, 5, 6, 7, 8, 9:
                return  "MTN"
	}
	}else if num[1] == 8{
        switch n {
        case 0, 1, 2, 3:
                return  "MTN"
        }
	}else if num[1] == 9{
        switch n {
        case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
                return  "Orange"
        }
	}

	return "Please enter a valid mtn or orange number\n mtn valid number(1st 3 digits)\n - 65(0,...,4)\n - 67(0,...,9)\n - 68(0,..,3)\n Orange valid number(1st 3 digit)\n - 65(5,...,8)\n - 69(0,..,9)\n"
}

func printDivider() {
	width := 60
	fmt.Println(strings.Repeat("-",width))
}
