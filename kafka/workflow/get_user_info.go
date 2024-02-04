package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

func GetUserInfo(name string) (User, error) {

	var user User

	req, err := http.NewRequest("GET", "https://api.github.com/users/"+name, nil)
	if err != nil {
		return user, err
	}

	req.Header.Set("Authorization", "token "+"")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("Error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return user, err
	}

	return user, nil
}
