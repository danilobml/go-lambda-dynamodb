#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { GoLambdaDynamoStack } from '../lib/go-lambda-dynamo-stack';

const app = new cdk.App();
new GoLambdaDynamoStack(app, 'GoLambdaDynamoStack', {
});