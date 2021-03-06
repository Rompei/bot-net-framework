package main
import (
	"github.com/ChimeraCoder/anaconda"
	"runtime"
	"github.com/jinzhu/gorm"
)

// TalkController is controller of talk.
type TalkController struct {
	Talk Talk
	Seq int
	db *gorm.DB
}

// NewTalkController is constructor of TalkController.
func NewTalkController(talk Talk, db *gorm.DB)*TalkController{
	return &TalkController{
		Talk:talk,
		Seq:0,
		db:db,
	}
}

// PostOne posts one tweet and inclement sequence.
func (tc *TalkController) PostOne() (Tweets, error){



	numTweet := len(tc.Talk.Tweets[tc.Seq])
	resultCh := make(chan Tweet, numTweet)
	errCh := make(chan error, runtime.NumCPU())

	for _,v  := range tc.Talk.Tweets[tc.Seq]{
		go postTweet(v, resultCh, errCh)
	}

	count := 0
	var err error
	L1:
	for{
		select {
		case result := <-resultCh:
			tc.db.Save(result)
			count++
			if count == numTweet {
				break L1
			}
		case err = <-errCh:
			break L1
		default:
		}
	}

	// Increment sequence
	tc.Seq++
	return tc.Talk.Tweets[tc.Seq-1], err
}

func postTweet(tweet Tweet, resultCh chan Tweet, errCh chan error) {

	// This function is for posting tweet asynchronously

	api := anaconda.NewTwitterApi(tweet.Bot.AccessToken, tweet.Bot.AccessTokenSecret)

	result, err := api.PostTweet(tweet.Text, nil)
	if err != nil {
		errCh <- err
		return
	}
	tweet.TweetIdStr = result.IdStr
	tweet.Bot = Bot{}
	resultCh <- tweet
}