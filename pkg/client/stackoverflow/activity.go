package stackoverflow

type ActivityType string

const (
	ActivityTypeComment  ActivityType = "comment"
	ActivityTypeAnswer   ActivityType = "answer"
	ActivityTypeQuestion ActivityType = "question"
)

type Activity struct {
	Type      ActivityType
	CreatedAt int64
	Body      string
	Tags      []string
	UserName  string
}

func NewActivity(activityType ActivityType, createdAt int64, body string, tags []string, userName string) *Activity {
	return &Activity{
		Type:      activityType,
		CreatedAt: createdAt,
		Body:      body,
		Tags:      tags,
		UserName:  userName,
	}
}
