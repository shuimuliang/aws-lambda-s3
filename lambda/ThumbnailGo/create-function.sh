aws lambda --profile test create-function --function-name ThumbnailGo --runtime go1.x \
  --zip-file fileb://function.zip --handler main  \
  --role arn:aws:iam::566009277786:role/lambda-s3-role-nft
