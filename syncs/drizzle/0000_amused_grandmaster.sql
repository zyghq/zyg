-- Current sql file was generated after introspecting the database
-- If you want to run this migration please uncomment this code before executing migrations
/*
CREATE TABLE "workspace" (
	"workspace_id" varchar(255) PRIMARY KEY NOT NULL,
	"name" varchar(255) NOT NULL,
	"public_name" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"version_id" uuid NOT NULL,
	"synced_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
CREATE TABLE "member" (
	"member_id" varchar(255) PRIMARY KEY NOT NULL,
	"workspace_id" varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	"public_name" varchar(255) NOT NULL,
	"role" varchar(255) NOT NULL,
	"permissions" jsonb NOT NULL,
	"avatar_url" varchar(511) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"version_id" uuid NOT NULL,
	"synced_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
CREATE TABLE "customer" (
	"customer_id" varchar(255) PRIMARY KEY NOT NULL,
	"workspace_id" varchar(255) NOT NULL,
	"external_id" varchar(255),
	"email" varchar(255),
	"phone" varchar(255),
	"name" varchar(255) NOT NULL,
	"role" varchar(255) NOT NULL,
	"avatar_url" varchar(511) NOT NULL,
	"is_email_verified" boolean DEFAULT false NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"version_id" uuid NOT NULL,
	"synced_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT "customer_workspace_id_external_id_key" UNIQUE("workspace_id","external_id"),
	CONSTRAINT "customer_workspace_id_email_key" UNIQUE("workspace_id","email"),
	CONSTRAINT "customer_workspace_id_phone_key" UNIQUE("workspace_id","phone")
);
--> statement-breakpoint
CREATE TABLE "thread" (
	"thread_id" varchar(255) PRIMARY KEY NOT NULL,
	"workspace_id" varchar(255) NOT NULL,
	"customer_id" varchar(255) NOT NULL,
	"assignee_id" varchar(255),
	"assigned_at" timestamp,
	"title" text NOT NULL,
	"description" text NOT NULL,
	"preview_text" text NOT NULL,
	"status" varchar(127) NOT NULL,
	"status_changed_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"status_changed_by_id" varchar(255) NOT NULL,
	"stage" varchar(127) NOT NULL,
	"replied" boolean DEFAULT false NOT NULL,
	"priority" varchar(255) NOT NULL,
	"channel" varchar(127) NOT NULL,
	"created_by_id" varchar(255) NOT NULL,
	"updated_by_id" varchar(255) NOT NULL,
	"labels" jsonb,
	"inbound_seq_id" varchar(255),
	"outbound_seq_id" varchar(255),
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"version_id" uuid NOT NULL,
	"synced_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
ALTER TABLE "member" ADD CONSTRAINT "member_workspace_id_fkey" FOREIGN KEY ("workspace_id") REFERENCES "public"."workspace"("workspace_id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "customer" ADD CONSTRAINT "customer_workspace_id_fkey" FOREIGN KEY ("workspace_id") REFERENCES "public"."workspace"("workspace_id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "thread" ADD CONSTRAINT "thread_workspace_id_fkey" FOREIGN KEY ("workspace_id") REFERENCES "public"."workspace"("workspace_id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "thread" ADD CONSTRAINT "thread_customer_id_fkey" FOREIGN KEY ("customer_id") REFERENCES "public"."customer"("customer_id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "thread" ADD CONSTRAINT "thread_assignee_id_fkey" FOREIGN KEY ("assignee_id") REFERENCES "public"."member"("member_id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "thread" ADD CONSTRAINT "thread_status_changed_by_id_fkey" FOREIGN KEY ("status_changed_by_id") REFERENCES "public"."member"("member_id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "thread" ADD CONSTRAINT "thread_created_by_id_fkey" FOREIGN KEY ("created_by_id") REFERENCES "public"."member"("member_id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "thread" ADD CONSTRAINT "thread_updated_by_id_fkey" FOREIGN KEY ("updated_by_id") REFERENCES "public"."member"("member_id") ON DELETE no action ON UPDATE no action;
*/