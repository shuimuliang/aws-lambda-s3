# Generate thumbnail on AWS

**Table of Contents**

- [Generate thumbnail on AWS](#)
  - [Workflow](#workflow)
  - [Setup AWS S3 bucket](#)
    - [Setup S3 subdirectory](#)
    - [Setup S3 Properities](#)
  - [Setup AWS SNS Topic for S3 origin image put](#)
  - [Setup AWS SNS Subscription for topic](#)
  - [Setup AWS Lambda function](#)
    - [Create AWS IAM Policy and Role](#)
    - [Create AWS Lambda function](#)
    - [Build AWS Lambda function](#)
    - [Upload AWS Lambda function](#)
  - [Upload Demo](#)
  - [Setup AWS SNS Topic for S3 thumb image put](#)
  - [Setup S3 Properities nftavatarthumb ](#)
  - [Setup AWS SNS Subscription for nftavatarthumb topic](#)

## Workflow

![workflow](./000workflow.jpg)

## Setup AWS S3 bucket

创建S3 bucket类似: nftavatar.domain.com

|字段|数值|说明|
|:-----|:------|:-----------------------------|
|bucketname | nftavatar.ascendex.io | |

### Setup S3 subdirectory
创建子目录: origin

|字段|数值|说明|
|:-----|:------|:-----------------------------|
|sub directory | origin | 存放上传的原始文件 |

### Setup S3 Properities
创建S3事件
|字段|数值|说明|
|:-----|:------|:-----------------------------|
|名称| nftavatarput |  |
|事件类型| 发送 | s3:ObjectCreated:Put |
|筛选条件| origin/ | origin目录上传的文件才进行转换 |
|目标类型| SNS主题|  |
|目标| nftavatarput | SNS Topic创建后绑定 |

![S3Properity](./S3Properity.png)

## Setup AWS SNS Topic for S3 origin image put
创建AWS SNS主题
|字段|数值|说明|
|:-----|:------|:-----------------------------|
|topic name | nftavatarput| |
|ARN | arn:aws:sns:ap-northeast-1:566009277786:nftavatarput | 需要修改主题所有者<566009277786> |

访问策略json, 需要修改主题所有者<566009277786>
```json
{
  "Version": "2008-10-17",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__default_statement_ID",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": [
        "SNS:Publish",
        "SNS:RemovePermission",
        "SNS:SetTopicAttributes",
        "SNS:DeleteTopic",
        "SNS:ListSubscriptionsByTopic",
        "SNS:GetTopicAttributes",
        "SNS:Receive",
        "SNS:AddPermission",
        "SNS:Subscribe"
      ],
      "Resource": "arn:aws:sns:ap-northeast-1:566009277786:nftavatarput",
      "Condition": {
        "StringEquals": {
          "AWS:SourceOwner": "566009277786"
        }
      }
    },
    {
      "Sid": "s3",
      "Effect": "Allow",
      "Principal": {
        "Service": "s3.amazonaws.com"
      },
      "Action": "SNS:Publish",
      "Resource": "arn:aws:sns:ap-northeast-1:566009277786:nftavatarput"
    }
  ]
}
```

## Setup AWS SNS Subscription for topic

创建AWS 主题订阅
|字段|数值|说明|
|:-----|:------|:-----------------------------|
|终端节点| arn:aws:lambda:ap-northeast-1:566009277786:function:ThumbnailGo | 修改主题所有者 |
|协议 | LAMBDA | |

## Setup AWS Lambda function

### Create IAM Policy and Role
创建IAM策略: AWSLambdaS3PolicyNFT

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "logs:CreateLogStream",
                "logs:CreateLogGroup",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:*:*:*",
                "arn:aws:s3:::nftavatar.ascendex.io/*"
            ]
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::nftavatar.ascendex.io/*"
        }
    ]
}
```

创建IAM Role:lambda-s3-role-nft, 并附加策略: AWSLambdaS3PolicyNFT

其信任关系json为
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": [
                    "s3.amazonaws.com",
                    "lambda.amazonaws.com"
                ]
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```

基于git https://github.com/gdm-exchange/act-aws-lambda-s3.git
路径 lambda/ThumbnailGo 

### Create AWS Lambda function
|字段|数值|说明|
|:-----|:------|:-----------------------------|
|function-name| ThumbnailGo ||
|runtime| go1.x||
|handler|main||
|role |arn:aws:iam::566009277786:role/lambda-s3-role-nft|需修改所有者56xxxx|
```
aws lambda --profile test create-function --function-name ThumbnailGo --runtime go1.x \
  --zip-file fileb://function.zip --handler main  \
  --role arn:aws:iam::566009277786:role/lambda-s3-role-nft
```

Lambda配置 -> 环境变量

|字段|数值|说明|
|:-----|:------|:-----------------------------|
|DestinationBucket | nftavatar.ascendex.io| |
|ThumbSets | 96x96,48x48| |

### Build AWS Lambda function

```sh
GOOS=linux go build -ldflags "-s -w" main.go
zip function.zip main
```

### Upload AWS Lambda function
```sh
aws lambda --profile <devops-profile-name> update-function-code \
    --function-name  ThumbnailGo \
    --zip-file fileb://function.zip
```

## Upload Demo

1. 通过网页，或者服务端SDK，在origin目录上传cat-program.jpeg
2. 可以通过CloudWatch > Log groups > /aws/lambda/ThumbnailGo 看到日志事件，通过s3 cmd看到缩略图成功生成
```sh
aws s3 ls s3://nftavatar.ascendex.io/ --recursive --profile <devops-profile-name>
```

![Lambda-demo](./lambda-demo.jpg)

### Setup AWS SNS Topic for S3 converted image put
创建AWS SNS主题
|字段|数值|说明|
|:-----|:------|:-----------------------------|
|topic name | nftavatarconvert | |
|ARN | arn:aws:sns:ap-northeast-1:566009277786:nftavatarconvert | 需要修改主题所有者<566009277786> |

访问策略json, 需要修改主题所有者<566009277786>
```json
{
  "Version": "2008-10-17",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__default_statement_ID",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": [
        "SNS:Publish",
        "SNS:RemovePermission",
        "SNS:SetTopicAttributes",
        "SNS:DeleteTopic",
        "SNS:ListSubscriptionsByTopic",
        "SNS:GetTopicAttributes",
        "SNS:Receive",
        "SNS:AddPermission",
        "SNS:Subscribe"
      ],
      "Resource": "arn:aws:sns:ap-northeast-1:566009277786:nftavatarthumb",
      "Condition": {
        "StringEquals": {
          "AWS:SourceOwner": "566009277786"
        }
      }
    },
    {
      "Sid": "s3",
      "Effect": "Allow",
      "Principal": {
        "Service": "s3.amazonaws.com"
      },
      "Action": "SNS:Publish",
      "Resource": "arn:aws:sns:ap-northeast-1:566009277786:nftavatarthumb"
    }
  ]
}
```

### Setup S3 Properities nftavatarthumb
创建S3事件
|字段|数值|说明|
|:-----|:------|:-----------------------------|
|名称| nftavatarthumb |  |
|事件类型| 发送 | s3:ObjectCreated:Put |
|筛选条件| 48x48/ | 选取其中一个转换目录 |
|目标类型| SNS主题|  |
|目标| nftavatarthumb | SNS Topic创建后绑定 |

### Setup AWS SNS Subscription for nftavatarthumb topic
创建SNS订阅
|字段|数值|说明|
|:-----|:------|:-----------------------------|
|主题ARN| arn:aws:sns:ap-northeast-1:566009277786:nftavatarthumb |  |
|协议| HTTP | |
|终端节点endpoint| http://35.78.198.128:1230/v1/sns/s3/callback | 域名,端口,uri需要按生产配置 |

SNS第一个payload是, 该payload用于确认终端节点可以联通
```json
{"uuid": 
   "{
    "Type":"SubscriptionConfirmation",
    "MessageId" : "62997681-c3c9-49ed-b8d2-0c605f3c3f4a",
    "Token":"2336412f37fb687f5d51e6e2425dacbbaa2960518b60ee5f7522150f536585e075e65b50164cb1213c2820332307304ed08302ee3e20c2fb5ee7d755ab3ab5ff78e0f13cb4375d39ee6437f94ab5a8acfdc1372539113f9222c56291dfade5cee0eed2ce10edd651861e247b161a73f7fc00a02114270dfb59b76bd559f35063",
    "TopicArn":"arn:aws:sns:ap-northeast-1:566009277786:nftavatarthumb",
    "Message":"You have chosen to subscribe to the topic arn:aws:sns:ap-northeast-1:566009277786:nftavatarthumb.To confirm the subscription, visit the SubscribeURL included in this message.",
    "SubscribeURL":"https://sns.ap-northeast-1.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:ap-northeast-1:566009277786:nftavatarthumb&Token=2336412f37fb687f5d51e6e2425dacbbaa2960518b60ee5f7522150f536585e075e65b50164cb1213c2820332307304ed08302ee3e20c2fb5ee7d755ab3ab5ff78e0f13cb4375d39ee6437f94ab5a8acfdc1372539113f9222c56291dfade5cee0eed2ce10edd651861e247b161a73f7fc00a02114270dfb59b76bd559f35063",
    "Timestamp":"2022-05-05T20:34:05.377Z",
    "SignatureVersion":"1",
 "Signature":"OEftVWiPZbNkpT+kr+ckqxxRFZKnNwG8rn69ZzqjHY/q6nNzojbiE6ZHyTzYd0MJw/tDx0l5b3NMEHTt4JSNTr+RUtZ8Mim4pGGl9dMRKOmdQLeDw9WTzg4B0yzEymr/If4S18gugFLtIQzO1UlnLvUcDLfDUKcipuvJ/OSHrWGXhU2o78uFJ9C3xdQzragpNaudPzrXwPHV70quR9XaPqxBsod5xcTQO56+v+3ldEl3SH3+J/b2RsfSRaCgsMzql011NkuU0YX54KFjsBJ2dwedtn9Sa9iNhVI7VXAM2OjLBxA++ww96axxpvloo4nbgh55THOcj7OKofJJCw54Lw==",
    "SigningCertURL":"https://sns.ap-northeast-1.amazonaws.com/SimpleNotificationService-7ff5318490ec183fbaddaa2a969abfda.pem"
   }"
}
```

endpoint服务器, 需要
1) 放行http://domain/v1/sns/s3/callback这个路径, 接受SNS第一个payload
2) 使用json库反序列化, 拿到SubscribeURL
3) 使用HTTP GET请求一次SubscribeURL
4) 返回HTTP 200

SNS第二个payload是
```json
{
   "uuid":"{""Type"" : ""Notification"",""MessageId"" : "34f1ee0f-60b7-5920-87c9-b4e1598e5f34",""TopicArn"" : ""arn":"aws":"sns":"ap-northeast-1":"566009277786":"avatarthumb"",""Subject"" : "Amazon S3 Notification",""Message"" : "{
      "Records":[
         {
            "eventVersion":"2.1",
            "eventSource":"aws:s3",
            "awsRegion":"ap-northeast-1",
            "eventTime":"2022-05-05T20:43:36.316Z",
            "eventName":"ObjectCreated:Put",
            "userIdentity":{
               "principalId":"A14QFIMDECLDKR"
            },
            "requestParameters":{
               "sourceIPAddress":"10.6.1.214"
            },
            "responseElements":{
               "x-amz-request-id":"WPVSR9HM0BSS649R",
               "x-amz-id-2":"zfO8ffv8bKKs4/Qlzzmc+ThWrwH8OcWpuQub9OzxRpuITYWTOaQnX8d972hWxMxEEeW7T7Z6Ke14GaYAFt8f22G5qSKrgex6HgOWw5mO1lk="
            },
            "s3":{
               "s3SchemaVersion":"1.0",
               "configurationId":"avatarthumb",
               "bucket":{
                  "name":"avatar.movieous.video",
                  "ownerIdentity":{
                     "principalId":"A14QFIMDECLDKR"
                  },
                  "arn":"arn:aws:s3:::avatar.movieous.video"
               },
               "object":{
                  "key":"48x48/chicken-beef.png",
                  "size":93268,
                  "eTag":"dd09fc3edfa0d3f11d97fad84b20c385",
                  "sequencer":"00627436F841ACF414"
               }
            }
         }
      ]
   }",""Timestamp"" : ""2022-05-05T20":"43":37.363Z",""SignatureVersion"" : "1",""Signature"" : "Awu4o9PRlO5frOC62XIKF72ffDc/C3o0CKrwy+uavd+FXs6NLGqpQ0yG9okAUgZOB7BbBDYvtdfTiNgXpz00qCUUNdI6P3Hwh/ittDm/y2yzaVxoc1uYzpNZMmcbruRHLNQf0eieqM8ppd1Li+hiiWq1yeLW2RC4xFZ4ShixzMKW+SsLq6GDiWsJEbM5hMNOQuAqF4egnKHwOICj17hZz+/RT0kLb2V6Hynt+EHsHUu+VgVxbO/WuKhOjvEv86Gby0c/IuP8DMaNO3FbUPP8wwcaIaQKqrReubd36JjyAGG2wtQ+ZIoOX+m8Ic5URwE0fX2thTvp32mYkIH5AQZuLw==",""SigningCertURL"" : ""https"://sns.ap-northeast-1.amazonaws.com/SimpleNotificationService-7ff5318490ec183fbaddaa2a969abfda.pem",""UnsubscribeURL"" : ""https"::"aws":"sns":"ap-northeast-1":"566009277786":"avatarthumb":c5c337a5-f819-4077-993e-9caa573f5564"n}"
}
```

endpoint服务器需要
1) 拿到返回的json payload
2) 提取s3 key中的bucket name和key
3) 在服务端标记该key所对应的文件已转换成功, 比如状态0(未上传), 状态1(已上传), 状态2(已生成缩略图)
4) 在服务端存一张映射表，记录s3 bucket name和CDN域名的映射关系, 用于拼接图片的CDN路径

|字段|数值|说明|
|:-----|:------|:-----------------------------|
|AWS S3 bucket name| avatarnft.ascendex.io | |
|AWS CloudFront加速域名| avatarnft.ascendex.io | 配成同一个比较省事 |
