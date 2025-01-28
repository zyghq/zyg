import * as restate from "@restatedev/restate-sdk/fetch";
import { drizzle } from "drizzle-orm/neon-http";
import { eq, sql, and } from "drizzle-orm";
import * as schema from "./db/schema";

const syncdb = drizzle(process.env.DATABASE_URL as string);

// Defines a service to check if the in-sync service is up and running
const inSync = restate.service({
  name: "inSync",
  handlers: {
    ping: async (ctx: restate.Context) => {
      console.log(`got ping...`);
      const pingId = ctx.rand.uuidv4();
      console.log(`Generated ping ID: ${pingId}`);
      return {
        message: "pong",
        pingId,
      };
    },
    dbtime: async (ctx: restate.Context) => {
      console.log(`invoke dbtime...`);
      const time = await syncdb.execute(sql`select now()`);
      console.log(`Got time: ${time.rows[0].now}`);
      return {
        time: time.rows[0].now,
      };
    },
  },
});

type InsertWorkspace = typeof schema.workspace.$inferInsert;
type InsertMember = typeof schema.member.$inferInsert;
type InsertCustomer = typeof schema.customer.$inferInsert;
type InsertThread = typeof schema.thread.$inferInsert;

const sync = restate.object({
  name: "sync",
  handlers: {
    upsertWorkspace: async (
      ctx: restate.ObjectContext,
      workspace: InsertWorkspace,
    ): Promise<{
      workspaceId: string;
      versionId: string;
    }> => {
      console.log("invoking upsertWorkspace...");
      const currVersionId = await ctx.run(
        "read current workspace version",
        async () => {
          // first read the current member version from the database.
          const results = await syncdb
            .select({ versionId: schema.workspace.versionId })
            .from(schema.workspace)
            .where(eq(schema.workspace.workspaceId, workspace.workspaceId));
          // if the workspace does not exist, version is null
          if (results.length !== 1) {
            return null;
          }
          // found the workspace with version ID return it.
          return results[0].versionId;
        },
      );
      if (!currVersionId) {
        // workspace does not exist, insert new with initial version
        const newVersionId = ctx.rand.uuidv4();
        const inserts = {
          ...workspace,
          versionId: newVersionId,
        } as InsertWorkspace;
        const dbInserted: { workspaceId: string; versionId: string }[] =
          await syncdb
            .insert(schema.workspace)
            .values(inserts)
            .onConflictDoUpdate({
              target: schema.workspace.workspaceId,
              set: { ...inserts },
            })
            .returning({
              workspaceId: schema.workspace.workspaceId,
              versionId: schema.workspace.versionId,
            });
        if (dbInserted.length !== 1) {
          throw new restate.TerminalError(
            `Failed to insert workspace with workspace ID ${workspace.workspaceId} version ID ${newVersionId}`,
          );
        }
        const inserted = dbInserted[0];
        const { workspaceId, versionId } = inserted;
        console.log(
          `Inserted workspace with workspace ID ${workspaceId} version ID ${versionId}`,
        );
        return { workspaceId: workspaceId, versionId: versionId };
      }
      // apply update to workspace with current version, updating to new version.
      const newVersionId = ctx.rand.uuidv4();
      const updates = {
        ...workspace,
        versionId: newVersionId,
      } as InsertWorkspace;
      const dbUpdated: { workspaceId: string; versionId: string }[] =
        await syncdb
          .update(schema.workspace)
          .set({ ...updates })
          .where(
            and(
              eq(schema.workspace.workspaceId, workspace.workspaceId),
              eq(schema.workspace.versionId, currVersionId),
            ),
          )
          .returning({
            workspaceId: schema.workspace.workspaceId,
            versionId: schema.workspace.versionId,
          });
      if (dbUpdated.length !== 1) {
        throw new restate.TerminalError(
          `Failed to update workspace with workspace ID ${workspace.workspaceId} version ID ${newVersionId}`,
        );
      }
      const updated = dbUpdated[0];
      const { workspaceId, versionId } = updated;
      console.log(
        `Updated workspace with workspace ID ${workspaceId} version ID ${versionId}`,
      );
      return { workspaceId: workspaceId, versionId: versionId };
    },
    upsertMember: async (
      ctx: restate.ObjectContext,
      member: InsertMember,
    ): Promise<{
      memberId: string;
      versionId: string;
    }> => {
      console.log("invoking upsertMember...");
      const currVersionId = await ctx.run(
        "read current member version",
        async () => {
          // first read the current member version from the database.
          const results = await syncdb
            .select({ versionId: schema.member.versionId })
            .from(schema.member)
            .where(eq(schema.member.memberId, member.memberId));
          // if the member does not exist, version is null
          if (results.length !== 1) {
            return null;
          }
          // found the member with version ID return it.
          return results[0].versionId;
        },
      );
      if (!currVersionId) {
        // member does not exist, insert new with initial version
        const newVersionId = ctx.rand.uuidv4();
        const inserts = { ...member, versionId: newVersionId } as InsertMember;
        const upsertedMember: { memberId: string; versionId: string }[] =
          await syncdb
            .insert(schema.member)
            .values(inserts)
            .onConflictDoUpdate({
              target: schema.member.memberId,
              set: { ...inserts },
            })
            .returning({
              memberId: schema.member.memberId,
              versionId: schema.member.versionId,
            });
        if (upsertedMember.length !== 1) {
          throw new restate.TerminalError(
            `Failed to insert member with member ID ${member.memberId} version ID ${newVersionId}`,
          );
        }
        console.log(
          `Inserted member with member ID ${member.memberId} version ID ${newVersionId}`,
        );
        return { memberId: member.memberId, versionId: newVersionId };
      }
      // apply update to member with current version, updating to new version.
      const newVersionId = ctx.rand.uuidv4();
      const updates = { ...member, versionId: newVersionId } as InsertMember;
      const updatedMember: { memberId: string; versionId: string }[] =
        await syncdb
          .update(schema.member)
          .set({ ...updates })
          .where(
            and(
              eq(schema.member.memberId, member.memberId),
              eq(schema.member.versionId, currVersionId),
            ),
          )
          .returning({
            memberId: schema.member.memberId,
            versionId: schema.member.versionId,
          });
      if (updatedMember.length !== 1) {
        throw new restate.TerminalError(
          `Failed to update member with member ID ${member.memberId} version ID ${newVersionId}`,
        );
      }
      console.log(
        `Updated member with member ID ${member.memberId} version ID ${newVersionId}`,
      );
      return { memberId: member.memberId, versionId: newVersionId };
    },
    upsertCustomer: async (
      ctx: restate.ObjectContext,
      customer: InsertCustomer,
    ): Promise<{
      customerId: string;
      versionId: string;
    }> => {
      console.log("invoking upsertCustomer...");
      const currVersionId = await ctx.run(
        "read current customer version",
        async () => {
          // first read the current member version from the database.
          const results = await syncdb
            .select({ versionId: schema.customer.versionId })
            .from(schema.customer)
            .where(eq(schema.customer.customerId, customer.customerId));
          // if the customer does not exist, version is null
          if (results.length !== 1) {
            return null;
          }
          // found the customer with version ID return it.
          return results[0].versionId;
        },
      );
      if (!currVersionId) {
        // customer does not exist, insert new with initial version
        const newVersionId = ctx.rand.uuidv4();
        const inserts = {
          ...customer,
          versionId: newVersionId,
        } as InsertCustomer;
        const upsertedCustomer: { customerId: string; versionId: string }[] =
          await syncdb
            .insert(schema.customer)
            .values(inserts)
            .onConflictDoUpdate({
              target: schema.customer.customerId,
              set: { ...inserts },
            })
            .returning({
              customerId: schema.customer.customerId,
              versionId: schema.customer.versionId,
            });
        if (upsertedCustomer.length !== 1) {
          throw new restate.TerminalError(
            `Failed to insert customer with customer ID ${customer.customerId} version ID ${newVersionId}`,
          );
        }
        console.log(
          `Inserted customer with customer ID ${customer.customerId} version ID ${newVersionId}`,
        );
        return { customerId: customer.customerId, versionId: newVersionId };
      }
      // apply update to customer with current version, updating to new version.
      const newVersionId = ctx.rand.uuidv4();
      const updates = {
        ...customer,
        versionId: newVersionId,
      } as InsertCustomer;
      const updatedCustomer: { customerId: string; versionId: string }[] =
        await syncdb
          .update(schema.customer)
          .set({ ...updates })
          .where(
            and(
              eq(schema.customer.customerId, customer.customerId),
              eq(schema.customer.versionId, currVersionId),
            ),
          )
          .returning({
            customerId: schema.customer.customerId,
            versionId: schema.customer.versionId,
          });
      if (updatedCustomer.length !== 1) {
        throw new restate.TerminalError(
          `Failed to update customer with customer ID ${customer.customerId} version ID ${newVersionId}`,
        );
      }
      console.log(
        `Updated customer with customer ID ${customer.customerId} version ID ${newVersionId}`,
      );
      return { customerId: customer.customerId, versionId: newVersionId };
    },
    upsertThread: async (
      ctx: restate.ObjectContext,
      thread: InsertThread,
    ): Promise<{
      threadId: string;
      versionId: string;
    }> => {
      console.log("invoking upsertThread...");
      const currVersionId = await ctx.run(
        "read current thread version",
        async () => {
          // first read the current thread version from the database.
          const results = await syncdb
            .select({ versionId: schema.thread.versionId })
            .from(schema.thread)
            .where(eq(schema.thread.threadId, thread.threadId));
          // if the thread does not exist, version is null
          if (results.length !== 1) {
            return null;
          }
          // found the thread with version ID return it.
          return results[0].versionId;
        },
      );
      if (!currVersionId) {
        // thread does not exist, insert new with initial version
        const newVersionId = ctx.rand.uuidv4();
        const inserts = {
          ...thread,
          versionId: newVersionId,
        } as InsertThread;
        const upsertedThread: { threadId: string; versionId: string }[] =
          await syncdb
            .insert(schema.thread)
            .values(inserts)
            .onConflictDoUpdate({
              target: schema.thread.threadId,
              set: { ...inserts },
            })
            .returning({
              threadId: schema.thread.threadId,
              versionId: schema.thread.versionId,
            });
        if (upsertedThread.length !== 1) {
          throw new restate.TerminalError(
            `Failed to insert thread with thread ID ${thread.threadId} version ID ${newVersionId}`,
          );
        }
        console.log(
          `Inserted thread with thread ID ${thread.threadId} version ID ${newVersionId}`,
        );
        return { threadId: thread.threadId, versionId: newVersionId };
      }
      // apply update to thread with current version, updating to new version.
      const newVersionId = ctx.rand.uuidv4();
      const updates = {
        ...thread,
        versionId: newVersionId,
      } as InsertThread;
      const updatedThread: { threadId: string; versionId: string }[] =
        await syncdb
          .update(schema.thread)
          .set({ ...updates })
          .where(
            and(
              eq(schema.thread.threadId, thread.threadId),
              eq(schema.thread.versionId, currVersionId),
            ),
          )
          .returning({
            threadId: schema.thread.threadId,
            versionId: schema.thread.versionId,
          });
      if (updatedThread.length !== 1) {
        throw new restate.TerminalError(
          `Failed to update thread with thread ID ${thread.threadId} version ID ${newVersionId}`,
        );
      }
      console.log(
        `Updated thread with thread ID ${thread.threadId} version ID ${newVersionId}`,
      );
      return { threadId: thread.threadId, versionId: newVersionId };
    },
  },
});

const handler = restate
  .endpoint()
  .bind(inSync)
  .bind(sync)
  .bidirectional()
  .handler();

const server = Bun.serve({
  port: 9080,
  ...handler,
});

console.log(`sync service listening on ${server.url}`);
