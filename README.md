# API Server

This is the api server for the bip software.

``` 
   go run ./internal/seeders/start-seeders/main.go
```

### Prerequisites 

- Git  
- App 
- cmd/integration_task


#### Git

In order to ensure the smooth functioning of the backend, it is essential to have the Git Server running on a publicly accessible IP address. To facilitate this, we require the deployment of the BIP GIT repository, which can be found at https://github.com/bip-so/git.

To configure the environment for the backend, please include the following variables in your .env file:

GIT_HOSTS:

This variable represents the hostname or IP address of the Git Server.
Example: `GIT_HOSTS=<your_git_server_hostname>`
GIT_SECRET_KEY:

This variable should be set to the same value as the SECRET_KEY used by the Git Server.
Example:`GIT_SECRET_KEY=<your_git_server_secret_key>`
By providing the appropriate values for GIT_HOSTS and GIT_SECRET_KEY in the .env file, you will ensure the backend can establish a connection with the Git Server and perform necessary operations seamlessly.



#### App

FRONTEND_HOST:
This environment variable represents the endpoint of your frontend code repository.
Example: FRONTEND_HOST=https://github.com/bip-so/app
By including the FRONTEND_HOST endpoint in your .env file, your backend services will be able to establish a connection with the frontend application and enable seamless data transfer and interaction.


#### Integration Task Server

You can find the code required for integration server in the `./cmd/integration_task` folder. This code contains the necessary logic to handle events from Discord and Slack effectively. Please make sure to review and deploy the Lambda functions correctly to ensure proper functionality.


### Services 

We require following services to get the app started. 

* Algolia Search: Visit Algolia's website (https://www.algolia.com/) and navigate to their documentation section for more details.
* Discord app: To integrate Discord into your application, you can utilize the Discord API. The official Discord Developer Portal (https://discord.com/developers/docs/intro) offers comprehensive documentation and guides for setting up a Discord app, creating a bot, and interacting with the API.
* Get Stream: Get Stream provides detailed installation instructions for various platforms and frameworks. You can visit their official website (https://getstream.io/) and explore the documentation section to get started. They offer SDKs and code examples for easy integration.
* Kafka: Kafka has comprehensive documentation available on the Apache Kafka website (https://kafka.apache.org/). It provides step-by-step guides for installing Kafka on different operating systems, along with detailed explanations of its architecture and usage.
* Redis: Redis provides installation instructions tailored to various operating systems on their official website (https://redis.io/). You can refer to their documentation to set up Redis based on your specific requirements.
* S3 Bucket: To use S3, you need to sign up for an AWS account and follow their documentation (https://aws.amazon.com/s3/) for creating an S3 bucket and configuring access.
* Cloudfront URL: You can refer to the AWS CloudFront documentation (https://aws.amazon.com/cloudfront/) for detailed instructions on setting up CloudFront and obtaining a CloudFront URL for your content.
* SES (Amazon Simple Email Service): You can find installation and setup instructions in the Amazon SES documentation (https://aws.amazon.com/ses/).
* Supabase: Installation: Supabase offers a comprehensive documentation portal (https://supabase.io/docs/) that guides you through the installation and usage of their services. You can refer to their documentation to get started.
* Slack app: Installation: To create and integrate a Slack app, you can refer to the Slack API documentation (https://api.slack.com/). It provides step-by-step guides for creating a Slack app, configuring permissions, and interacting with the Slack API.


We have added a sample `sample.env.txt` file 

Finally you can run the code with

`go run main.go`