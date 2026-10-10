package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ippanelAPIKey() string {
	return os.Getenv("IPPANEL_API_KEY")
}

func ippanelOriginator() string {
	return envOr("IPPANEL_ORIGINATOR", "+983000505")
}

func ippanelPatternVerify() string {
	return os.Getenv("IPPANEL_PATTERN_VERIFY")
}

func ippanelPatternPayment() string {
	return os.Getenv("IPPANEL_PATTERN_PAYMENT")
}

func ippanelPatternReadyProduct() string {
	return os.Getenv("IPPANEL_PATTERN_READY_PRODUCT")
}

func ippanelPatternSendProduct() string {
	return os.Getenv("IPPANEL_PATTERN_SEND_PRODUCT")
}

func ippanelPatternAdmin() string {
	return os.Getenv("IPPANEL_PATTERN_ADMIN")
}

func ippanelBroadcastSender() string {
	return envOr("IPPANEL_SMS_SENDER", "+98Clubs")
}

// smsPatternConfigured returns false when API key or pattern code is missing,
// so callers skip the HTTP request instead of sending an invalid payload.
func smsPatternConfigured(patternCode, context string) bool {
	if ippanelAPIKey() == "" {
		log.Printf("sms: skipping %s — IPPANEL_API_KEY is not set", context)
		return false
	}
	if patternCode == "" {
		log.Printf("sms: skipping %s — pattern code env var is not set", context)
		return false
	}
	return true
}

func SendSms(verifyCode, phoneNumber string) {
	pattern := ippanelPatternVerify()
	if !smsPatternConfigured(pattern, "SendSms") {
		return
	}

	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := fmt.Sprintf(`
  {
	"pattern_code": "%s",
	"originator": "%s",
	"recipient": "%s",
	"values": {
	  "verification-code": "%s"
	}
  }
`, pattern, ippanelOriginator(), phoneNumber, verifyCode)
	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "AccessKey "+ippanelAPIKey())
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
	pattern := ippanelPatternPayment()
	if !smsPatternConfigured(pattern, "BackPaymentSendSms") {
		return
	}

	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := fmt.Sprintf(`
  {
	"pattern_code": "%s",
	"originator": "%s",
	"recipient": "%s",
	"values": {
	  "verification-code": "%s"
	}
  }
`, pattern, ippanelOriginator(), phoneNumber, verifyCode)
	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "AccessKey "+ippanelAPIKey())
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
	pattern := ippanelPatternReadyProduct()
	if !smsPatternConfigured(pattern, "SendSmstoCustomerForReadyProduct") {
		return
	}

	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := fmt.Sprintf(`
  {
	"pattern_code": "%s",
	"originator": "%s",
	"recipient": "%s",
	"values": {
	  "verification-code": "%s"
	}
  }
`, pattern, ippanelOriginator(), phoneNumber, verifyCode)
	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "AccessKey "+ippanelAPIKey())
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
	pattern := ippanelPatternSendProduct()
	if !smsPatternConfigured(pattern, "SendSmstoCustomerForSendProduct") {
		return
	}

	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := map[string]interface{}{
		"pattern_code": pattern,
		"originator":   ippanelOriginator(),
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
	req.Header.Add("Authorization", "AccessKey "+ippanelAPIKey())
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
	pattern := ippanelPatternAdmin()
	if !smsPatternConfigured(pattern, "SendSmstoAdmin") {
		return nil
	}

	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	for _, phoneNumber := range phoneNumbers {
		payload := fmt.Sprintf(`
        {
            "pattern_code": "%s",
            "originator": "%s",
            "recipient": "%s",
            "values": {
                "verification-code": "%s"
            }
        }
        `, pattern, ippanelOriginator(), phoneNumber, verifyCode)

		payloadReader := strings.NewReader(payload)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payloadReader)
		if err != nil {
			fmt.Println(err)
			return err
		}
		req.Header.Add("Authorization", "AccessKey "+ippanelAPIKey())
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
	apiKey := ippanelAPIKey()
	if apiKey == "" {
		log.Printf("sms: skipping SendSmstoAllUsers — IPPANEL_API_KEY is not set")
		return nil
	}

	url := "https://api2.ippanel.com/api/v1/sms/send/webservice/single"
	method := "POST"
	sender := ippanelBroadcastSender()

	var sentNumbers []string
	var failedNumbers []string

	message = strings.ReplaceAll(message, "\n", "\\n")

	for _, phoneNumber := range phoneNumbers {
		payload := fmt.Sprintf(`
        {
            "recipient": ["%s"],
            "sender": "%s",
            "message": "%s"
        }
        `, phoneNumber, sender, message)

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

		status, _ := response["status"].(string)
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
