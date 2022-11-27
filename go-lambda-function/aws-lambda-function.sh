aws iam create-role --role-name lambda-ex --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}'
aws iam create-role --role-name lambda-ex --assume-role-policy-document file://trust-policy.json
aws iam attach-role-policy --role-name lambda-ex --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
aws lambda create-function --function-name go-lambda-function --zip-file fileb://function.zip --handler main --runtime go1.x --role arn:aws:iam::390229387745:role/lambda-ex
aws lambda invoke --function-name go-lambda-function --cli-binary-format raw-in-base64-out --payload '{"what is your name?": "Dat", "how old are you?": 21}' output.txt