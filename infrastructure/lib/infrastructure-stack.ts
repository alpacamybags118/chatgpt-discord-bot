import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigatewayv2 from '@aws-cdk/aws-apigatewayv2';
import * as integrations from '@aws-cdk/aws-apigatewayv2-integrations';
import { Duration } from 'aws-cdk-lib';



export class InfrastructureStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const discordLambda = new lambda.Function(this, "discordbotfunction", {
      handler: "main",
      code: lambda.Code.fromAsset("./assets/main.zip"),
      runtime: lambda.Runtime.GO_1_X,
      environment: {
        DISCORD_PUBLIC_KEY: process.env.DISCORD_PUBLIC_KEY || '',
        DISCORD_TOKEN: process.env.DISCORD_TOKEN || '',
        OPEN_AI_API_KEY: process.env.OPEN_AI_API_KEY || '',
      },
      timeout: cdk.Duration.seconds(15),
    })

    const apiGw = new apigatewayv2.HttpApi(this, "discordbotapigw", {
      corsPreflight: {
        allowOrigins: ["*"]
      }
    });

    const lambdaIntegration = new integrations.HttpLambdaIntegration("discordbotapigw", discordLambda);
  
    apiGw.addRoutes({
      path: '/discord',
      methods: [apigatewayv2.HttpMethod.ANY],
      integration: lambdaIntegration,
      
    });

  }
}
