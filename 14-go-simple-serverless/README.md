## Some Important AWS-CLI commands to run

You need to configure your AWS CLI. As a developer you would have access to 
different profiles (development and production) and all of those profiles can 
stay in your `./aws` config and credentials files as different user profiles.

### Approach 1
```bash
aws iam create-role -role-name lambda-ex --asume-role-policy-document '{"Version": "2012-10-17", "Statement": [{ "Effect": "Allow", "Principal": {"service": "lambda.amazonaws.com"}, "Action": "sts"AssumeRole"}]}'
```

## Approach 2

This needs to followed by adding the `trust-policy.json` file.
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```
Creating roles based on `trust-policy.json` as a source of truth instead of 
writing down the entire policy by yourself.
```bash
aws iam create-role --role-name lambda-ex --asume-role-policy-document \
file://trust-policy.json
```

### Continued from Both Approaches

Attaching a role policy to an existing role
```bash
aws iam attach-role-policy --role-name lambda-ex --policy-arn \
arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

For deploying to Lambda
```bash
aws lambda create-function --function-name go-lambda-1 --zip-file \
file://function.zip --handler main --runtime go1.x \
arn:aws> --role arn:aws:iam:<account-id>:role/lambda-ex
```

### Invoking the Function
```bash
aws lambda invoke --function-name go-lambda --clli-binary-format \
raw-in-base64-out --payload '{"What is your name?": "Jim", "How old are you?": 33}' output.txt
```
