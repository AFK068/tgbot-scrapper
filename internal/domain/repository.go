package domain

type ChatLinkRepository interface {
	// Registers a chat in the repository.
	RegisterChat(chatID int64) error

	// Deletes a chat from the repository.
	DeleteChat(chatID int64) error

	// Save a link in the repository.
	SaveLink(chatID int64, link *Link) error

	// Delete a link from the repository.
	DeleteLink(chatID int64, link *Link) error

	// Get a list of links from the repository.
	GetListLinks(chatID int64) ([]*Link, error)

	// Check if the chat exists in the repository.
	CheckUserExistence(chatID int64) bool

	// Get all links from the repository.
	GetAllLinks() []*Link
}
