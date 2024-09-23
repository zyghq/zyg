import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import { createPat } from "@/db/api";
import { useWorkspaceStore } from "@/providers";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  CheckCircledIcon,
  ClipboardCopyIcon,
  ExclamationTriangleIcon,
} from "@radix-ui/react-icons";
import { createFileRoute, Link, useRouterState } from "@tanstack/react-router";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { ArrowLeftIcon } from "lucide-react";
import * as React from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/pats/add"
)({
  component: AddNewPat,
});

type FormInputs = {
  description: string;
  name: string;
};

const formSchema = z.object({
  description: z.string().default(""),
  name: z.string().min(3),
});

function AddNewPat() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();
  const workspaceStore = useWorkspaceStore();
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const { toast } = useToast();
  const [generatedToken, setGeneratedToken] = React.useState("");

  const [copiedText, copyToClipboard] = useCopyToClipboard();
  const hasCopiedText = Boolean(copiedText);

  const form = useForm<FormInputs>({
    defaultValues: {
      description: "",
      name: "",
    },
    resolver: zodResolver(formSchema),
  });

  const { formState } = form;
  const { errors, isSubmitSuccessful, isSubmitting } = formState;

  const isAdding = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { description, name } = inputs;
      const { data, error } = await createPat(token, {
        description,
        name,
      });
      if (error) {
        console.error(error);
        form.setError("root", {
          message: "Something went wrong. Please try again later.",
          type: "serverError",
        });
        return;
      }
      if (data) {
        const { name, token, ...rest } = data;
        setGeneratedToken(token);
        workspaceStore.getState().addPat({ name, token: "", ...rest });
        toast({
          description: `Token ${name} is now added.`,
        });
      }
    } catch (err) {
      console.error(err);
      form.setError("root", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    }
  };

  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="my-12">
            <header className="text-xl font-semibold">
              Add new Personal Access Token
            </header>
          </div>
          <Separator />
        </div>
        <Form {...form}>
          <form className="space-y-8" onSubmit={form.handleSubmit(onSubmit)}>
            {errors?.root?.type === "serverError" && (
              <Alert variant="destructive">
                <ExclamationTriangleIcon className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>
                  {`${errors?.root?.message || "Please try again later."}`}
                </AlertDescription>
              </Alert>
            )}
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input
                      {...field}
                      autoComplete="off"
                      disabled={isSubmitSuccessful}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea
                      className="resize-none"
                      {...field}
                      autoComplete="off"
                      disabled={isSubmitSuccessful}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {isSubmitSuccessful && (
              <div className="flex my-auto">
                <div className="flex flex-col space-y-1">
                  <code className="mr-1 rounded-lg border bg-muted p-1">
                    {generatedToken}
                  </code>
                  <div className="text-sm">
                    {`Please copy this generated token. It won't be shown again.`}
                  </div>
                </div>
                <Button
                  aria-disabled={hasCopiedText}
                  disabled={hasCopiedText}
                  onClick={() => copyToClipboard(generatedToken)}
                  size="sm"
                  type="button"
                  variant="ghost"
                >
                  {hasCopiedText ? (
                    <CheckCircledIcon className="h-5 w-5 text-green-500" />
                  ) : (
                    <ClipboardCopyIcon className="h-4 w-4" />
                  )}
                </Button>
              </div>
            )}
            <div className="flex space-x-4">
              {isSubmitSuccessful ? (
                <Button asChild>
                  <Link
                    params={{ workspaceId }}
                    to={`/workspaces/${workspaceId}/settings/pats`}
                  >
                    <ArrowLeftIcon className="h-4 w-4 mr-1" />
                    Yes, Copied
                  </Link>
                </Button>
              ) : (
                <Button
                  aria-disabled={isAdding}
                  aria-label="Add new Personal Access Token"
                  disabled={isAdding}
                  type="submit"
                >
                  Add Token
                </Button>
              )}
            </div>
          </form>
        </Form>
      </div>
    </div>
  );
}
