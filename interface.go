package slack

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

type EpochTime int64

// DefaultSlackAPIEndpoint contains the prefix used for Slack REST API
const (
	DefaultAPIEndpoint         = "https://slack.com/api/"
	DefaultOAuth2AuthEndpoint  = "https://slack.com/oauth/authorize"
	DefaultOAuth2TokenEndpoint = "https://slack.com/api/oauth.access"
)

// Oauth2Endpoint contains the Slack OAuth2 endpoint configuration
var OAuth2Endpoint = oauth2.Endpoint{
	AuthURL:  DefaultOAuth2AuthEndpoint,
	TokenURL: DefaultOAuth2TokenEndpoint,
}

type Client struct {
	auth     *AuthService
	chat     *ChatService
	users    *UsersService
	slackURL string
	token    string
}

type SlackResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type AuthService struct {
	client *httpClient
	token  string
}

type AuthTestResponse struct {
	URL    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type Attachment interface{} // TODO
type ChatService struct {
	client *httpClient
	token  string
}

type ChatResponse struct {
	Channel   string      `json:"channel"`
	Timestamp string      `json:"ts"`
	Message   interface{} `json:"message"` // TODO
}

type Edited struct {
	Timestamp string `json:"ts"`
	User      string `json:"user"`
}

type Comment struct {
	ID        string    `json:"id,omitempty"`
	Created   EpochTime `json:"created,omitempty"`
	Timestamp EpochTime `json:"timestamp,omitempty"`
	User      string    `json:"user,omitempty"`
	Comment   string    `json:"comment,omitempty"`
}

// Message is a representation of a message, as obtained
// by the RTM or Events API. This is NOT what you use when
// you are posting a message. See ChatService#PostMessage
// and MessageParams for that.
type Message struct {
	Attachments []Attachment `json:"attachments"`
	Channel     string       `json:"channel"`
	Edited      *Edited      `json:"edited"`
	IsStarred   bool         `json:"is_starred"`
	PinnedTo    []string     `json:"pinned_to"`
	Text        string       `json:"text"`
	Timestamp   string       `json:"timestamp"`
	Type        string       `json:"type"`
	User        string       `json:"user"`

	// Message Subtypes
	Subtype string `json:"subtype"`

	// Hidden Subtypes
	Hidden           bool   `json:"hidden,omitempty"`     // message_changed, message_deleted, unpinned_item
	DeletedTimestamp string `json:"deleted_ts,omitempty"` // message_deleted
	EventTimestamp   string `json:"event_ts,omitempty"`

	// bot_message (https://api.slack.com/events/message/bot_message)
	BotID    string `json:"bot_id,omitempty"`
	Username string `json:"username,omitempty"`
	Icons    *Icon  `json:"icons,omitempty"`

	// channel_join, group_join
	Inviter string `json:"inviter,omitempty"`

	// channel_topic, group_topic
	Topic string `json:"topic,omitempty"`

	// channel_purpose, group_purpose
	Purpose string `json:"purpose,omitempty"`

	// channel_name, group_name
	Name    string `json:"name,omitempty"`
	OldName string `json:"old_name,omitempty"`

	// channel_archive, group_archive
	Members []string `json:"members,omitempty"`

	// file_share, file_comment, file_mention
	//	File *File `json:"file,omitempty"`

	// file_share
	Upload bool `json:"upload,omitempty"`

	// file_comment
	Comment *Comment `json:"comment,omitempty"`

	// pinned_item
	ItemType string `json:"item_type,omitempty"`

	// https://api.slack.com/rtm
	ReplyTo int    `json:"reply_to,omitempty"`
	Team    string `json:"team,omitempty"`

	// reactions
	Reactions []ItemReaction `json:"reactions,omitempty"`
}

type Icon struct {
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

type ItemReaction struct {
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Users []string `json:"users"`
}

// MessageParans is used when posting a new message
type MessageParams struct {
	AsUser      bool
	Attachments []Attachment
	EscapeText  bool
	IconEmoji   string
	IconURL     string
	LinkNames   bool
	Markdown    bool `json:"mrkdwn,omitempty"`
	Parse       string
	UnfurlLinks bool
	UnfurlMedia bool
	Username    string
}

type UsersService struct {
	client *httpClient
	token  string
}

type UserProfile struct {
	AlwaysActive       bool   `json:"always_active"`
	AvatarHash         string `json:"avatar_hash"`
	FirstName          string `json:"first_name"`
	Image24            string `json:"image_24"`
	Image32            string `json:"image_32"`
	Image48            string `json:"image_48"`
	Image72            string `json:"image_72"`
	Image192           string `json:"image_192"`
	Image512           string `json:"image_512"`
	LastName           string `json:"last_name"`
	RealName           string `json:"real_name"`
	RealNameNormalized string `json:"real_name_normalized"`
}

type User struct {
	Color             string      `json:"color"`
	Deleted           bool        `json:"deleted"`
	ID                string      `json:"id"`
	IsAdmin           bool        `json:"is_admin"`
	IsBot             bool        `json:"is_bot"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	Name              string      `json:"name"`
	Profile           UserProfile `json:"profile"`
	RealName          string      `json:"real_name"`
	Status            string      `json:"status,omitempty"`
	TeamID            string      `json:"team_id"`
	TZ                string      `json:"tz,omitempty"`
	TZLabel           string      `json:"tz_label"`
	TZOffset          int         `json:"tz_offset"`
	Update            int         `json:"updated"`
}

type UserList []*User

type Presence string

const (
	Presencective Presence = "away"
	PresenceAway  Presence = "away"
)

type UserPresence struct {
	AutoAway        bool     `json:"auto_away,omitempty"`
	ConnectionCount int      `json:"connection_count,omitempty"`
	LastActivity    int      `json:"last_activity,omitempty"`
	ManualAway      bool     `json:"manual_away,omitempty"`
	Online          bool     `json:"online"`
	Presence        Presence `json:"presence"`
}

type EventType int

const (
	RTMConnectingEvent EventType = iota
	MaxEvent
)

type RTM struct {
	outch chan Event
}

type RTMEvent struct {
	typ  EventType
	data interface{}
}

type eventUnmarshalProxy struct {
	EventTimestamp string
	Item           json.RawMessage
	Timestamp      string
	Type           string
	User           string
}

type Event struct {
	EventTimestamp string
	Item           interface{}
	Timestamp      string
	Type           string
	User           string
}
