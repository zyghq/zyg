import { PriorityIcons } from "@/components/icons";
import { Button } from "@/components/ui/button";
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
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { updateThread } from "@/db/api";
import { threadTransformer } from "@/db/models";
import { ThreadResponse } from "@/db/schema";
import { WorkspaceStoreState } from "@/db/store";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  CaretSortIcon,
  CheckIcon,
  PlusCircledIcon,
} from "@radix-ui/react-icons";
import { useMutation } from "@tanstack/react-query";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import * as React from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";
import { useStore } from "zustand";

const FormSchema = z.object({
  assignee: z.string(),
  priority: z.string(),
});

type FormInputs = z.infer<typeof FormSchema>;

function SetAssignee({
  onValueChange,
  value,
}: {
  onValueChange: (value: string) => void;
  value: string;
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
    <Dialog onOpenChange={setShowNewTeamDialog} open={showNewTeamDialog}>
      <Popover onOpenChange={setOpen} open={open}>
        <PopoverTrigger asChild>
          <Button
            aria-expanded={open}
            aria-label="Select a team"
            className="flex gap-x-2"
            role="combobox"
            variant="outline"
          >
            {value === "unassigned" || !value ? (
              <Avatar className="h-5 w-5">
                <AvatarImage src={`https://avatar.vercel.sh/unassigned`} />
                <AvatarFallback>U</AvatarFallback>
              </Avatar>
            ) : (
              <Avatar className="h-5 w-5">
                <AvatarImage src={`https://avatar.vercel.sh/${value}`} />
                <AvatarFallback>M</AvatarFallback>
              </Avatar>
            )}
            {renderMemberName()}
            <CaretSortIcon className="ml-auto h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-[200px] p-0">
          <Command>
            <CommandList>
              <CommandInput placeholder="Search member..." />
              <CommandEmpty>No member found.</CommandEmpty>
              {membersUpdated.map((member) => (
                <CommandItem
                  className="text-xs flex gap-2"
                  key={member.memberId}
                  onSelect={() => {
                    onValueChange(member.memberId);
                    setOpen(false);
                  }}
                >
                  {member.memberId === "unassigned" || !member.memberId ? (
                    <Avatar className="h-5 w-5">
                      <AvatarImage
                        src={`https://avatar.vercel.sh/unassigned`}
                      />
                      <AvatarFallback>M</AvatarFallback>
                    </Avatar>
                  ) : (
                    <Avatar className="h-5 w-5">
                      <AvatarImage
                        src={`https://avatar.vercel.sh/${member.memberId}`}
                      />
                      <AvatarFallback>M</AvatarFallback>
                    </Avatar>
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
                    className="text-xs flex gap-2"
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
              <Input autoComplete="off" id="name" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Email</Label>
              <Input id="email" placeholder="name@example.com" type="email" />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button onClick={() => setShowNewTeamDialog(false)} variant="outline">
            Cancel
          </Button>
          <Button type="submit">Continue</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export function PropertiesForm({
  assigneeId,
  priority,
  threadId,
  token,
  workspaceId,
}: {
  assigneeId: string;
  priority: string;
  threadId: string;
  token: string;
  workspaceId: string;
}) {
  const refSubmitButtom = React.useRef<HTMLButtonElement>(null);

  const form = useForm<z.infer<typeof FormSchema>>({
    defaultValues: {
      assignee: assigneeId,
      priority: priority,
    },
    resolver: zodResolver(FormSchema),
  });
  const workspaceStore = useWorkspaceStore();

  const { formState } = form;
  const { errors } = formState;

  React.useEffect(() => {
    form.reset({ assignee: assigneeId, priority: priority });
  }, [form, assigneeId, priority]);

  const mutation = useMutation({
    mutationFn: async (inputs: FormInputs) => {
      const { assignee, priority } = inputs;
      const memberId = assignee === "unassigned" ? null : assignee;
      const body = {
        assignee: memberId,
        priority,
      };
      const { data, error } = await updateThread(
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
      return data as ThreadResponse;
    },
    onError: (error) => {
      console.error(error);
      form.setError("root", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    },
    onSuccess: (data) => {
      const transformer = threadTransformer();
      const [, thread] = transformer.normalize(data);
      workspaceStore.getState().updateThread(thread);
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
        className="flex flex-col gap-2 px-4 py-2"
        onSubmit={form.handleSubmit(onSubmit)}
      >
        <div className="flex gap-1">
          <FormField
            control={form.control}
            name="priority"
            render={({ field }) => (
              <FormItem>
                <Select
                  defaultValue={field.value}
                  onValueChange={onPriorityChange}
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
                  onValueChange={onAssigneeChange}
                  value={field.value}
                />
                <FormMessage />
              </FormItem>
            )}
          />
          <button hidden ref={refSubmitButtom} type="submit">
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
