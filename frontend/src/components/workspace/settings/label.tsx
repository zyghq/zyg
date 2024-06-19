import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Pencil1Icon, DotsHorizontalIcon } from "@radix-ui/react-icons";
import { TagIcon } from "lucide-react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useMutation } from "@tanstack/react-query";
import { updateWorkspaceLabel } from "@/db/api";
import { useWorkspaceStore } from "@/providers";

type FormInputs = {
  name: string;
};

const formSchema = z.object({
  name: z.string().min(2),
});

export function LabelItem({
  token,
  workspaceId,
  labelId,
  label,
}: {
  token: string;
  workspaceId: string;
  labelId: string;
  label: string;
}) {
  const [editMode, setEditMode] = React.useState(false);
  const workspaceStore = useWorkspaceStore();

  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: label,
    },
  });

  const mutation = useMutation({
    mutationFn: async (inputs: { name: string }) => {
      const { name } = inputs;
      const { error, data } = await updateWorkspaceLabel(
        token,
        workspaceId,
        labelId,
        {
          name,
        }
      );
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data;
    },
    onError: (error) => {
      console.error(error);
      form.setError("name", {
        type: "serverError",
        message: "Something went wrong. Please try again later.",
      });
    },
    onSuccess: (data) => {
      const { labelId, ...rest } = data;
      workspaceStore.getState().updateLabel(labelId, { labelId, ...rest });
      form.reset();
      setEditMode(false);
    },
  });

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    await mutation.mutateAsync(inputs);
  };

  const { formState } = form;
  const { isSubmitting } = formState;

  const isUpdating = isSubmitting || mutation.isPending;

  return (
    <div className="flex flex-col items-start gap-2 rounded-lg border p-3">
      {editMode ? (
        <div className="flex w-full flex-col">
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="flex justify-between space-x-2"
            >
              <div className="flex items-center w-full">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem className="w-full">
                      <FormControl>
                        <Input
                          autoFocus
                          required
                          autoComplete="off"
                          placeholder="Label"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <div className="flex gap-1">
                <Button
                  disabled={isUpdating}
                  type="button"
                  variant="outline"
                  onClick={() => setEditMode(false)}
                  size="default"
                >
                  Cancel
                </Button>
                <Button disabled={isUpdating} type="submit" size="default">
                  Save
                </Button>
              </div>
            </form>
          </Form>
        </div>
      ) : (
        <div className="flex w-full flex-col">
          <div className="flex items-center">
            <div className="flex items-center gap-2">
              <TagIcon className="h-4 w-4 text-muted-foreground" />
              <div className="font-normal">{label}</div>
            </div>
            <div className="ml-auto">
              <Button
                onClick={() => setEditMode(true)}
                variant="ghost"
                size="icon"
              >
                <Pencil1Icon className="h-4 w-4" />
              </Button>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <DotsHorizontalIcon className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem>
                    <div>Copy Label ID</div>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
