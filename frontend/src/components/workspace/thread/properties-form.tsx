import * as React from "react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import { z } from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { PriorityIcons } from "@/components/icons";

import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import Avatar from "boring-avatars";
import {
  CaretSortIcon,
  CheckIcon,
  PlusCircledIcon,
} from "@radix-ui/react-icons";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { useStore } from "zustand";
import { useMutation } from "@tanstack/react-query";
import { updateThread } from "@/db/api";

const FormSchema = z.object({
  priority: z.string(),
  assignee: z.string(),
});

type FormInputs = z.infer<typeof FormSchema>;

function SetAssignee({
  value,
  onValueChange,
}: {
  value: string;
  onValueChange: (value: string) => void;
}) {
  const workspaceStore = useWorkspaceStore();
  const [open, setOpen] = React.useState(false);
  const [showNewTeamDialog, setShowNewTeamDialog] = React.useState(false);

  const members = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewMembers(state)
  );
  const memberName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewMemberName(state, value || "")
  );

  const renderMemberName = () => {
    if (value === "unassigned") {
      return "Unassigned";
    }
    if (memberName === "") {
      return "n/a";
    }
    return memberName;
  };

  const membersUpdated = [{ memberId: "unassigned", name: "Unassigned" }];
  members.forEach((member) => {
    membersUpdated.push({
      memberId: member.memberId,
      name: member.name || "n/a",
    });
  });

  return (
    <Dialog open={showNewTeamDialog} onOpenChange={setShowNewTeamDialog}>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            aria-label="Select a team"
            className="flex gap-1"
          >
            {value === "unassigned" || !value ? (
              <Avatar
                name={value || "unassigned"}
                size={18}
                colors={["#e4e4ea"]}
              />
            ) : (
              <Avatar name={value} size={18} />
            )}
            {renderMemberName()}
            <CaretSortIcon className="ml-auto h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0" align="start">
          <Command>
            <CommandList>
              <CommandInput placeholder="Search member..." />
              <CommandEmpty>No member found.</CommandEmpty>
              {membersUpdated.map((member) => (
                <CommandItem
                  key={member.memberId}
                  onSelect={() => {
                    onValueChange(member.memberId);
                    setOpen(false);
                  }}
                  className="text-xs flex gap-1"
                >
                  {member.memberId === "unassigned" || !member.memberId ? (
                    <Avatar
                      name={member.memberId || "unassigned"}
                      size={18}
                      colors={["#e4e4ea"]}
                    />
                  ) : (
                    <Avatar name={member.memberId} size={18} />
                  )}
                  {member.name}
                  <CheckIcon
                    className={cn(
                      "ml-auto h-4 w-4",
                      value === member.memberId ? "opacity-100" : "opacity-0"
                    )}
                  />
                </CommandItem>
              ))}
            </CommandList>
            <CommandSeparator />
            <CommandList>
              <CommandGroup>
                <DialogTrigger asChild>
                  <CommandItem
                    className="text-xs flex gap-1"
                    onSelect={() => {
                      setOpen(false);
                      setShowNewTeamDialog(true);
                    }}
                  >
                    <PlusCircledIcon className="h-4 w-4" />
                    Invite Member
                  </CommandItem>
                </DialogTrigger>
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Invite Member</DialogTitle>
          <DialogDescription>Invite a member to your team.</DialogDescription>
        </DialogHeader>
        <div>
          <div className="space-y-4 py-2 pb-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" autoComplete="off" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Email</Label>
              <Input type="email" id="email" placeholder="name@example.com" />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setShowNewTeamDialog(false)}>
            Cancel
          </Button>
          <Button type="submit">Continue</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export function PropertiesForm({
  token,
  workspaceId,
  threadId,
  priority,
  assigneeId,
}: {
  token: string;
  workspaceId: string;
  threadId: string;
  priority: string;
  assigneeId: string;
}) {
  const refSubmitButtom = React.useRef<HTMLButtonElement>(null);

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      priority: priority,
      assignee: assigneeId,
    },
  });
  const workspaceStore = useWorkspaceStore();

  const { formState } = form;
  const { errors } = formState;

  React.useEffect(() => {
    form.reset({ priority: priority, assignee: assigneeId });
  }, [form, assigneeId, priority]);

  const mutation = useMutation({
    mutationFn: async (inputs: FormInputs) => {
      const { assignee, priority } = inputs;
      const memberId = assignee === "unassigned" ? null : assignee;
      const body = {
        assignee: memberId,
        priority,
      };
      const { error, data } = await updateThread(
        token,
        workspaceId,
        threadId,
        body
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
      form.setError("root", {
        type: "serverError",
        message: "Something went wrong. Please try again later.",
      });
    },
    onSuccess: (data) => {
      workspaceStore.getState().updateThread(data);
    },
  });

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    await mutation.mutateAsync(inputs);
  };

  async function onPriorityChange(value: string) {
    form.setValue("priority", value);
    setTimeout(() => {
      refSubmitButtom?.current?.click();
    }, 0);
  }

  async function onAssigneeChange(value: string) {
    form.setValue("assignee", value);
    setTimeout(() => {
      refSubmitButtom?.current?.click();
    }, 0);
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col gap-2"
      >
        <div className="flex gap-1">
          <FormField
            control={form.control}
            name="priority"
            render={({ field }) => (
              <FormItem>
                <Select
                  onValueChange={onPriorityChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Priority" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="urgent">
                      <div className="flex items-center gap-1">
                        <PriorityIcons.urgent className="h-5 w-5" />
                        <span>Urgent</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="high">
                      <div className="flex items-center gap-1">
                        <PriorityIcons.high className="h-5 w-5" />
                        <span>High</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="normal">
                      <div className="flex items-center gap-1">
                        <PriorityIcons.normal className="h-5 w-5" />
                        <span>Normal</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="low">
                      <div className="flex items-center gap-1">
                        <PriorityIcons.low className="h-5 w-5" />
                        <span>Low</span>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="assignee"
            render={({ field }) => (
              <FormItem>
                <SetAssignee
                  value={field.value}
                  onValueChange={onAssigneeChange}
                />
                <FormMessage />
              </FormItem>
            )}
          />
          <button ref={refSubmitButtom} hidden type="submit">
            Submit
          </button>
        </div>
        {errors?.root && (
          <div className="text-xs text-red-500">Something went wrong</div>
        )}
      </form>
    </Form>
  );
}
