package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func SendSmstoAdmin(verifyCode string, phoneNumbers []string) error {
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
			return err
		}
		req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(string(body))
	}
	return nil

}
func SendSmstoAllUsers(message string, phoneNumbers []string) error {
	url := "https://api2.ippanel.com/api/v1/sms/send/webservice/single"
	method := "POST"
	sender := "+98Clubs"
	//10004223
	apiKey := "wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ="

	var sentNumbers []string
	var failedNumbers []string

	// Escape newlines in the message
	message = strings.ReplaceAll(message, "\n", "\\n")

	for _, phoneNumber := range phoneNumbers {
		payload := fmt.Sprintf(`
        {
            "recipient": ["%s"],
            "sender": "%s",
            "message": "%s"
        }
        `, phoneNumber, sender, message)

		// Convert payload to io.Reader
		payloadReader := strings.NewReader(payload)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payloadReader)
		if err != nil {
			fmt.Println("Error creating request:", err)
			failedNumbers = append(failedNumbers, phoneNumber)
			continue
		}
		req.Header.Add("apikey", apiKey)
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			failedNumbers = append(failedNumbers, phoneNumber)
			continue
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			failedNumbers = append(failedNumbers, phoneNumber)
			continue
		}
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error unmarshalling response:", err)
			failedNumbers = append(failedNumbers, phoneNumber)
			continue
		}

		status := response["status"].(string)
		if status == "Forbidden" || status == "Unprocessable Entity" {
			fmt.Printf("Failed to send SMS to %s: %s\n", phoneNumber, string(body))
			failedNumbers = append(failedNumbers, phoneNumber)
		} else {
			fmt.Printf("Successfully sent SMS to %s: %s\n", phoneNumber, string(body))
			sentNumbers = append(sentNumbers, phoneNumber)
		}
	}
	log.Printf("Successfully sent SMS to: %v", sentNumbers)
	log.Printf("Failed to send SMS to: %v", failedNumbers)

	return nil
}
