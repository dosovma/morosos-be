build:
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go

zip:
	zip myFunction.zip bootstrap

make auth:
	aws sso login

create:
	aws lambda create-function \
    --function-name testGolambda_v1 \
    --zip-file fileb://myFunction.zip \
    --handler bootstrap \
    --runtime provided.al2023 \
    --architecture arm64 \
    --role arn:aws:iam::539247474956:role/lambda-apigatewayEUS-policy

update:
	aws lambda update-function-code \
	--function-name testGolambda_v1 \
	--zip-file fileb://myFunction.zip

invoke:
	aws lambda invoke \
	--function-name testGolambda_v1 \
	--payload file://test/input.json ./test/outputfile.txt \
	--cli-binary-format raw-in-base64-out

invoke2:
	aws lambda invoke \
	--function-name testGolambda_v1 \
	--payload '{"PK":{"S":"form1"}, "SK":{"S":"01236"}}' ./test/outputfile.txt \
	--cli-binary-format raw-in-base64-out


put-item:
	aws dynamodb put-item \
	--table-name testTable \
  	--item \
	'{"PK":{"S":"form1"}, "SK":{"S":"01236"}}'

curl:
	curl https://dmq5zvn1g5.execute-api.eu-south-2.amazonaws.com/test/DynamoDBManager \
    -d '{"PK":"form123", "SK":"success"}'