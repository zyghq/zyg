import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { createWorkspace } from "@/db/api";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRightIcon, ExclamationTriangleIcon } from "@radix-ui/react-icons";
import {
  createFileRoute,
  Link,
  useNavigate,
  useRouter,
  useRouterState,
} from "@tanstack/react-router";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";

export const Route = createFileRoute("/_account/workspaces/add")({
  component: CreateWorkspaceComponent,
});

type FormInputs = {
  name: string;
};

const formSchema = z.object({
  name: z.string().min(3),
});

function CreateWorkspaceComponent() {
  const { token } = Route.useRouteContext();
  const router = useRouter();
  const navigate = useNavigate();
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const { toast } = useToast();

  const form = useForm({
    defaultValues: {
      name: "",
    },
    resolver: zodResolver(formSchema),
  });

  const { formState } = form;
  const { errors, isSubmitSuccessful, isSubmitting } = formState;

  const isCreating = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { name } = inputs;
      const { data, error } = await createWorkspace(token, { name });
      if (error) {
        console.error(error);
        form.setError("root", {
          message: "Something went wrong. Please try again later.",
          type: "serverError",
        });
        return;
      }
      if (data) {
        toast({
          description: `Workspace ${data.workspaceName} created.`,
        });
        await router.invalidate();
        await navigate({
          params: { workspaceId: data.workspaceId },
          to: "/workspaces/$workspaceId",
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
    <div className="flex min-h-screen flex-col justify-center p-4">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <Card className="mx-auto w-full max-w-sm">
            <CardHeader>
              <CardTitle>Create a new Workspace.</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
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
                      <Input type="text" {...field} required />
                    </FormControl>
                    <FormDescription>
                      Typically your company name.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
            <CardFooter className="flex justify-between">
              <Button
                aria-disabled={isCreating || isSubmitSuccessful}
                aria-label="Create Workspace"
                disabled={isCreating || isSubmitSuccessful}
                type="submit"
              >
                Save & Continue
                <ArrowRightIcon className="ml-2 h-4 w-4" />
              </Button>
            </CardFooter>
            <CardFooter className="flex justify-center">
              <Button asChild variant="link">
                <Link to="/workspaces">Go Back</Link>
              </Button>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </div>
  );
}
