![PayPal Developer Cover](https://github.com/paypaldev/.github/blob/main/pp-cover.png)
<div align="center">
  <a href="https://twitter.com/paypaldev" target="_blank">
    <img alt="Twitter: PayPal Developer" src="https://img.shields.io/twitter/follow/paypaldev?style=social" />
  </a>
  <br />
  <a href="https://twitter.com/paypaldev" target="_blank">Twitter</a>
    <span>&nbsp;&nbsp;-&nbsp;&nbsp;</span>
  <a href="https://www.paypal.com/us/home" target="_blank">PayPal</a>
    <span>&nbsp;&nbsp;-&nbsp;&nbsp;</span>
  <a href="https://developer.paypal.com/home" target="_blank">Docs</a>
    <span>&nbsp;&nbsp;-&nbsp;&nbsp;</span>
  <a href="https://github.com/paypaldev" target="_blank">Code Samples</a>
    <span>&nbsp;&nbsp;-&nbsp;&nbsp;</span>
  <a href="https://dev.to/paypaldeveloper" target="_blank">Blog</a>
  <br />
  <hr />
</div>

# PayPal GO FullStack Standard Checkout
This sample app shows you how to integrate PayPal into your app for the [standard checkout(https://developer.paypal.com/docs/checkout/standard/integrate/) workflow.

To create this application from scratch, follow the [Standard Checkout integration](https://developer.paypal.com/docs/checkout/standard/integrate/) guide from the [PayPal Developer](https://developer.paypal.com/home) docs.

## Public/index.html

### Create Order
The `createOrder()` callback allows you to create the request of your order with the following properties from the [V2 orders create call](https://developer.paypal.com/api/orders/v2/#orders-create-request-body): item_total, purchase_units, and more.

```javascript
    createOrder: function (data, actions) {
      return fetch("/api/orders", {
        method: "post",
      })
        .then((response) => response.json())
        .then((order) => order.id);
    },
```

### Approve Order
The `onApprove()` allows doing something with the order details after the order has been created. You can learn more about the [onApprove()](https://developer.paypal.com/sdk/js/reference/#link-onapprove) callback in your docs.

```javascript
onApprove: function (data, actions) {
      return fetch(`/api/orders/${data.orderID}/capture`, {
        method: "post",
      })
        .then((response) => response.json())
        .then((orderData) => {
          // Successful capture! For dev/demo purposes:
          console.log(
            "Capture result",
            orderData,
            JSON.stringify(orderData, null, 2)
          );
          const transaction = orderData.purchase_units[0].payments.captures[0];
          alert(`Transaction ${transaction.status}: ${transaction.id}

            See console for all available details
          `);
        });
    },
```

## Main.go
This Go server is using [Gorilla/Mux](https://github.com/gorilla/mux) to create the request router.

### Create Order
The `createOrder()` is called on the `createOrder()`callback from the PayPal button to create an order. In this function, you can create the request of your order with the following properties from the [V2 orders create call](https://developer.paypal.com/api/orders/v2/#orders-create-request-body): item_total, purchase_units, and more.

Endpoint: `/api/orders`

```go
func createOrder(w http.ResponseWriter, r *http.Request) {
	accessToken := generateAccessToken()
	body := []byte(`{
		"intent":"CAPTURE",
		"purchase_units":[
		   {
			  "amount":{
				 "currency_code":"USD",
				 "value":"100.00"
			  }
		   }
		]
	 }`)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	req, err := http.NewRequest("POST", base+"/v2/checkout/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Println("request failed with status:", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}
```

### Capture Payment
The `capturePayment()` function is called on the `onApprove()` callback from the PayPal button once an order has been completed. In this function, you can use the value returned to be stored in a database or to do something else with it.

Endpoint: `/api/orders/:orderID/capture`

```javascript
orderID := mux.Vars(r)["orderID"]
	accessToken := generateAccessToken()
	body := []byte(`{}`)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	req, err := http.NewRequest("POST", base+"/v2/checkout/orders/"+orderID+"/capture", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Println("request failed with status:", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		return
	}

	// copy response from external to frontend
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}
```

### Run this App

You will need to create a `.env` file with the following environment variables:

```shell
CLIENT_ID=
APP_SECRET=
```

Complete the steps in [Get started](https://developer.paypal.com/api/rest/) to get the following sandbox account information from the Developer Dashboard:
- Sandbox client ID and the secret of [a REST app](https://www.paypal.com/signin?returnUri=https%3A%2F%2Fdeveloper.paypal.com%2Fdeveloper%2Fapplications&_ga=1.252581760.841672670.1664266268).
- Access token to use the PayPal REST API server.

![paypal developer credentials](env.png)

Now, run the following command in your terminal to start the server:

`go run main.go`

and open the `index.html` in your browswer.

## PayPal Developer Community
The PayPal Developer community helps you build your career while improving your products and the developer experience. You’ll be able to contribute code and documentation, meet new people and learn from the open-source community.
 
* Website: [developer.paypal.com](https://developer.paypal.com)
* Twitter: [@paypaldev](https://twitter.com/paypaldev)
* GitHub:  [@paypal](https://github.com/paypal)

