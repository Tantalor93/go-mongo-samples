package domain

type Work struct {
	Author   string `bson:"author"`
	Country  string `bson:"country"`
	Language string `bson:"language"`
	Link     string `bson:"link"`
	Pages    int    `bson:"pages"`
	Title    string `bson:"title"`
	Year     int    `bson:"year"`
}
