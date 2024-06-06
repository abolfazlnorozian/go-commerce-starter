package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendSms(verifyCode, phoneNumber string) {
	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := fmt.Sprintf(`
  {
	"pattern_code": "wshbtvj694w8uni",
	"originator": "+983000505",
	"recipient": "%s",
	"values": {
	  "verification-code": "%s"
	}
  }
`, phoneNumber, verifyCode)
	// Convert payload to io.Reader
	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

func BackPaymentSendSms(verifyCode, phoneNumber string) {
	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := fmt.Sprintf(`
  {
	"pattern_code": "05rz9m57ae49j6u",
	"originator": "+983000505",
	"recipient": "%s",
	"values": {
	  "verification-code": "%s"
	}
  }
`, phoneNumber, verifyCode)
	// Convert payload to io.Reader
	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

// func SendSmstoAdmin(verifyCode, phoneNumber string) {
// 	url := "http://rest.ippanel.com/v1/messages/patterns/send"
// 	method := "POST"

// 	payload := fmt.Sprintf(`
//   {
// 	"pattern_code": "cqxonayv27ogoth",
// 	"originator": "+983000505",
// 	"recipient": "%s",
// 	"values": {
// 	  "verification-code": "%s"
// 	}
//   }
// `, phoneNumber, verifyCode)
// 	// Convert payload to io.Reader
// 	payloadReader := strings.NewReader(payload)

// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, payloadReader)

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
// 	req.Header.Add("Content-Type", "application/json")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(body))

// }

func SendSmstoCustomerForReadyProduct(verifyCode, phoneNumber string) {
	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := fmt.Sprintf(`
  {
	"pattern_code": "c19mua1xv9o9q3n",
	"originator": "+983000505",
	"recipient": "%s",
	"values": {
	  "verification-code": "%s"
	}
  }
`, phoneNumber, verifyCode)
	// Convert payload to io.Reader
	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
func SendSmstoCustomerForSendProduct(verifyCode, postalTrackingCode, phoneNumber string) {
	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := map[string]interface{}{
		"pattern_code": "kwp8i54jjcdug03",
		"originator":   "+983000505",
		"recipient":    phoneNumber,
		"values": map[string]string{
			"verification-code":   verifyCode,
			"postalTracking-code": postalTrackingCode,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return
	}
	payloadReader := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK HTTP status: %s\n", res.Status)
		fmt.Printf("Response body: %s\n", string(body))
		return
	}

	fmt.Println("Response:", string(body))
}

// func SendSmstoCustomerForSendProduct(verifyCode, postalTrackingCode, phoneNumber string) {
// 	url := "http://rest.ippanel.com/v1/messages/patterns/send"
// 	method := "POST"

// 	payload := fmt.Sprintf(`
//   {
// 	"pattern_code": "kwp8i54jjcdug03",
// 	"originator": "+983000505",
// 	"recipient": "%s",
// 	"values": {
// 	  "verification-code": "%s"
// 	}
//   }
// `, phoneNumber, postalTrackingCode, verifyCode)
// 	// Convert payload to io.Reader
// 	payloadReader := strings.NewReader(payload)

// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, payloadReader)

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
// 	req.Header.Add("Content-Type", "application/json")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(body))

// }

func SendSmstoAdmin(verifyCode string, phoneNumbers []string) {
	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	for _, phoneNumber := range phoneNumbers {
		payload := fmt.Sprintf(`
        {
            "pattern_code": "cqxonayv27ogoth",
            "originator": "+983000505",
            "recipient": "%s",
            "values": {
                "verification-code": "%s"
            }
        }
        `, phoneNumber, verifyCode)

		// Convert payload to io.Reader
		payloadReader := strings.NewReader(payload)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payloadReader)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
	}
}
