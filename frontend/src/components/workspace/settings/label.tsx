import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { updateWorkspaceLabel } from "@/db/api";
import { useWorkspaceStore } from "@/providers";
import { zodResolver } from "@hookform/resolvers/zod";
import { DotsHorizontalIcon, Pencil1Icon } from "@radix-ui/react-icons";
import { useMutation } from "@tanstack/react-query";
import { TagIcon } from "lucide-react";
import * as React from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";

type FormInputs = {
  name: string;
};

const formSchema = z.object({
  name: z.string().min(2),
});

export function LabelItem({
  label,
  labelId,
  token,
  workspaceId,
}: {
  label: string;
  labelId: string;
  token: string;
  workspaceId: string;
}) {
  const [editMode, setEditMode] = React.useState(false);
  const workspaceStore = useWorkspaceStore();

  const form = useForm({
    defaultValues: {
      name: label,
    },
    resolver: zodResolver(formSchema),
  });

  const mutation = useMutation({
    mutationFn: async (inputs: { name: string }) => {
      const { name } = inputs;
      const { data, error } = await updateWorkspaceLabel(
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
        message: "Something went wrong. Please try again later.",
        type: "serverError",
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
              className="flex justify-between space-x-2"
              onSubmit={form.handleSubmit(onSubmit)}
            >
              <div className="flex items-center w-full">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem className="w-full">
                      <FormControl>
                        <Input
                          autoComplete="off"
                          autoFocus
                          placeholder="Label"
                          required
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
                  onClick={() => setEditMode(false)}
                  size="default"
                  type="button"
                  variant="outline"
                >
                  Cancel
                </Button>
                <Button disabled={isUpdating} size="default" type="submit">
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
                size="icon"
                variant="ghost"
              >
                <Pencil1Icon className="h-4 w-4" />
              </Button>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button size="icon" variant="ghost">
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
