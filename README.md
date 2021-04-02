# Sigma

Sigma is a tool that generates a serverless REST or GraphQL api in Javascript.

## Usage

Sigma generates an api from an api specification file. You can define an api spec as JSON or as a GraphQL schema.

1. Create a file named `api_spec.gql` or `api_spec.json`
2. Install the api compiler
3. In the terminal, run `sigma compile -o my_api`
4. Change into the deployment directory `cd my_api/.deploy`
5. Deploy the api by executing `./linux_deploy` for Linux, `./macos_deploy` for macOS, and `Start-Process -FilePath WinDeploy.exe` for Windows.

## Adding Custom Endpoints

1. Create a file named `api_ext.json`
2. Create the root object
```json
{
  "endpoints": [],
  "variables": []
}
```
3. Add an endpoint with the following format for REST apis
```json
{
  "route": "/my/custom/route",
  "methodFunctions": {
    "get": {
      "cloudProvider": "aws",
      "credentials": {
        "accessId": "${{ variables.AWS_ACCESS_ID }}",
        "secretKey": "${{ variables.AWS_SECRET_KEY }}",
        "roleArn": "${{ variables.AWS_ROLE_ARN }}"
      },
      "functionArn": "arn:aws:lambda:us-west-2:123456789012:function:my-endpoint-handler"
    }
  }
}
```