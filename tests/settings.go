package tests

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = ``

// type signInResponse struct {
// 	Token string `json:"token"`
// }

// func getToken() string {
// 	password := os.Getenv("TODO_PASSWORD")
// 	port := os.Getenv("TODO_PORT")
// 	if port == "" {
// 		port = "7540"
// 	}
// 	if len(password) == 0 {
// 		return ""
// 	}

// 	payload, err := json.Marshal(map[string]string{"password": password})
// 	if err != nil {
// 		return ""
// 	}

// 	client := &http.Client{Timeout: 15 * time.Second}
// 	resp, err := client.Post("http://localhost:"+port+"/api/signin", "application/json", bytes.NewBuffer(payload))
// 	if err != nil {
// 		return ""
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return ""
// 	}

// 	var body signInResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
// 		return ""
// 	}

// 	return body.Token
// }
