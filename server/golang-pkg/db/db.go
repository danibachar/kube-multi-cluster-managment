import {

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
}

func NewDBRef() Firestore {
	ctx := context.Background()
	opt := option.WithCredentialsFile("../../kmcm-2f5fc-firebase-adminsdk-3p6in-10b5fde3c9.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
	  return nil, fmt.Errorf("error initializing app: %v", err)
	}
	
	client, err := app.Firestore(ctx)
	if err != nil {
	  log.Fatalln(err)
	}
	return client
}

