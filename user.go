package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
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

type Subscription struct {
	ID          string    `json:"_id"`
	SubPlan     string    `json:"sub_plan"`
	SubPlanName string    `json:"sub_plan_name"`
	Channel     Channel   `json:"channel"`
	CreatedAt   time.Time `json:"created_at"`
}

type Follow struct {
	Notifications bool      `json:"notifications"`
	Channel       Channel   `json:"channel"`
	CreatedAt     time.Time `json:"created_at"`
}

type Follows struct {
	Total   uint     `json:"_total"`
	Follows []Follow `json:"follows"`
}

type Block struct {
	ID        string    `json:"_id"`
	User      User      `json:"user"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUser retrieves the user based on the access token attached to the client.
func (ac AccessClient) GetUser() (User, error) {
	var user User

	if err := ac.validateScope("user_read"); err != nil {
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
func (c Client) GetUserByID(id string) (User, error) {
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

// GetUsersByName accepts a list of usernames to fetch from Twitch. You can include
// up to 100 names and get a slice of basic user information back, including the ID.
func (c Client) GetUsersByName(names ...string) ([]User, error) {
	var userRes struct {
		Total uint `json:"_total"`
		Users []User
	}

	uri := fmt.Sprintf("%s/users?login=%s", baseURI, strings.Join(names, ","))

	res, err := c.makeGetRequest(uri)
	if err != nil {
		return []User{}, err
	}

	if err := json.NewDecoder(res.Body).Decode(&userRes); err != nil {
		return []User{}, err
	}

	return userRes.Users, nil
}

// GetUserSubscription returns whether or not a given user ID is subscribed to the given channel ID.
func (ac AccessClient) GetUserSubscription(userID, channelID string) (Subscription, error) {
	var subscription Subscription
	uri := fmt.Sprintf("%s/users/%s/subscriptions/%s", baseURI, userID, channelID)

	if err := ac.validateScope("user_subscriptions"); err != nil {
		return subscription, err
	}

	res, err := ac.makeGetRequest(uri)
	if err != nil {
		return subscription, err
	}

	if res.StatusCode == 404 {
		return subscription, fmt.Errorf("user is not subscribed")
	}

	if res.StatusCode == 422 {
		return subscription, fmt.Errorf("channel does not have a subscription program")
	}

	if err := json.NewDecoder(res.Body).Decode(&subscription); err != nil {
		return subscription, err
	}

	log.Printf("StatusCode: %d", res.StatusCode)
	log.Println("MADE IT")

	return subscription, nil
}

// GetUserFollows retrieves the list of channels that a given user follows.
// TODO: Add in support for direction/sortby.
func (c Client) GetUserFollows(userID string, limit, offset int) (Follows, error) {
	var follows Follows
	query := url.Values{}
	query.Add("limit", strconv.Itoa(limit))
	query.Add("offset", strconv.Itoa(offset))

	uri := fmt.Sprintf("%s/users/%s/follows/channels?%s", baseURI, userID, query.Encode())

	res, err := c.makeGetRequest(uri)
	if err != nil {
		return follows, err
	}

	if err := json.NewDecoder(res.Body).Decode(&follows); err != nil {
		return follows, err
	}

	return follows, nil
}

// CheckUserFollowsChannel returns a Follow payload if the user follows the given channel.
// If the user doesn't follow the channel, an error is returned.
func (c Client) CheckUserFollowsChannel(userID, channelID string) (Follow, error) {
	var follow Follow

	uri := fmt.Sprintf("%s/users/%s/follows/channels/%s", baseURI, userID, channelID)

	res, err := c.makeGetRequest(uri)
	if err != nil {
		return follow, err
	}

	if res.StatusCode == 404 {
		return follow, fmt.Errorf("User %s does not follow channel %s", userID, channelID)
	}

	if err := json.NewDecoder(res.Body).Decode(&follow); err != nil {
		return follow, err
	}

	return follow, nil
}

func (ac AccessClient) FollowChannel(userID, channelID string, notify bool) (Follow, error) {
	var follow Follow
	if err := ac.validateScope("user_follows_edit"); err != nil {
		return follow, err
	}

	uri := fmt.Sprintf("%s/users/%s/follows/channels/%s", baseURI, userID, channelID)

	var notifications = struct {
		Notifications bool `json:"notifications"`
	}{notify}

	payload, err := json.Marshal(&notifications)
	if err != nil {
		return follow, err
	}

	res, err := ac.makePutRequest(uri, payload)
	if err != nil {
		return follow, err
	}

	if res.StatusCode == 422 {
		return follow, fmt.Errorf("User %s could not follow Channel %s", userID, channelID)
	}

	if err := json.NewDecoder(res.Body).Decode(&follow); err != nil {
		return follow, err
	}

	return follow, err
}

func (ac AccessClient) UnfollowChannel(userID, channelID string) error {
	if err := ac.validateScope("user_follows_edit"); err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/users/%s/follows/channels/%s", baseURI, userID, channelID)

	res, err := ac.makeDeleteRequest(uri)
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("Failed to unfollow User %s from Channel %s", userID, channelID)
	}

	return nil
}

func (ac AccessClient) BlockUser(userID, blockedID string) (Block, error) {
	var block Block
	if err := ac.validateScope("user_blocks_edit"); err != nil {
		return block, err
	}

	uri := fmt.Sprintf("%s/users/%s/blocks/%s", baseURI, userID, blockedID)

	res, err := ac.makePutRequest(uri, nil)
	if err != nil {
		return block, err
	}

	if err := json.NewDecoder(res.Body).Decode(&block); err != nil {
		return block, err
	}

	return block, nil
}
