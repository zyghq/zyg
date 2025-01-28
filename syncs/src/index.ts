import * as restate from "@restatedev/restate-sdk/fetch";

interface MessageBody {
  email: string;
  messageId: string;
}

const handler = restate
  .endpoint()
  .bind(
    restate.service({
      name: "EmailTemplateService",
      handlers: {
        Message: async (ctx: restate.Context, body: MessageBody) => {
          // Durably execute a set of steps; resilient against failures
          const greetingId = ctx.rand.uuidv4();
          await ctx.sleep(1000);
          console.log(`Generated greeting ID: ${greetingId}`);
          console.log(`Got Email: ${body.email}`);
          console.log(`Got MessageId: ${body.messageId}`);

          // Respond to caller
          return `Got GreetingId ${greetingId}!`;
        },
      },
    })
  )
  .bidirectional()
  .handler();

const server = Bun.serve({
  port: 9080,
  ...handler,
});

console.log(`Listening on ${server.url}`);
