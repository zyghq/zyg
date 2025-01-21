import { Icons, stageIcon } from "@/components/icons";
import { PriorityIcons } from "@/components/icons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
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
  DropdownMenu,
  DropdownMenuContent,
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
import { Label as LabelComponent } from "@/components/ui/label";
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
import { deleteThreadLabel, getThreadLabels, putThreadLabel } from "@/db/api";
import { threadStatusVerboseName } from "@/db/helpers";
import { ThreadResponse, threadTransformer } from "@/db/models";
import { Label, ThreadLabelResponse } from "@/db/models";
import { WorkspaceStoreState } from "@/db/store";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  BorderDashedIcon,
  CaretSortIcon,
  CheckIcon,
  PersonIcon,
  PlusCircledIcon,
  PlusIcon,
} from "@radix-ui/react-icons";
import { useQuery } from "@tanstack/react-query";
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

  const { formState, setError, setValue } = form;
  const { errors } = formState;

  const workspaceStore = useWorkspaceStore();

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
      setError("root", {
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
    setValue("stage", value);
    setTimeout(() => {
      refSubmitButton?.current?.click();
    }, 0);
  }

  React.useEffect(() => {
    setValue("stage", stage);
  }, [stage, setValue]);

  return (
    <Form {...form}>
      <form
        className="flex w-56 flex-col gap-1"
        onSubmit={form.handleSubmit(onSubmit)}
      >
        <FormField
          control={form.control}
          name="stage"
          render={({ field }) => (
            <FormItem>
              <Select onValueChange={onStageChange} value={field.value}>
                <FormControl>
                  <SelectTrigger
                    aria-label="Select Thread Status"
                    className="border-none shadow-none hover:bg-accent focus:ring-0"
                  >
                    <SelectValue placeholder={threadStatusVerboseName(stage)} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent align="end">
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
            aria-label="Select Assignee"
            className="flex w-56 gap-x-2 border-none shadow-none hover:bg-accent focus:ring-0"
            role="combobox"
            variant="ghost"
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
        <PopoverContent align="end" className="w-56 p-0">
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
              <LabelComponent htmlFor="name">Name</LabelComponent>
              <Input autoComplete="off" id="name" />
            </div>
            <div className="space-y-2">
              <LabelComponent htmlFor="name">Email</LabelComponent>
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
        className="flex w-56 flex-col gap-1"
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

  const { formState, setError, setValue } = form;
  const { errors } = formState;

  const workspaceStore = useWorkspaceStore();

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
      setError("root", {
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
    setValue("priority", value);
    setTimeout(() => {
      refSubmitButton?.current?.click();
    }, 0);
  }

  React.useEffect(() => {
    setValue("priority", priority);
  }, [form, priority]);

  return (
    <Form {...form}>
      <form
        className="flex w-56 flex-col gap-1"
        onSubmit={form.handleSubmit(onSubmit)}
      >
        <FormField
          control={form.control}
          name="priority"
          render={({ field }) => (
            <FormItem>
              <Select onValueChange={onPriorityChange} value={field.value}>
                <FormControl>
                  <SelectTrigger
                    aria-label="Select Thread Priority"
                    className="border-none shadow-none hover:bg-accent focus:ring-0"
                  >
                    <SelectValue placeholder="Priority" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent align="end">
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

export function ThreadLabels({
  threadId,
  token,
  workspaceId,
  workspaceLabels,
}: {
  threadId: string;
  token: string;
  workspaceId: string;
  workspaceLabels: Label[];
}) {
  const {
    data: threadLabels,
    error,
    isPending,
    refetch,
  } = useQuery({
    enabled: !!threadId,
    queryFn: async () => {
      const { data, error } = await getThreadLabels(
        token,
        workspaceId,
        threadId,
      );
      if (error) throw new Error("failed to fetch thread labels");
      return data as ThreadLabelResponse[];
    },
    queryKey: ["threadLabels", workspaceId, threadId, token],
  });

  const threadLabelMutation = useMutation({
    mutationFn: async (values: { icon: string; name: string }) => {
      const { data, error } = await putThreadLabel(
        token,
        workspaceId,
        threadId,
        values,
      );
      if (error) {
        throw new Error(error.message);
      }

      if (!data) {
        throw new Error("no data returned");
      }
      return data as ThreadLabelResponse;
    },
    onError: (error) => {
      console.error(error);
    },
    onSuccess: async () => {
      await refetch();
    },
  });

  const deleteThreadLabelMutation = useMutation({
    mutationFn: async (labelId: string) => {
      const { data, error } = await deleteThreadLabel(
        token,
        workspaceId,
        threadId,
        labelId,
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
    },
    onSuccess: async () => {
      await refetch();
    },
  });

  const isChecked = (labelId: string) => {
    return threadLabels?.some((label) => label.labelId === labelId);
  };

  function onSelect(labelId: string, name: string, icon?: string) {
    if (isChecked(labelId)) {
      deleteThreadLabelMutation.mutate(labelId);
    } else {
      threadLabelMutation.mutate({ icon: icon || "", name });
    }
  }

  const renderLabels = () => {
    if (isPending) {
      return (
        <Badge variant="outline">
          <BorderDashedIcon className="h-4 w-4" />
        </Badge>
      );
    }

    if (error) {
      return (
        <div className="flex items-center gap-1">
          <Icons.oops className="h-5 w-5" />
          <div className="text-xs text-red-500">Something went wrong</div>
        </div>
      );
    }

    return (
      <React.Fragment>
        {threadLabels.length > 0 ? (
          threadLabels?.map((label) => (
            <Badge key={label.labelId} variant="outline">
              <div className="flex items-center gap-1">
                <div>{label.icon}</div>
                <div className="capitalize text-muted-foreground">
                  {label.name}
                </div>
              </div>
            </Badge>
          ))
        ) : (
          <Badge variant="outline">
            <BorderDashedIcon className="h-4 w-4" />
          </Badge>
        )}
      </React.Fragment>
    );
  };

  return (
    <div className="flex flex-col">
      <div className="flex items-center justify-between">
        <div className="items-center font-serif text-sm font-medium">
          Labels
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button className="border-dashed" size="sm" variant="outline">
              <PlusIcon className="mr-1 h-3 w-3" />
              Add
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <Command>
              <CommandList>
                <CommandInput placeholder="Filter" />
                <CommandEmpty>No results</CommandEmpty>
                <CommandGroup>
                  {workspaceLabels.map((label) => (
                    <CommandItem
                      className="text-sm"
                      key={label.labelId}
                      onSelect={() =>
                        onSelect(label.labelId, label.name, label.icon)
                      }
                    >
                      <div className="flex gap-2">
                        <div>{label.icon}</div>
                        <div className="capitalize">{label.name}</div>
                      </div>
                      <CheckIcon
                        className={cn(
                          "ml-auto h-4 w-4",
                          isChecked(label.labelId)
                            ? "opacity-100"
                            : "opacity-0",
                        )}
                      />
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="flex flex-wrap gap-1 py-1">{renderLabels()}</div>
      {threadLabelMutation.isError && (
        <div className="mt-1 text-xs text-red-500">Something went wrong</div>
      )}
    </div>
  );
}
