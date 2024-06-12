import {
  createFileRoute,
  useNavigate,
  useRouter,
  useRouterState,
  Link,
} from "@tanstack/react-router";
import { useToast } from "@/components/ui/use-toast";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ExclamationTriangleIcon, ArrowRightIcon } from "@radix-ui/react-icons";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import { z } from "zod";
import { createWorkspace } from "@/db/api";

export const Route = createFileRoute("/_auth/workspaces/add")({
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
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  const isCreating = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { name } = inputs;
      const { error, data } = await createWorkspace(token, { name });
      if (error) {
        console.error(error);
        form.setError("root", {
          type: "serverError",
          message: "Something went wrong. Please try again later.",
        });
        return;
      }
      if (data) {
        toast({
          description: `Workspace ${data.workspaceName} created.`,
        });
        await router.invalidate();
        await navigate({
          to: "/workspaces/$workspaceId",
          params: { workspaceId: data.workspaceId },
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
    <div className="flex flex-col justify-center p-4 min-h-screen">
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
                type="submit"
                disabled={isCreating || isSubmitSuccessful}
                aria-disabled={isCreating || isSubmitSuccessful}
                aria-label="Create Workspace"
              >
                Save & Continue
                <ArrowRightIcon className="ml-2 h-4 w-4" />
              </Button>
            </CardFooter>
            <CardFooter className="flex justify-center">
              <Button variant="link" asChild>
                <Link to="/workspaces">Go Back</Link>
              </Button>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </div>
  );
}
