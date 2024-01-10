module github.com/cfc-servers/cfc_suggestions

go 1.14

require (
	github.com/aws/aws-lambda-go v1.44.0
	github.com/aws/aws-sdk-go v1.36.25
	github.com/go-resty/resty/v2 v2.4.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/guregu/dynamo v1.10.2
	github.com/plally/goslash v0.1.2-0.20210501142734-b97cdda6fdfb
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.2
)

replace github.com/bwmarrin/discordgo => github.com/plally/discordgo v0.23.3-0.20210504211656-c9fe4fa407ce
