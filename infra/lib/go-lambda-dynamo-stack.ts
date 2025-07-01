import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { RestApi, LambdaIntegration } from 'aws-cdk-lib/aws-apigateway';

export class GoLambdaDynamoStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Lambda function
    const lambdaFunction = new lambda.Function(this, 'GoLambdaDynamoHandler', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset("../lambdas")
    })

    // DynamoDB table
    const dynamoDbPeopleTable = new dynamodb.TableV2(this, "GoLambdaDynamoPeopleTable", {
      tableName: 'people',
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      billing: dynamodb.Billing.onDemand()
    })
    // IAM Policy
    dynamoDbPeopleTable.grantReadWriteData(lambdaFunction)

    // Api Gateway REST api
    const api = new RestApi(this, "GoLambdaDynamoApi", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["OPTIONS", "GET", "POST", "PUT", "DELETE"],
        allowHeaders: ["Content-Type"],
        allowCredentials: true
      }
    })

    // API Proxy
    const integration = new LambdaIntegration(lambdaFunction)
    api.root.addProxy({
      defaultIntegration: integration,
      anyMethod: true
    })
  }
}
