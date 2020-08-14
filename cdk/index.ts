import events = require('@aws-cdk/aws-events');
import targets = require('@aws-cdk/aws-events-targets');
import lambda = require('@aws-cdk/aws-lambda');
import cdk = require('@aws-cdk/core');

require('dotenv').config();

const {
  ENDPOINT = "",
  WEBHOOK = "",
  FROM = "9",
  TO = "2",
  DATE = "2020-09-01",
  TIMETABLE = "18:00",
  TICKETCOUNT = "2",
  CARRIAGECATEGORY = "0",
  ONLYSHOWDISCOUNT = "0",
  COLLEGESTUDENTS = "0",
  DEVICEID = "",
  DEVICEIDHASH = "",
  DEVICECATEGORY = "I",
  APPVERSION = "5.20",
  PARAMETERVERSION = "20200527"
} = process.env;

export class LambdaCronStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const env = {
      ENV: 'production',
      ENDPOINT,
      WEBHOOK,
      FROM,
      TO,
      DATE,
      TIMETABLE,
      TICKETCOUNT,
      CARRIAGECATEGORY,
      ONLYSHOWDISCOUNT,
      COLLEGESTUDENTS,
      DEVICEID,
      DEVICEIDHASH,
      DEVICECATEGORY,
      APPVERSION,
      PARAMETERVERSION
    };
    console.log(env)
    const lambdaFn = new lambda.Function(this, 'Singleton', {
      code: lambda.Code.fromAsset('./lambda-hsr_reminder.zip'),
      handler: 'lambda-hsr_reminder',
      runtime: lambda.Runtime.GO_1_X,
      environment: env,
      timeout: cdk.Duration.seconds(5),
      retryAttempts: 0
    });

    const rule = new events.Rule(this, 'Rule', {
      schedule: events.Schedule.expression('rate(10 minutes)')
    });

    rule.addTarget(new targets.LambdaFunction(lambdaFn));
  }
}

const app = new cdk.App();
new LambdaCronStack(app, 'hsr-reminder-go');
app.synth();
