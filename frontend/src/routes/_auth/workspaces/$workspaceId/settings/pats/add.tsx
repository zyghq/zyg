import * as React from "react";
import { createFileRoute, useRouterState, Link } from "@tanstack/react-router";
import { Separator } from "@/components/ui/separator";
import { z } from "zod";
import { useToast } from "@/components/ui/use-toast";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  ExclamationTriangleIcon,
  ClipboardCopyIcon,
  CheckCircledIcon,
} from "@radix-ui/react-icons";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { createAccountPat } from "@/db/api";
import { useWorkspaceStore } from "@/providers";
import { ArrowLeftIcon } from "lucide-react";
import { useCopyToClipboard } from "@uidotdev/usehooks";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/settings/pats/add"
)({
  component: AddNewPat,
});

type FormInputs = {
  name: string;
  description: string;
};

const formSchema = z.object({
  name: z.string().min(3),
  description: z.string().default(""),
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
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      description: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  const isAdding = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { name, description } = inputs;
      const { error, data } = await createAccountPat(token, {
        name,
        description,
      });
      if (error) {
        console.error(error);
        form.setError("root", {
          type: "serverError",
          message: "Something went wrong. Please try again later.",
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
        type: "serverError",
        message: "Something went wrong. Please try again later.",
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
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
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
                  type="button"
                  variant="ghost"
                  size="sm"
                  disabled={hasCopiedText}
                  aria-disabled={hasCopiedText}
                  onClick={() => copyToClipboard(generatedToken)}
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
                    to={`/workspaces/${workspaceId}/settings/pats`}
                    params={{ workspaceId }}
                  >
                    <ArrowLeftIcon className="h-4 w-4 mr-1" />
                    Yes, Copied
                  </Link>
                </Button>
              ) : (
                <Button
                  type="submit"
                  disabled={isAdding}
                  aria-disabled={isAdding}
                  aria-label="Add new Personal Access Token"
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
