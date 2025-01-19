import { stageIcon } from "@/components/icons";
import { PriorityIcons } from "@/components/icons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
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
import { threadStatusVerboseName } from "@/db/helpers";
import { ThreadResponse, threadTransformer } from "@/db/models";
import { WorkspaceStoreState } from "@/db/store";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  CaretSortIcon,
  CheckIcon,
  PersonIcon,
  PlusCircledIcon,
} from "@radix-ui/react-icons";
import { useMutation } from "@tanstack/react-query";
import * as React from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";
import { useStore } from "zustand";

const StageFormSchema = z.object({
  stage: z.string(),
});

type StageFormInputs = z.infer<typeof StageFormSchema>;

export function SetThreadStatusForm({
  stage,
  threadId,
  token,
  workspaceId,
}: {
  stage: string;
  threadId: string;
  token: string;
  workspaceId: string;
}) {
  const refSubmitButton = React.useRef<HTMLButtonElement>(null);

  const form = useForm<z.infer<typeof StageFormSchema>>({
    defaultValues: {
      stage: stage,
    },
    resolver: zodResolver(StageFormSchema),
  });
  const workspaceStore = useWorkspaceStore();

  const { formState } = form;
  const { errors } = formState;

  const mutation = useMutation({
    mutationFn: async (inputs: StageFormInputs) => {
      const { stage } = inputs;
      const { data, error } = await updateThread(token, workspaceId, threadId, {
        stage: stage,
      });
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

  const onSubmit: SubmitHandler<StageFormInputs> = async (inputs) => {
    await mutation.mutateAsync(inputs);
  };

  async function onStageChange(value: string) {
    form.setValue("stage", value);
    setTimeout(() => {
      refSubmitButton?.current?.click();
    }, 0);
  }

  return (
    <Form {...form}>
      <form
        className="flex w-52 flex-col gap-1"
        onSubmit={form.handleSubmit(onSubmit)}
      >
        <FormField
          control={form.control}
          name="stage"
          render={({ field }) => (
            <FormItem>
              <Select defaultValue={field.value} onValueChange={onStageChange}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={threadStatusVerboseName(stage)} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="needs_first_response">
                    <div className="flex items-center gap-2">
                      {stageIcon("needs_first_response", {
                        className: "h-4 w-4 text-indigo-500 my-auto",
                      })}
                      <span className="text-sm">
                        {threadStatusVerboseName("needs_first_response")}
                      </span>
                    </div>
                  </SelectItem>
                  <SelectItem value="needs_next_response">
                    <div className="flex items-center gap-2">
                      {stageIcon("needs_next_response", {
                        className: "h-4 w-4 text-indigo-500 my-auto",
                      })}
                      <span className="text-sm">
                        {threadStatusVerboseName("needs_next_response")}
                      </span>
                    </div>
                  </SelectItem>
                  <SelectItem value="waiting_on_customer">
                    <div className="flex items-center gap-2">
                      {stageIcon("waiting_on_customer", {
                        className: "h-4 w-4 text-indigo-500 my-auto",
                      })}
                      <span className="text-sm">
                        {threadStatusVerboseName("waiting_on_customer")}
                      </span>
                    </div>
                  </SelectItem>
                  <SelectItem value="hold">
                    <div className="flex items-center gap-2">
                      {stageIcon("hold", {
                        className: "h-4 w-4 text-indigo-500 my-auto",
                      })}
                      <span className="text-sm">
                        {threadStatusVerboseName("hold")}
                      </span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />
        <button hidden ref={refSubmitButton} type="submit" />
        {errors?.root && (
          <div className="text-xs text-red-500">Something went wrong</div>
        )}
      </form>
    </Form>
  );
}

const AssigneeFormSchema = z.object({
  assignee: z.string(),
});

type AssigneeFormInputs = z.infer<typeof AssigneeFormSchema>;

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
    state.viewMembers(state),
  );
  const memberName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewMemberName(state, value || ""),
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

  const unassignedMembers = [{ memberId: "unassigned", name: "Unassigned" }];
  const membersUpdated = unassignedMembers.concat(members);

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
              <PersonIcon className="h-4 w-4 text-muted-foreground" />
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
                  className="flex gap-2 text-xs"
                  key={member.memberId}
                  onSelect={() => {
                    onValueChange(member.memberId);
                    setOpen(false);
                  }}
                >
                  {member.memberId === "unassigned" || !member.memberId ? (
                    <PersonIcon className="h-4 w-4 text-muted-foreground" />
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
                      value === member.memberId ? "opacity-100" : "opacity-0",
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
                    className="flex gap-2 text-xs"
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

export function SetThreadAssigneeForm({
  assigneeId,
  threadId,
  token,
  workspaceId,
}: {
  assigneeId: string;
  threadId: string;
  token: string;
  workspaceId: string;
}) {
  const refSubmitButton = React.useRef<HTMLButtonElement>(null);

  const form = useForm<z.infer<typeof AssigneeFormSchema>>({
    defaultValues: {
      assignee: assigneeId,
    },
    resolver: zodResolver(AssigneeFormSchema),
  });

  const workspaceStore = useWorkspaceStore();

  const { formState } = form;
  const { errors } = formState;

  React.useEffect(() => {
    form.reset({ assignee: assigneeId });
  }, [form, assigneeId]);

  const mutation = useMutation({
    mutationFn: async (inputs: AssigneeFormInputs) => {
      const { assignee } = inputs;
      const memberId = assignee === "unassigned" ? null : assignee;
      const body = {
        assignee: memberId,
      };
      const { data, error } = await updateThread(
        token,
        workspaceId,
        threadId,
        body,
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

  const onSubmit: SubmitHandler<AssigneeFormInputs> = async (inputs) => {
    await mutation.mutateAsync(inputs);
  };

  async function onAssigneeChange(value: string) {
    form.setValue("assignee", value);
    setTimeout(() => {
      refSubmitButton?.current?.click();
    }, 0);
  }

  return (
    <Form {...form}>
      <form
        className="flex flex-col gap-1"
        onSubmit={form.handleSubmit(onSubmit)}
      >
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
        <button hidden ref={refSubmitButton} type="submit" />
        {errors?.root && (
          <div className="text-xs text-red-500">Something went wrong</div>
        )}
      </form>
    </Form>
  );
}

const PriorityFormSchema = z.object({
  priority: z.string(),
});

type PriorityFormInputs = z.infer<typeof PriorityFormSchema>;

export function SetThreadPriorityForm({
  priority,
  threadId,
  token,
  workspaceId,
}: {
  priority: string;
  threadId: string;
  token: string;
  workspaceId: string;
}) {
  const refSubmitButton = React.useRef<HTMLButtonElement>(null);

  const form = useForm<z.infer<typeof PriorityFormSchema>>({
    defaultValues: {
      priority: priority,
    },
    resolver: zodResolver(PriorityFormSchema),
  });
  const workspaceStore = useWorkspaceStore();

  const { formState } = form;
  const { errors } = formState;

  React.useEffect(() => {
    form.reset({ priority: priority });
  }, [form, priority]);

  const mutation = useMutation({
    mutationFn: async (inputs: PriorityFormInputs) => {
      const { priority } = inputs;
      const body = {
        priority: priority,
      };
      const { data, error } = await updateThread(
        token,
        workspaceId,
        threadId,
        body,
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

  const onSubmit: SubmitHandler<PriorityFormInputs> = async (inputs) => {
    await mutation.mutateAsync(inputs);
  };

  async function onPriorityChange(value: string) {
    form.setValue("priority", value);
    setTimeout(() => {
      refSubmitButton?.current?.click();
    }, 0);
  }

  return (
    <Form {...form}>
      <form
        className="flex flex-col gap-1"
        onSubmit={form.handleSubmit(onSubmit)}
      >
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
                      <span className="text-sm">Urgent</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="high">
                    <div className="flex items-center gap-1">
                      <PriorityIcons.high className="h-5 w-5" />
                      <span className="text-sm">High</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="normal">
                    <div className="flex items-center gap-1">
                      <PriorityIcons.normal className="h-5 w-5" />
                      <span className="text-sm">Normal</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="low">
                    <div className="flex items-center gap-1">
                      <PriorityIcons.low className="h-5 w-5" />
                      <span className="text-sm">Low</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />
        <button hidden ref={refSubmitButton} type="submit">
          Submit
        </button>
        {errors?.root && (
          <div className="text-xs text-red-500">Something went wrong</div>
        )}
      </form>
    </Form>
  );
}
