package twitch

import (
	"encoding/json"
	"fmt"
	"time"
)

// A User contains all of the data Twitch returns about a given user. Based on the method
// of retrieval, some of these fields may be omitted.
type User struct {
	ID               string        `json:"_id"`
	Bio              string        `json:"bio"`
	DisplayName      string        `json:"display_name"`
	Email            string        `json:"email,omitempty"`
	EmailVerified    bool          `json:"email_verified,omitempty"`
	Logo             string        `json:"logo"`
	Name             string        `json:"name"`
	Notifications    Notifications `json:"notifications,omitempty"`
	Partnered        bool          `json:"partnered,omitempty"`
	TwitterConnected bool          `json:"twitter_connected,omitempty"`
	Type             string        `json:"type"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type Notifications struct {
	Email bool `json:"email"`
	Push  bool `json:"push"`
}

// GetUser retrieves the user based on the access token attached to the client.
func (ac AccessClient) GetUser() (User, error) {
	var user User

	if err := ac.access.validateScope("user_read"); err != nil {
		return user, err
	}

	uri := fmt.Sprintf("%s/user", baseURI)

	res, err := ac.makeGetRequest(uri)
	if err != nil {
		return user, err
	}

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByID retrieves a user based on the given user ID.
func (c *Client) GetUserByID(id uint) (User, error) {
	var user User
	uri := fmt.Sprintf("%s/user", baseURI)

	res, err := c.makeGetRequest(uri)
	if err != nil {
		return user, err
	}

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}
